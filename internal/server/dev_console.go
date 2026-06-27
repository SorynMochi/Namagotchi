package server

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

const devConsoleHTML = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Namigotchi Dev Console</title>
  <style>
    * {
      box-sizing: border-box;
    }

    body {
      margin: 0;
      min-height: 100vh;
      font-family: system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
      background: #120c18;
      color: #fcefff;
    }

    main {
      max-width: 1120px;
      margin: 0 auto;
      padding: 28px 18px 40px;
    }

    a {
      color: #ff8fc7;
      font-weight: 800;
    }

    h1 {
      margin: 28px 0 12px;
      font-size: clamp(1.8rem, 4vw, 2.5rem);
      line-height: 1.05;
    }

    p {
      color: #d9c6e5;
      line-height: 1.45;
    }

    .grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(170px, 1fr));
      gap: 10px;
      margin-top: 20px;
    }

    button {
      width: 100%;
      min-width: 0;
      min-height: 40px;
      border: 0;
      border-radius: 999px;
      padding: 9px 14px;
      cursor: pointer;
      background: #ff8fc7;
      color: #1f1022;
      font-size: 0.9rem;
      font-weight: 900;
      line-height: 1.1;
      text-align: center;
      white-space: nowrap;
    }

    button:hover {
      filter: brightness(1.05);
    }

    button:focus-visible {
      outline: 2px solid rgba(255,255,255,0.75);
      outline-offset: 2px;
    }

    button.secondary {
      width: 100%;
      margin-bottom: 18px;
      color: #fcefff;
      background: rgba(255,255,255,0.12);
      border: 1px solid rgba(255,255,255,0.18);
    }

    pre {
      min-height: 190px;
      max-height: 420px;
      margin-top: 22px;
      padding: 14px;
      overflow: auto;
      white-space: pre-wrap;
      border-radius: 14px;
      border: 1px solid rgba(255,255,255,0.16);
      background: rgba(0,0,0,0.35);
    }

    .reset-server-zone {
      margin-top: 16px;
      padding-top: 16px;
      border-top: 1px solid rgba(255,255,255,0.14);
    }

    .reset-server-zone-label {
      margin: 0 0 10px;
      color: #ff8a98;
      font-size: 0.78rem;
      font-weight: 900;
      letter-spacing: 0.12em;
      text-transform: uppercase;
    }

    .server-reset-button {
      width: min(260px, 100%);
      background: #ff4f64;
      color: #fff7fb;
      letter-spacing: 0.08em;
    }

    .reset-server-overlay {
      position: fixed;
      inset: 0;
      z-index: 9999;
      display: grid;
      place-items: center;
      padding: 18px;
      background: rgba(0,0,0,0.72);
    }

    .reset-server-modal {
      width: min(560px, 100%);
      padding: 22px;
      border: 1px solid rgba(255,255,255,0.22);
      border-radius: 22px;
      background: #24162d;
      box-shadow: 0 30px 90px rgba(0,0,0,0.55);
    }

    .reset-server-modal h2 {
      margin: 0 0 10px;
      color: #ff8a98;
      font-size: clamp(1.4rem, 5vw, 2rem);
      letter-spacing: 0.04em;
    }

    .reset-server-modal p,
    .reset-server-modal li {
      color: #f1ddeb;
    }

    .reset-server-modal ul {
      margin: 12px 0 18px;
      padding-left: 22px;
    }

    .reset-server-actions {
      display: flex;
      flex-wrap: wrap;
      gap: 10px;
      margin-top: 16px;
    }

    .reset-server-actions button {
      width: auto;
      min-width: 180px;
    }

    .reset-server-confirm {
      background: #ff4f64;
      color: #fff7fb;
    }

    .reset-server-cancel {
      background: rgba(255,255,255,0.12);
      color: #fcefff;
      border: 1px solid rgba(255,255,255,0.18);
    }

    @media (max-width: 560px) {
      .grid {
        grid-template-columns: 1fr 1fr;
      }

      button {
        min-height: 38px;
        font-size: 0.82rem;
        white-space: normal;
      }

      .reset-server-actions button {
        width: 100%;
      }
    }
  </style>
</head>
<body>
  <main>
    <p><a href="/">Back to game</a></p>
    <button id="dev-lock-button" class="secondary" type="button">Lock Dev Console</button>
    <h1>Namigotchi Dev Console</h1>
    <p>This page is served only after the backend verifies your dev access.</p>
    <section class="grid" aria-label="Dev commands">
      <button data-endpoint="/api/dev/seed-player">Seed Player</button>
      <button data-endpoint="/api/dev/force-tick">Force Tick</button>
      <button data-endpoint="/api/dev/reset-chain" data-method="POST" data-inputs="playerName">Reset Chain</button>
      <button data-endpoint="/api/dev/reset-max-chain" data-method="POST" data-inputs="playerName">Reset Max Chain</button>
      <button data-endpoint="/api/dev/spawn-wardrobe-item">Spawn Random Item</button>
      <button data-endpoint="/api/dev/clear-wardrobe" data-method="POST" data-inputs="playerName">Clear Wardrobe</button>
      <button data-endpoint="/api/dev/add-currency" data-method="POST" data-inputs="playerName,currencyType,currencyAmount">Add Currency</button>
      <button data-endpoint="/api/dev/remove-currency" data-method="POST" data-inputs="playerName,currencyType,currencyAmount">Remove Currency</button>
      <button data-endpoint="/api/dev/reset-levels" data-method="POST" data-inputs="playerName,activityName">Reset Levels</button>
      <button data-endpoint="/api/dev/audit-logs" data-method="GET">Refresh Audit Logs</button>
      <button data-endpoint="/api/dev/finish-care" data-method="POST">Finish Care</button>
      <button data-endpoint="/api/dev/security-events" data-method="GET">Security Logs</button>
      <button data-endpoint="/api/dev/account-age" data-method="GET">Show Account Age</button>
    </section>

    <pre id="dev-log">Ready.</pre>

    <section class="reset-server-zone" aria-label="Dangerous dev commands">
      <p class="reset-server-zone-label">Danger Zone</p>
      <button id="reset-server-button" class="server-reset-button" type="button">RESET SERVER</button>
    </section>
  </main>

  <script>
    const log = document.querySelector("#dev-log");
    const lockButton = document.querySelector("#dev-lock-button");
    const resetServerButton = document.querySelector("#reset-server-button");
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

    function writeLog(value) {
      if (typeof value === "string") {
        log.textContent = value;
        return;
      }

      log.textContent = JSON.stringify(value, null, 2);
    }

    async function lockDevConsole() {
      if (!lockButton) {
        return;
      }

      const originalText = lockButton.textContent;
      lockButton.disabled = true;
      lockButton.textContent = "Locking...";

      try {
        await csrfFetch("/api/dev/lock", { method: "POST" });
      } finally {
        lockButton.disabled = false;
        lockButton.textContent = originalText;
        window.location.reload();
      }
    }

    const DEV_INPUT_LABELS = {
      playerName: "Player name, exact match only. Use RESETALL where supported, or GIVEALL for Add Currency.",
      currencyType: "Asset type: Credits, Nibbles, NamiCoins, Fans, Memes, Lost Items, Confidence, Receipts, Patterns, or ALL.",
      currencyAmount: "Amount. Supports 50k, 1m, 100b, etc. Remove Currency also supports ALL.",
      activityName: "Activity/level to reset: Playdeck Level, Nami Level, Streaming, Doom Scrolling, Cleaning, Exercise, Shopping, Designing, or ALL."
    };

    function devCommandInputKeys(button) {
      return String(button.dataset.inputs || "")
        .split(",")
        .map((value) => value.trim())
        .filter(Boolean);
    }

    function buildDevCommandPayload(button) {
      const inputKeys = devCommandInputKeys(button);

      if (inputKeys.length === 0) {
        return null;
      }

      const payload = {};

      for (const key of inputKeys) {
        const label = DEV_INPUT_LABELS[key] || key;
        const value = window.prompt(label);

        if (value === null) {
          const cancelled = new Error("cancelled");
          cancelled.name = "DevCommandCancelled";
          throw cancelled;
        }

        payload[key] = value.trim();
      }

      return payload;
    }

    async function runDevCommand(endpoint, button) {
      const originalText = button.textContent;

      button.disabled = true;
      button.textContent = "Running...";

      try {
        const method = button.dataset.method || "POST";
        const requestOptions = { method };
        const payload = buildDevCommandPayload(button);

        if (payload) {
          requestOptions.headers = {
            "Content-Type": "application/json"
          };
          requestOptions.body = JSON.stringify(payload);
        }

        const response = await csrfFetch(endpoint, requestOptions);
        const text = await response.text();

        let responsePayload = text;
        try {
          responsePayload = text ? JSON.parse(text) : {};
        } catch {
          responsePayload = text;
        }

        if (!response.ok) {
          writeLog({
            ok: false,
            status: response.status,
            response: responsePayload
          });
          return;
        }

        writeLog(responsePayload);
      } catch (error) {
        if (error && error.name === "DevCommandCancelled") {
          return;
        }

        writeLog({
          ok: false,
          message: error.message
        });
      } finally {
        button.disabled = false;
        button.textContent = originalText;
      }
    }

    function openResetServerWarning() {
      let clickCount = 0;

      const overlay = document.createElement("div");
      overlay.className = "reset-server-overlay";
      overlay.innerHTML =
        "<section class=\"reset-server-modal\" role=\"dialog\" aria-modal=\"true\" aria-labelledby=\"reset-server-title\">" +
          "<h2 id=\"reset-server-title\">RESET SERVER</h2>" +
          "<p><strong>This is destructive.</strong> Clicking outside this warning cancels the action.</p>" +
          "<ul>" +
            "<li>Deletes all player accounts except Soryn.</li>" +
            "<li>Deletes every wardrobe item and restarts item IDs.</li>" +
            "<li>Resets currencies, work resources, work levels, Playdeck level, and Nami level.</li>" +
            "<li>Returns Nami care status to the initial starting values.</li>" +
          "</ul>" +
          "<div class=\"reset-server-actions\">" +
            "<button class=\"reset-server-confirm\" type=\"button\">Confirm Reset (0 / 3)</button>" +
            "<button class=\"reset-server-cancel\" type=\"button\">Cancel</button>" +
          "</div>" +
        "</section>";

      const confirmButton = overlay.querySelector(".reset-server-confirm");
      const cancelButton = overlay.querySelector(".reset-server-cancel");

      function closeModal() {
        document.removeEventListener("keydown", handleKeyDown);
        overlay.remove();
      }

      function handleKeyDown(event) {
        if (event.key === "Escape") {
          closeModal();
        }
      }

      overlay.addEventListener("click", (event) => {
        if (event.target === overlay) {
          closeModal();
        }
      });

      cancelButton.addEventListener("click", closeModal);

      confirmButton.addEventListener("click", async () => {
        clickCount += 1;

        if (clickCount < 3) {
          confirmButton.textContent = "Confirm Reset (" + clickCount + " / 3)";
          return;
        }

        confirmButton.disabled = true;
        cancelButton.disabled = true;
        confirmButton.textContent = "Resetting...";

        try {
          const response = await csrfFetch("/api/dev/reset-server", { method: "POST" });
          const text = await response.text();

          let responsePayload = text;
          try {
            responsePayload = text ? JSON.parse(text) : {};
          } catch {
            responsePayload = text;
          }

          if (!response.ok) {
            writeLog({
              ok: false,
              status: response.status,
              response: responsePayload
            });
            closeModal();
            return;
          }

          writeLog(responsePayload);
          closeModal();
        } catch (error) {
          writeLog({
            ok: false,
            message: error.message
          });
          closeModal();
        }
      });

      document.addEventListener("keydown", handleKeyDown);
      document.body.appendChild(overlay);
      confirmButton.focus();
    }
    if (lockButton) {
      lockButton.addEventListener("click", lockDevConsole);
    }

    if (resetServerButton) {
      resetServerButton.addEventListener("click", openResetServerWarning);
    }

    document.querySelectorAll("[data-endpoint]").forEach((button) => {
      button.addEventListener("click", () => {
        runDevCommand(button.dataset.endpoint, button);
      });
    });
  </script>
</body>
</html>`

type DevAccountAgeResponse struct {
	OK         bool   `json:"ok"`
	Display    string `json:"display"`
	ExactAge   string `json:"exactAge"`
	AgeSeconds int64  `json:"ageSeconds"`
	CreatedAt  string `json:"createdAt"`
	Now        string `json:"now"`
	Source     string `json:"source"`
}

func formatExactAccountAge(duration time.Duration) (string, int64) {
	totalSeconds := int64(duration / time.Second)
	if totalSeconds < 0 {
		totalSeconds = 0
	}

	days := totalSeconds / 86400
	hours := (totalSeconds % 86400) / 3600
	minutes := (totalSeconds % 3600) / 60
	seconds := totalSeconds % 60

	return fmt.Sprintf("%dd%dh%dm%ds", days, hours, minutes, seconds), totalSeconds
}

func (s *Server) HandleDevAccountAge(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	account, ok := s.AuthAccountFromRequest(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "login required")
		return
	}

	now := time.Now().UTC()
	createdAt := account.CreatedAt.UTC()
	exactAge, ageSeconds := formatExactAccountAge(now.Sub(createdAt))

	writeJSON(w, http.StatusOK, DevAccountAgeResponse{
		OK:         true,
		Display:    "Account Age: " + exactAge,
		ExactAge:   exactAge,
		AgeSeconds: ageSeconds,
		CreatedAt:  createdAt.Format(time.RFC3339),
		Now:        now.Format(time.RFC3339),
		Source:     "auth_accounts.created_at",
	})
}
func (s *Server) HandleDevConsolePage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	path := strings.TrimRight(r.URL.Path, "/")
	if path == "" {
		path = "/dev"
	}

	if path != "/dev" {
		http.NotFound(w, r)
		return
	}

	if !s.requestHasDevUnlock(r) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		_, _ = w.Write([]byte(devConsoleUnlockHTML))
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	_, _ = w.Write([]byte(devConsoleHTML))
}
