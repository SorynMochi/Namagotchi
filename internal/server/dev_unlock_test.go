package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestDevUnlockTokenVerification(t *testing.T) {
	secret := "test-dev-unlock-secret"
	accountID := int64(42)
	expiresAt := time.Now().Add(5 * time.Minute)

	token := makeDevUnlockToken(accountID, expiresAt, secret)

	if !verifyDevUnlockToken(token, accountID, secret, time.Now()) {
		t.Fatal("expected valid dev unlock token to verify")
	}

	if verifyDevUnlockToken(token, accountID+1, secret, time.Now()) {
		t.Fatal("expected token for different account to fail")
	}

	if verifyDevUnlockToken(token, accountID, "wrong-secret", time.Now()) {
		t.Fatal("expected token signed with different secret to fail")
	}

	if verifyDevUnlockToken(token, accountID, secret, expiresAt.Add(time.Second)) {
		t.Fatal("expected expired token to fail")
	}

	if verifyDevUnlockToken("not-a-valid-token", accountID, secret, time.Now()) {
		t.Fatal("expected malformed token to fail")
	}
}

func TestHandleDevLockClearsUnlockCookie(t *testing.T) {
	s := &Server{}

	request := httptest.NewRequest(http.MethodPost, "/api/dev/lock", nil)
	recorder := httptest.NewRecorder()

	s.HandleDevLock(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	response := recorder.Result()
	defer response.Body.Close()

	var lockCookie *http.Cookie
	for _, cookie := range response.Cookies() {
		if cookie.Name == devUnlockCookieName {
			lockCookie = cookie
			break
		}
	}

	if lockCookie == nil {
		t.Fatal("expected dev unlock cookie to be cleared")
	}

	if lockCookie.MaxAge >= 0 {
		t.Fatalf("expected dev unlock cookie MaxAge to be negative, got %d", lockCookie.MaxAge)
	}

	if !strings.Contains(recorder.Body.String(), "Dev console locked.") {
		t.Fatalf("expected lock response message, got %q", recorder.Body.String())
	}
}

func TestHandleDevLockRejectsGet(t *testing.T) {
	s := &Server{}

	request := httptest.NewRequest(http.MethodGet, "/api/dev/lock", nil)
	recorder := httptest.NewRecorder()

	s.HandleDevLock(recorder, request)

	if recorder.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status %d, got %d", http.StatusMethodNotAllowed, recorder.Code)
	}
}
