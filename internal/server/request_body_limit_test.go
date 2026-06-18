package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRequestBodyLimitBlocksLargeAPIPost(t *testing.T) {
	server := &Server{}

	handler := server.withRequestBodyLimit(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	body := strings.NewReader(strings.Repeat("x", int(apiRequestBodyLimitBytes)+1))
	request := httptest.NewRequest(http.MethodPost, "/api/auth/login", body)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	response := recorder.Result()
	defer response.Body.Close()

	if response.StatusCode != http.StatusRequestEntityTooLarge {
		t.Fatalf("expected status %d, got %d", http.StatusRequestEntityTooLarge, response.StatusCode)
	}
}

func TestRequestBodyLimitAllowsSmallAPIPost(t *testing.T) {
	server := &Server{}

	called := false
	handler := server.withRequestBodyLimit(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusNoContent)
	}))

	request := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(`{"ok":true}`))
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	response := recorder.Result()
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, response.StatusCode)
	}

	if !called {
		t.Fatal("expected wrapped handler to be called")
	}
}

func TestRequestBodyLimitDoesNotLimitStaticPost(t *testing.T) {
	server := &Server{}

	called := false
	handler := server.withRequestBodyLimit(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusNoContent)
	}))

	body := strings.NewReader(strings.Repeat("x", int(apiRequestBodyLimitBytes)+1))
	request := httptest.NewRequest(http.MethodPost, "/some-static-path", body)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	response := recorder.Result()
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, response.StatusCode)
	}

	if !called {
		t.Fatal("expected wrapped handler to be called")
	}
}

func TestShouldLimitRequestBody(t *testing.T) {
	tests := []struct {
		method string
		path   string
		want   bool
	}{
		{method: http.MethodPost, path: "/api/player/care", want: true},
		{method: http.MethodPut, path: "/api/player/care", want: true},
		{method: http.MethodPatch, path: "/api/player/care", want: true},
		{method: http.MethodGet, path: "/api/player/status", want: false},
		{method: http.MethodPost, path: "/app.js", want: false},
	}

	for _, test := range tests {
		request := httptest.NewRequest(test.method, test.path, strings.NewReader("{}"))

		got := shouldLimitRequestBody(request)
		if got != test.want {
			t.Fatalf("method=%s path=%s expected %v, got %v", test.method, test.path, test.want, got)
		}
	}
}
