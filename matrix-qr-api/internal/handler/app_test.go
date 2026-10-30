package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"matrix-qr-api/internal/auth"
	"matrix-qr-api/internal/stats"
)

var secret = []byte("test-secret-with-enough-length")

type fakeCalculator struct {
	token    string
	matrices []stats.NamedMatrix
	err      error
}

func (f *fakeCalculator) Calculate(_ context.Context, token string, matrices []stats.NamedMatrix) (*stats.Statistics, error) {
	f.token = token
	f.matrices = matrices
	if f.err != nil {
		return nil, f.err
	}
	return &stats.Statistics{Max: 1, Min: 0, Count: len(matrices)}, nil
}

func do(t *testing.T, calc stats.Calculator, method, path, token string, body any) (*http.Response, map[string]any) {
	t.Helper()
	app := New(Options{
		JWTSecret:     secret,
		TokenTTL:      time.Minute,
		Username:      "admin",
		Password:      "secret",
		AllowedOrigin: "*",
		Stats:         calc,
	})

	var reader io.Reader
	if body != nil {
		raw, _ := json.Marshal(body)
		reader = bytes.NewReader(raw)
	}
	req := httptest.NewRequest(method, path, reader)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	var out map[string]any
	_ = json.NewDecoder(res.Body).Decode(&out)
	return res, out
}

func validToken(t *testing.T) string {
	t.Helper()
	token, err := auth.Issue(secret, "admin", time.Minute)
	if err != nil {
		t.Fatalf("issue token: %v", err)
	}
	return token
}

func TestLogin(t *testing.T) {
	res, body := do(t, &fakeCalculator{}, http.MethodPost, "/api/auth/login", "", map[string]string{"username": "admin", "password": "secret"})
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status %d", res.StatusCode)
	}
	if body["token"] == "" || body["token"] == nil {
		t.Fatal("expected token in response")
	}

	res, _ = do(t, &fakeCalculator{}, http.MethodPost, "/api/auth/login", "", map[string]string{"username": "admin", "password": "wrong"})
	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("status %d, want 401", res.StatusCode)
	}
}

func TestQRRequiresToken(t *testing.T) {
	res, _ := do(t, &fakeCalculator{}, http.MethodPost, "/api/qr", "", map[string]any{"matrix": [][]float64{{1}}})
	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("status %d, want 401", res.StatusCode)
	}

	res, _ = do(t, &fakeCalculator{}, http.MethodPost, "/api/qr", "not-a-token", map[string]any{"matrix": [][]float64{{1}}})
	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("status %d, want 401", res.StatusCode)
	}
}

func TestQRReturnsFactorsAndStatistics(t *testing.T) {
	calc := &fakeCalculator{}
	token := validToken(t)

	res, body := do(t, calc, http.MethodPost, "/api/qr", token, map[string]any{"matrix": [][]float64{{1, 2}, {3, 4}, {5, 6}}})
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status %d: %v", res.StatusCode, body)
	}
	if len(body["q"].([]any)) != 3 || len(body["r"].([]any)) != 3 {
		t.Fatalf("unexpected shapes: %v", body)
	}
	if body["statistics"] == nil {
		t.Fatal("expected statistics")
	}
	if calc.token != token {
		t.Fatal("token was not forwarded to statistics service")
	}
	if len(calc.matrices) != 2 || calc.matrices[0].Name != "Q" || calc.matrices[1].Name != "R" {
		t.Fatalf("unexpected matrices sent: %+v", calc.matrices)
	}
}

func TestQRValidation(t *testing.T) {
	token := validToken(t)
	cases := map[string]any{
		"jagged":  map[string]any{"matrix": [][]float64{{1, 2}, {3}}},
		"empty":   map[string]any{"matrix": [][]float64{}},
		"missing": map[string]any{},
		"text":    map[string]any{"matrix": [][]string{{"a"}}},
	}
	for name, body := range cases {
		t.Run(name, func(t *testing.T) {
			res, _ := do(t, &fakeCalculator{}, http.MethodPost, "/api/qr", token, body)
			if res.StatusCode != http.StatusBadRequest {
				t.Fatalf("status %d, want 400", res.StatusCode)
			}
		})
	}
}

func TestQRStatisticsFailure(t *testing.T) {
	res, _ := do(t, &fakeCalculator{err: errors.New("down")}, http.MethodPost, "/api/qr", validToken(t), map[string]any{"matrix": [][]float64{{1}}})
	if res.StatusCode != http.StatusBadGateway {
		t.Fatalf("status %d, want 502", res.StatusCode)
	}
}
