package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/ChaPerx64/dobby/apps/backend/internal/adapters/oas"
)

type contextKey string

const (
	userIDKey contextKey = "userID"
)

type introspectionResponse struct {
	Active bool   `json:"active"`
	Sub    string `json:"sub"`
}

type dobbySecurity struct {
	introspectionURL string
	clientID         string
	clientSecret     string
	httpClient       *http.Client
}

func (s *dobbySecurity) HandleBearerAuth(ctx context.Context, operationName oas.OperationName, t oas.BearerAuth) (context.Context, error) {
	if s.introspectionURL == "" {
		return ctx, fmt.Errorf("authorization not configured: introspection URL missing")
	}

	data := url.Values{}
	data.Set("token", t.Token)

	req, err := http.NewRequestWithContext(ctx, "POST", s.introspectionURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create introspection request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(s.clientID, s.clientSecret)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("introspection request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("introspection returned status %d: %s", resp.StatusCode, string(body))
	}

	var ir introspectionResponse
	if err := json.NewDecoder(resp.Body).Decode(&ir); err != nil {
		return nil, fmt.Errorf("failed to decode introspection response: %w", err)
	}


	if !ir.Active {
		return nil, fmt.Errorf("invalid token")
	}

	// Extract the subject (user ID) from the authorized context
	return context.WithValue(ctx, userIDKey, ir.Sub), nil
}

// GetUserID retrieves the user ID from the context if it exists.
func GetUserID(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(userIDKey).(string)
	return id, ok
}
