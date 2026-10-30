package stats

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"matrix-qr-api/internal/matrix"
)

type NamedMatrix struct {
	Name string        `json:"name"`
	Data matrix.Matrix `json:"data"`
}

type MatrixSummary struct {
	Name       string `json:"name"`
	Rows       int    `json:"rows"`
	Cols       int    `json:"cols"`
	IsDiagonal bool   `json:"isDiagonal"`
}

type Statistics struct {
	Max         float64         `json:"max"`
	Min         float64         `json:"min"`
	Average     float64         `json:"average"`
	Sum         float64         `json:"sum"`
	Count       int             `json:"count"`
	AnyDiagonal bool            `json:"anyDiagonal"`
	Matrices    []MatrixSummary `json:"matrices"`
}

type Calculator interface {
	Calculate(ctx context.Context, token string, matrices []NamedMatrix) (*Statistics, error)
}

type HTTPClient struct {
	baseURL string
	http    *http.Client
}

func NewHTTPClient(baseURL string, timeout time.Duration) *HTTPClient {
	return &HTTPClient{
		baseURL: baseURL,
		http:    &http.Client{Timeout: timeout},
	}
}

func (c *HTTPClient) Calculate(ctx context.Context, token string, matrices []NamedMatrix) (*Statistics, error) {
	body, err := json.Marshal(map[string]any{"matrices": matrices})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/api/statistics", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("statistics service unreachable: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("statistics service responded with status %d", res.StatusCode)
	}

	var out Statistics
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("invalid statistics response: %w", err)
	}
	return &out, nil
}
