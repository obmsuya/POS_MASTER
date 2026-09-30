package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type fakeBackend struct {
	server        *httptest.Server
	refreshCalls  int
	historyCalls  int
	historyTokens []string
	refreshBodies []string
	historyStatus func(call int) int
}

func newFakeBackend(t *testing.T) *fakeBackend {
	backend := &fakeBackend{}
	mux := http.NewServeMux()
	mux.HandleFunc("/auth/token/refresh/", func(w http.ResponseWriter, r *http.Request) {
		backend.refreshCalls++
		var body map[string]string
		_ = json.NewDecoder(r.Body).Decode(&body)
		backend.refreshBodies = append(backend.refreshBodies, body["refresh"])
		_ = json.NewEncoder(w).Encode(map[string]string{"access": "access-2", "refresh": "refresh-2"})
	})
	mux.HandleFunc("/payments/balce/admin/payments/", func(w http.ResponseWriter, r *http.Request) {
		backend.historyCalls++
		backend.historyTokens = append(backend.historyTokens, r.Header.Get("Authorization"))
		w.WriteHeader(backend.historyStatus(backend.historyCalls))
		_, _ = w.Write([]byte(`{"data":[]}`))
	})
	backend.server = httptest.NewServer(mux)
	previousBaseURL := BaseURL
	BaseURL = backend.server.URL
	t.Cleanup(func() {
		BaseURL = previousBaseURL
		backend.server.Close()
	})
	return backend
}

func TestRefreshStoresRotatedRefreshToken(t *testing.T) {
	backend := newFakeBackend(t)
	client := New()
	client.RestoreSession("access-1", "refresh-1")

	if err := client.RefreshAccessToken(); err != nil {
		t.Fatalf("refresh failed: %v", err)
	}
	if client.AccessToken != "access-2" || client.RefreshToken != "refresh-2" {
		t.Fatalf("got access=%q refresh=%q, want access-2/refresh-2", client.AccessToken, client.RefreshToken)
	}
	if backend.refreshBodies[0] != "refresh-1" {
		t.Fatalf("refresh sent %q, want refresh-1", backend.refreshBodies[0])
	}
}

func TestUnauthorizedRefreshesOnceAndRetries(t *testing.T) {
	backend := newFakeBackend(t)
	backend.historyStatus = func(call int) int {
		if call == 1 {
			return http.StatusUnauthorized
		}
		return http.StatusOK
	}
	client := New()
	client.RestoreSession("access-1", "refresh-1")
	var savedAccess, savedRefresh string
	client.OnTokensRefreshed = func(access, refresh string) { savedAccess, savedRefresh = access, refresh }

	if _, err := client.PaymentHistory(""); err != nil {
		t.Fatalf("history failed: %v", err)
	}
	if backend.refreshCalls != 1 || backend.historyCalls != 2 {
		t.Fatalf("refresh calls=%d history calls=%d, want 1 and 2", backend.refreshCalls, backend.historyCalls)
	}
	if backend.historyTokens[1] != "Bearer access-2" {
		t.Fatalf("retry sent %q, want Bearer access-2", backend.historyTokens[1])
	}
	if savedAccess != "access-2" || savedRefresh != "refresh-2" {
		t.Fatalf("callback got access=%q refresh=%q", savedAccess, savedRefresh)
	}
}

func TestUnauthorizedAfterRetryReturnsError(t *testing.T) {
	backend := newFakeBackend(t)
	backend.historyStatus = func(int) int { return http.StatusUnauthorized }
	client := New()
	client.RestoreSession("access-1", "refresh-1")

	_, err := client.PaymentHistory("")
	var apiErr *APIError
	if !errors.As(err, &apiErr) || apiErr.StatusCode != http.StatusUnauthorized {
		t.Fatalf("got %v, want a 401 APIError", err)
	}
	if backend.refreshCalls != 1 || backend.historyCalls != 2 {
		t.Fatalf("refresh calls=%d history calls=%d, want 1 and 2", backend.refreshCalls, backend.historyCalls)
	}
}
