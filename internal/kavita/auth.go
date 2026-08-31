package kavita

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
)

type AuthTransport struct {
	baseURL    string
	apiKey     string
	pluginName string
	base       http.RoundTripper

	mu         sync.RWMutex
	userApiKey string
}

func NewAuthTransport(baseURL, apiKey, pluginName string, base http.RoundTripper) *AuthTransport {
	if base == nil {
		base = http.DefaultTransport
	}
	return &AuthTransport{
		baseURL:    baseURL,
		apiKey:     apiKey,
		pluginName: pluginName,
		base:       base,
	}
}

// RoundTrip executes the request, injecting the apiKey and handling 401 retries.
func (t *AuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Skip auth injection for the authenticate endpoint itself
	if req.URL.Path == "/api/Account" {
		return t.base.RoundTrip(req)
	}

	reqCopy := cloneRequest(req)
	token := t.GetUserAPIKey()
	if token != "" {
		reqCopy.Header.Set("x-api-key", token)
	}

	resp, err := t.base.RoundTrip(reqCopy)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusUnauthorized {
		return resp, nil
	}

	// TODO: Implement reactive authentication on 401 responses. This is currently commented
	// out to avoid potential infinite loops or deadlocks during re-authentication.
	// // Handle 401: Drain response body before re-authenticating
	// _, _ = io.Copy(io.Discard, resp.Body)
	// _ = resp.Body.Close()

	// // Re-authenticate (protected by lock inside Authenticate)
	// if err := t.Authenticate(req.Context()); err != nil {
	// 	return nil, fmt.Errorf("reactive auth failed: %w", err)
	// }

	// // Retry original request with new token
	// retryReq := cloneRequest(req)
	// retryReq.Header.Set("x-api-key", t.GetUserAPIKey())

	// return t.base.RoundTrip(retryReq)
	return resp, nil
}

func (t *AuthTransport) Authenticate(ctx context.Context) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	endpoint := fmt.Sprintf("%s/api/Account?apiKey=%s",
		t.baseURL,
		url.QueryEscape(t.apiKey),
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return fmt.Errorf("create auth request: %w", err)
	}

	resp, err := t.base.RoundTrip(req)
	if err != nil {
		return fmt.Errorf("execute auth request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("kavita auth failed with status %d: %s", resp.StatusCode, string(body))
	}

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return fmt.Errorf("decode auth response: %w", err)
	}

	t.userApiKey = authResp.ApiKey

	return nil
}

func (t *AuthTransport) GetUserAPIKey() string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	if t.userApiKey != "" {
		return t.userApiKey
	}
	return t.apiKey
}

func cloneRequest(req *http.Request) *http.Request {
	r := req.Clone(req.Context())
	r.Header = make(http.Header, len(req.Header))
	for k, s := range req.Header {
		r.Header[k] = append([]string(nil), s...)
	}

	if req.Body != nil && req.GetBody != nil {
		body, err := req.GetBody()
		if err == nil {
			r.Body = body
		}
	} else if req.Body != nil {
		buf, err := io.ReadAll(req.Body)
		if err == nil {
			req.Body = io.NopCloser(bytes.NewReader(buf))
			r.Body = io.NopCloser(bytes.NewReader(buf))
		}
	}

	return r
}
