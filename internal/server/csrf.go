package server

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"
)

const csrfCookieName = "namigotchi_csrf"
const csrfHeaderName = "X-CSRF-Token"

const csrfTokenTTL = 12 * time.Hour

type CSRFTokenResponse struct {
	OK        bool   `json:"ok"`
	CSRFToken string `json:"csrfToken"`
}

func (s *Server) HandleCSRFToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.recordSecurityEvent(r, http.StatusMethodNotAllowed, "csrf", "csrf token method not allowed")
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	token := ensureCSRFToken(w, r)

	writeJSON(w, http.StatusOK, CSRFTokenResponse{
		OK:        true,
		CSRFToken: token,
	})
}

func (s *Server) requireCSRF(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if isSafeCSRFMethod(r.Method) {
			next(w, r)
			return
		}

		cookie, err := r.Cookie(csrfCookieName)
		if err != nil {
			s.recordSecurityEvent(r, http.StatusForbidden, "csrf", "csrf token required")
			writeError(w, http.StatusForbidden, "csrf token required")
			return
		}

		cookieToken := cookie.Value
		headerToken := r.Header.Get(csrfHeaderName)

		if cookieToken == "" || headerToken == "" || !constantTimeStringEqual(cookieToken, headerToken) {
			s.recordSecurityEvent(r, http.StatusForbidden, "csrf", "csrf token invalid")
			writeError(w, http.StatusForbidden, "csrf token invalid")
			return
		}

		next(w, r)
	}
}

func ensureCSRFToken(w http.ResponseWriter, r *http.Request) string {
	if cookie, err := r.Cookie(csrfCookieName); err == nil && cookie.Value != "" {
		return cookie.Value
	}

	token := generateCSRFToken()
	setCSRFCookie(w, r, token)

	return token
}

func setCSRFCookie(w http.ResponseWriter, r *http.Request, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     csrfCookieName,
		Value:    token,
		Path:     "/",
		MaxAge:   int(csrfTokenTTL.Seconds()),
		Expires:  time.Now().Add(csrfTokenTTL),
		HttpOnly: false,
		SameSite: http.SameSiteLaxMode,
		Secure:   authCookieShouldBeSecure(r),
	})
}

func generateCSRFToken() string {
	bytes := make([]byte, 32)

	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}

	return base64.RawURLEncoding.EncodeToString(bytes)
}

func isSafeCSRFMethod(method string) bool {
	return method == http.MethodGet ||
		method == http.MethodHead ||
		method == http.MethodOptions
}
