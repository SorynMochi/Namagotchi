package server

import (
	"net/http"
	"strings"
)

const devConsoleHTML = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Namigotchi Dev Console</title>
  <style>
    body {
      margin: 0;
      min-height: 100vh;
      font-family: system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
      background: #120c18;
      color: #fcefff;
    }

    main {
      max-width: 960px;
      margin: 0 auto;
      padding: 32px 18px 48px;
    }

    a {
      color: #ff8fc7;
      font-weight: 800;
    }

    .grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
      gap: 14px;
      margin-top: 22px;
    }

    .card {
      padding: 16px;
      border: 1px solid rgba(255,255,255,0.16);
      border-radius: 16px;
      background: #21162e;
    }

    button {
      width: 100%;
      border: 0;
      border-radius: 999px;
      padding: 11px 14px;
      cursor: pointer;
      background: #ff8fc7;
      color: #1f1022;
      font-weight: 900;
    }

    button.danger {
      background: #ff6b6b;
    }

    button.secondary {
      margin-bottom: 18px;
      color: #fcefff;
      background: rgba(255,255,255,0.12);
      border: 1px solid rgba(255,255,255,0.18);
    }

    pre {
      min-height: 220px;
      margin-top: 18px;
      padding: 14px;
      overflow: auto;
      white-space: pre-wrap;
      border-radius: 14px;
      border: 1px solid rgba(255,255,255,0.16);
      background: rgba(0,0,0,0.35);
    }
  </style>
</head>
<body>
  <main>
    <p><a href="/">Back to game</a></p>
    <button id="dev-lock-button" class="secondary" type="button">Lock Dev Console</button>
    <h1>Namigotchi Dev Console</h1>
    <p>This page is served only after the backend verifies your dev access.</p>

    <section class="grid">
      <div class="card">
        <h2>Setup</h2>
        <button data-endpoint="/api/dev/seed-player">Seed Player</button>
      </div>

      <div class="card">
        <h2>Playdeck</h2>
        <button data-endpoint="/api/dev/force-tick">Force Tick</button>
      </div>

      <div class="card">
        <h2>Playdeck</h2>
        <button class="danger" data-endpoint="/api/dev/reset-playdeck-streak">Reset Streak</button>
      </div>

      <div class="card">
        <h2>Wardrobe</h2>
        <button data-endpoint="/api/dev/spawn-wardrobe-item">Spawn Random Item</button>
      </div>
      <div class="card">
  <h2>Audit</h2>
  <button data-endpoint="/api/dev/audit-logs" data-method="GET">Refresh Audit Logs</button>
</div>
      <div class="card">
        <h2>Security</h2>
        <button data-endpoint="/api/dev/security-events" data-method="GET">Refresh Security Events</button>
      </div>
    </section>

    <pre id="dev-log">Ready.</pre>
  </main>

  <script>
    const log = document.querySelector("#dev-log");
    const lockButton = document.querySelector("#dev-lock-button");
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

    async function runDevCommand(endpoint, button) {
      const originalText = button.textContent;

      button.disabled = true;
      button.textContent = "Running...";

      try {
        const method = button.dataset.method || "POST";
        const response = await csrfFetch(endpoint, { method });
        const text = await response.text();

        let payload = text;
        try {
          payload = text ? JSON.parse(text) : {};
        } catch {
          payload = text;
        }

        if (!response.ok) {
          writeLog({
            ok: false,
            status: response.status,
            response: payload
          });
          return;
        }

        writeLog(payload);
      } catch (error) {
        writeLog({
          ok: false,
          message: error.message
        });
      } finally {
        button.disabled = false;
        button.textContent = originalText;
      }
    }

    if (lockButton) {
      lockButton.addEventListener("click", lockDevConsole);
    }

    document.querySelectorAll("[data-endpoint]").forEach((button) => {
      button.addEventListener("click", () => {
        runDevCommand(button.dataset.endpoint, button);
      });
    });
  </script>
</body>
</html>`

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
