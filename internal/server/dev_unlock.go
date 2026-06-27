package server

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/SorynMochi/Namagotchi/internal/database"
)

const devUnlockCookieName = "namigotchi_dev_unlock"

const devUnlockTTL = 30 * time.Minute

type devUnlockRequest struct {
	Passphrase string `json:"passphrase"`
}

type devUnlockResponse struct {
	OK        bool      `json:"ok"`
	Message   string    `json:"message"`
	ExpiresAt time.Time `json:"expiresAt,omitempty"`
}

const devConsoleUnlockHTML = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Unlock Namigotchi Dev Console</title>
  <style>
    * {
      box-sizing: border-box;
    }

    body {
      margin: 0;
      min-height: 100vh;
      display: grid;
      place-items: center;
      padding: 24px;
      font-family: system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
      background: #120c18;
      color: #fcefff;
    }

    main {
      width: min(420px, calc(100vw - 32px));
      padding: 22px;
      border: 1px solid rgba(255,255,255,0.16);
      border-radius: 18px;
      background: #21162e;
      box-shadow: 0 24px 60px rgba(0,0,0,0.35);
    }

    h1 {
      margin: 0 0 10px;
      font-size: clamp(1.7rem, 5vw, 2.2rem);
      line-height: 1.05;
    }

    p {
      color: #d9c6e5;
      line-height: 1.45;
    }

    label {
      display: grid;
      gap: 8px;
      margin-top: 16px;
      font-weight: 800;
    }

    input {
      display: block;
      width: 100%;
      max-width: 100%;
      padding: 11px 13px;
      border: 1px solid rgba(255,255,255,0.18);
      border-radius: 12px;
      background: rgba(0,0,0,0.35);
      color: #fcefff;
      font: inherit;
    }

    input:focus {
      outline: 2px solid rgba(255,143,199,0.75);
      outline-offset: 2px;
    }

    button {
      width: 100%;
      margin-top: 16px;
      border: 0;
      border-radius: 999px;
      padding: 11px 14px;
      cursor: pointer;
      background: #ff8fc7;
      color: #1f1022;
      font-weight: 900;
    }

    a {
      color: #ff8fc7;
      font-weight: 800;
    }

    .error {
      color: #ff8a8a;
      min-height: 1.4em;
    }
  </style>
</head>
<body>
  <main>
    <p><a href="/">Back to game</a></p>
    <h1>Unlock Dev Console</h1>
    <p>Your account is allowed to use developer tools. Enter the separate dev unlock passphrase to open the console for this session.</p>

    <form id="dev-unlock-form">
      <label>
        Dev unlock passphrase
        <input id="dev-unlock-passphrase" type="password" autocomplete="current-password" required>
      </label>

      <button type="submit">Unlock Console</button>
      <p id="dev-unlock-error" class="error"></p>
    </form>
  </main>

  <script>
    const form = document.querySelector("#dev-unlock-form");
    const input = document.querySelector("#dev-unlock-passphrase");
    const error = document.querySelector("#dev-unlock-error");
    const CSRF_COOKIE_NAME = "namigotchi_csrf";
    const CSRF_HEADER_NAME = "X-CSRF-Token";
    let csrfTokenPromise = null;

    function readCookieValue(name) {
      return document.cookie
        .split(";")
        .map((part) => part.trim())
        .find((part) => part.startsWith(name + "="))
        ?.slice(name.length + 1) || "";
    }

    async function ensureCSRFToken() {
      const existingToken = readCookieValue(CSRF_COOKIE_NAME);
      if (existingToken) {
        return existingToken;
      }

      if (!csrfTokenPromise) {
        csrfTokenPromise = fetch("/api/auth/csrf")
          .then((response) => {
            if (!response.ok) {
              throw new Error("CSRF token request failed: " + response.status);
            }

            return response.json();
          })
          .then((payload) => payload.csrfToken || readCookieValue(CSRF_COOKIE_NAME))
          .finally(() => {
            csrfTokenPromise = null;
          });
      }

      return csrfTokenPromise;
    }

    async function csrfFetch(input, options = {}) {
      const requestOptions = { ...options };
      const method = String(requestOptions.method || "GET").toUpperCase();
      const headers = new Headers(requestOptions.headers || {});

      if (!["GET", "HEAD", "OPTIONS"].includes(method)) {
        const token = await ensureCSRFToken();
        headers.set(CSRF_HEADER_NAME, token);
      }

      requestOptions.headers = headers;

      return fetch(input, requestOptions);
    }

    form.addEventListener("submit", async (event) => {
      event.preventDefault();

      error.textContent = "";

      try {
        const response = await csrfFetch("/api/dev/unlock", {
          method: "POST",
          headers: {
            "Content-Type": "application/json"
          },
          body: JSON.stringify({
            passphrase: input.value
          })
        });

        const payload = await response.json().catch(() => null);

        if (!response.ok) {
          error.textContent = payload && payload.message ? payload.message : "Unlock failed.";
          return;
        }

        window.location.reload();
      } catch {
        error.textContent = "Unlock request failed.";
      }
    });
  </script>
</body>
</html>`

func (s *Server) requireDevUnlock(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !s.requestHasDevUnlock(r) {
			writeError(w, http.StatusForbidden, "dev console unlock required")
			return
		}

		next(w, r)
	}
}

func (s *Server) HandleDevUnlock(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	secret := devUnlockSecret()
	if secret == "" {
		writeError(w, http.StatusServiceUnavailable, "dev unlock secret is not configured")
		return
	}

	accountID, ok := database.AuthAccountIDFromContext(r.Context())
	if !ok || accountID < 1 {
		writeError(w, http.StatusUnauthorized, "login required")
		return
	}

	var request devUnlockRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid unlock request")
		return
	}

	if !constantTimeStringEqual(strings.TrimSpace(request.Passphrase), secret) {
		clearDevUnlockCookie(w, r)
		writeError(w, http.StatusForbidden, "dev unlock denied")
		return
	}

	expiresAt := time.Now().Add(devUnlockTTL)
	token := makeDevUnlockToken(accountID, expiresAt, secret)

	setDevUnlockCookie(w, r, token, expiresAt)

	writeJSON(w, http.StatusOK, devUnlockResponse{
		OK:        true,
		Message:   "Dev console unlocked.",
		ExpiresAt: expiresAt,
	})
}

func (s *Server) requestHasDevUnlock(r *http.Request) bool {
	secret := devUnlockSecret()
	if secret == "" {
		return false
	}

	accountID, ok := database.AuthAccountIDFromContext(r.Context())
	if !ok || accountID < 1 {
		return false
	}

	cookie, err := r.Cookie(devUnlockCookieName)
	if err != nil {
		return false
	}

	return verifyDevUnlockToken(strings.TrimSpace(cookie.Value), accountID, secret, time.Now())
}

func devUnlockSecret() string {
	return strings.TrimSpace(os.Getenv("NAMIGOTCHI_DEV_UNLOCK_SECRET"))
}

func setDevUnlockCookie(w http.ResponseWriter, r *http.Request, token string, expiresAt time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     devUnlockCookieName,
		Value:    token,
		Path:     "/",
		MaxAge:   int(devUnlockTTL.Seconds()),
		Expires:  expiresAt,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   authCookieShouldBeSecure(r),
	})
}

func clearDevUnlockCookie(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     devUnlockCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   authCookieShouldBeSecure(r),
	})
}

func makeDevUnlockToken(accountID int64, expiresAt time.Time, secret string) string {
	payload := fmt.Sprintf("%d.%d", accountID, expiresAt.Unix())
	signature := signDevUnlockPayload(payload, secret)

	return payload + "." + signature
}

func verifyDevUnlockToken(token string, accountID int64, secret string, now time.Time) bool {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return false
	}

	tokenAccountID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil || tokenAccountID != accountID {
		return false
	}

	expiresUnix, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return false
	}

	if !now.Before(time.Unix(expiresUnix, 0)) {
		return false
	}

	payload := parts[0] + "." + parts[1]
	expectedSignature := signDevUnlockPayload(payload, secret)

	return subtle.ConstantTimeCompare([]byte(parts[2]), []byte(expectedSignature)) == 1
}

func signDevUnlockPayload(payload string, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(payload))

	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func constantTimeStringEqual(left string, right string) bool {
	leftHash := sha256.Sum256([]byte(left))
	rightHash := sha256.Sum256([]byte(right))

	return subtle.ConstantTimeCompare(leftHash[:], rightHash[:]) == 1
}

func (s *Server) HandleDevLock(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	clearDevUnlockCookie(w, r)

	writeJSON(w, http.StatusOK, devUnlockResponse{
		OK:      true,
		Message: "Dev console locked.",
	})
}
