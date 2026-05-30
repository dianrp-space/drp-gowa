package chatwoot

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestChatwootClientSendsBothAuthHeaders(t *testing.T) {
	var receivedAPIAccessToken string
	var receivedLegacyToken string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedAPIAccessToken = r.Header.Get("Api-Access-Token")
		receivedLegacyToken = r.Header.Get("api_access_token")

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"payload":[]}`))
	}))
	defer server.Close()

	client := &Client{
		BaseURL:    server.URL,
		APIToken:   "test-token",
		AccountID:  123,
		InboxID:    456,
		HTTPClient: server.Client(),
	}

	_, err := client.FindContactByIdentifier("628123456789", true)
	if err != nil {
		t.Fatalf("FindContactByIdentifier returned error: %v", err)
	}

	if receivedAPIAccessToken != "test-token" {
		t.Fatalf("expected Api-Access-Token header to be set, got %q", receivedAPIAccessToken)
	}

	if receivedLegacyToken != "test-token" {
		t.Fatalf("expected api_access_token header to be set for backward compatibility, got %q", receivedLegacyToken)
	}
}
