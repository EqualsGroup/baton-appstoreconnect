package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// ListApps lists all apps in App Store Connect with pagination support.
func (c *Client) ListApps(ctx context.Context, cursor string) ([]App, *PaginationLinks, error) {
	url := BaseURL + "/v1/apps"
	if cursor != "" {
		url = cursor // App Store Connect pagination uses full URLs in links.next
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create list apps request: %w", err)
	}

	// Request 100 apps per page (max allowed).
	if cursor == "" {
		q := req.URL.Query()
		q.Set("limit", "100")
		req.URL.RawQuery = q.Encode()
	}

	resp, err := c.Do(ctx, req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list apps: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, nil, fmt.Errorf("failed to list apps: HTTP %d %s", resp.StatusCode, resp.Status)
	}

	var result Response[App]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, nil, fmt.Errorf("failed to decode list apps response: %w", err)
	}

	return result.Data, &result.Links, nil
}

// GetApp retrieves a single app by ID.
func (c *Client) GetApp(ctx context.Context, appID string) (*App, error) {
	url := fmt.Sprintf("%s/v1/apps/%s", BaseURL, appID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get app request: %w", err)
	}

	resp, err := c.Do(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get app: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("failed to get app: HTTP %d %s", resp.StatusCode, resp.Status)
	}

	var result SingleResponse[App]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode get app response: %w", err)
	}

	return &result.Data, nil
}
