// Package statsclient: cliente HTTP de la API de estadísticas.
package statsclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rbravo/interseguro-challenge/api-go/internal/matrix"
)

type NamedMatrix struct {
	Name string        `json:"name"`
	Data matrix.Matrix `json:"data"`
}

type Request struct {
	Matrices []NamedMatrix `json:"matrices"`
}

type Client struct {
	baseURL string
	http    *http.Client
}

// New crea un cliente con timeout explícito para no dejar requests colgados.
func New(baseURL string, timeout time.Duration) *Client {
	return &Client{
		baseURL: baseURL,
		http:    &http.Client{Timeout: timeout},
	}
}

// Devuelve un mapa genérico: Go solo orquesta, no valida el contrato.
func (c *Client) Statistics(ctx context.Context, token string, matrices []NamedMatrix) (map[string]any, error) {
	payload, err := json.Marshal(Request{Matrices: matrices})
	if err != nil {
		return nil, fmt.Errorf("no se pudo serializar el payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("no se pudo construir la petición: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("la API de estadísticas no respondió: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("no se pudo leer la respuesta: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("la API de estadísticas devolvió %d: %s", resp.StatusCode, string(body))
	}

	var out map[string]any
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, fmt.Errorf("respuesta inválida de la API de estadísticas: %w", err)
	}
	return out, nil
}
