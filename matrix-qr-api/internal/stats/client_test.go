package stats

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"matrix-qr-api/internal/matrix"
)

func TestCalculateSendsMatricesAndToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/statistics" {
			t.Errorf("unexpected path %s", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer abc" {
			t.Errorf("unexpected authorization header %q", got)
		}
		var body struct {
			Matrices []NamedMatrix `json:"matrices"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || len(body.Matrices) != 1 {
			t.Errorf("unexpected body: %v %+v", err, body)
		}
		_ = json.NewEncoder(w).Encode(Statistics{Max: 4, Min: 1, Sum: 10, Average: 2.5, Count: 4})
	}))
	defer server.Close()

	client := NewHTTPClient(server.URL, time.Second)
	out, err := client.Calculate(context.Background(), "abc", []NamedMatrix{{Name: "A", Data: matrix.Matrix{{1, 2}, {3, 4}}}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Max != 4 || out.Sum != 10 {
		t.Fatalf("unexpected statistics: %+v", out)
	}
}

func TestCalculatePropagatesUpstreamErrors(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	client := NewHTTPClient(server.URL, time.Second)
	if _, err := client.Calculate(context.Background(), "abc", nil); err == nil {
		t.Fatal("expected error")
	}
}
