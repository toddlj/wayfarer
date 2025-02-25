package telegram

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSendMessage_Success(t *testing.T) {
	// given
	requestSent := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		expectedPath := fmt.Sprintf("/bot%s/sendMessage", "FAKE_TOKEN")
		if r.URL.Path != expectedPath {
			t.Errorf("expected URL path %q, got %q", expectedPath, r.URL.Path)
		}
		requestSent = true

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok": true}`))
	}))
	defer ts.Close()

	client := NewClient(ts.URL, "FAKE_TOKEN")

	// when
	err := client.SendMessage(12345, "Hello, world!")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	// then
	if !requestSent {
		t.Fatal("expected request to be sent, but it wasn't")
	}
}

func TestSendMessage_Non200Response(t *testing.T) {
	// given
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer ts.Close()

	client := NewClient(ts.URL, "FAKE_TOKEN")

	// when
	err := client.SendMessage(12345, "Hello, world!")
	if err == nil {
		t.Fatal("expected an error due to non-200 response, got nil")
	}

	// then
	expectedErr := "bad status code received: 400"
	if err.Error() != expectedErr {
		t.Errorf("expected error %q, got %q", expectedErr, err.Error())
	}
}
