package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// ListBetaGroups lists all TestFlight beta groups across all apps, with
// pagination support. The related app is included so callers can associate a
// group with its app without an extra request.
func (c *Client) ListBetaGroups(ctx context.Context, cursor string) ([]BetaGroup, *PaginationLinks, error) {
	url := BaseURL + "/v1/betaGroups"
	if cursor != "" {
		url = cursor // App Store Connect pagination uses full URLs in links.next
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create list beta groups request: %w", err)
	}

	if cursor == "" {
		q := req.URL.Query()
		q.Set("limit", "200")
		q.Set("include", "app")
		req.URL.RawQuery = q.Encode()
	}

	resp, err := c.Do(ctx, req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list beta groups: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, nil, fmt.Errorf("failed to list beta groups: HTTP %d %s", resp.StatusCode, resp.Status)
	}

	var result Response[BetaGroup]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, nil, fmt.Errorf("failed to decode list beta groups response: %w", err)
	}

	return result.Data, &result.Links, nil
}

// GetBetaGroup retrieves a single beta group by ID.
func (c *Client) GetBetaGroup(ctx context.Context, betaGroupID string) (*BetaGroup, error) {
	url := fmt.Sprintf("%s/v1/betaGroups/%s", BaseURL, betaGroupID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get beta group request: %w", err)
	}

	resp, err := c.Do(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get beta group: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("failed to get beta group: HTTP %d %s", resp.StatusCode, resp.Status)
	}

	var result SingleResponse[BetaGroup]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode get beta group response: %w", err)
	}

	return &result.Data, nil
}

// ListBetaTesters lists all TestFlight beta testers with pagination support.
func (c *Client) ListBetaTesters(ctx context.Context, cursor string) ([]BetaTester, *PaginationLinks, error) {
	url := BaseURL + "/v1/betaTesters"
	if cursor != "" {
		url = cursor
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create list beta testers request: %w", err)
	}

	if cursor == "" {
		q := req.URL.Query()
		q.Set("limit", "200")
		req.URL.RawQuery = q.Encode()
	}

	resp, err := c.Do(ctx, req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list beta testers: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, nil, fmt.Errorf("failed to list beta testers: HTTP %d %s", resp.StatusCode, resp.Status)
	}

	var result Response[BetaTester]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, nil, fmt.Errorf("failed to decode list beta testers response: %w", err)
	}

	return result.Data, &result.Links, nil
}

// ListBetaGroupTesters lists the beta testers that belong to a beta group.
func (c *Client) ListBetaGroupTesters(ctx context.Context, betaGroupID, cursor string) ([]BetaTester, *PaginationLinks, error) {
	url := fmt.Sprintf("%s/v1/betaGroups/%s/betaTesters", BaseURL, betaGroupID)
	if cursor != "" {
		url = cursor
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create list beta group testers request: %w", err)
	}

	if cursor == "" {
		q := req.URL.Query()
		q.Set("limit", "200")
		req.URL.RawQuery = q.Encode()
	}

	resp, err := c.Do(ctx, req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list beta group testers: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, nil, fmt.Errorf("failed to list beta group testers: HTTP %d %s", resp.StatusCode, resp.Status)
	}

	var result Response[BetaTester]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, nil, fmt.Errorf("failed to decode list beta group testers response: %w", err)
	}

	return result.Data, &result.Links, nil
}

// AddBetaTesterToGroup adds an existing beta tester to a beta group.
func (c *Client) AddBetaTesterToGroup(ctx context.Context, betaGroupID, betaTesterID string) error {
	payload := RelationshipRequest{
		Data: []RelationshipData{
			{Type: "betaTesters", ID: betaTesterID},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal add beta tester request: %w", err)
	}

	url := fmt.Sprintf("%s/v1/betaGroups/%s/relationships/betaTesters", BaseURL, betaGroupID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create add beta tester request: %w", err)
	}

	resp, err := c.Do(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to add beta tester to group: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("failed to add beta tester to group: HTTP %d %s", resp.StatusCode, resp.Status)
	}

	return nil
}

// RemoveBetaTesterFromGroup removes a beta tester from a beta group. The tester
// itself is not deleted and keeps any other group memberships.
func (c *Client) RemoveBetaTesterFromGroup(ctx context.Context, betaGroupID, betaTesterID string) error {
	payload := RelationshipRequest{
		Data: []RelationshipData{
			{Type: "betaTesters", ID: betaTesterID},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal remove beta tester request: %w", err)
	}

	url := fmt.Sprintf("%s/v1/betaGroups/%s/relationships/betaTesters", BaseURL, betaGroupID)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create remove beta tester request: %w", err)
	}

	resp, err := c.Do(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to remove beta tester from group: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && (resp.StatusCode < 200 || resp.StatusCode >= 300) {
		return fmt.Errorf("failed to remove beta tester from group: HTTP %d %s", resp.StatusCode, resp.Status)
	}

	return nil
}

// CreateBetaTester invites a new beta tester by email. If betaGroupIDs are
// provided the tester is added to those groups on creation, which is what
// triggers the TestFlight invitation email.
func (c *Client) CreateBetaTester(ctx context.Context, email, firstName, lastName string, betaGroupIDs []string) (*BetaTester, error) {
	request := BetaTesterCreateRequest{
		Data: BetaTesterCreateData{
			Type: "betaTesters",
			Attributes: BetaTesterCreateAttributes{
				Email:     email,
				FirstName: firstName,
				LastName:  lastName,
			},
		},
	}

	if len(betaGroupIDs) > 0 {
		relData := make([]RelationshipData, 0, len(betaGroupIDs))
		for _, groupID := range betaGroupIDs {
			relData = append(relData, RelationshipData{Type: "betaGroups", ID: groupID})
		}
		request.Data.Relationships = &BetaTesterRelationships{
			BetaGroups: &RelationshipRequest{Data: relData},
		}
	}

	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal create beta tester request: %w", err)
	}

	url := BaseURL + "/v1/betaTesters"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create create beta tester request: %w", err)
	}

	resp, err := c.Do(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create beta tester: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("failed to create beta tester: HTTP %d %s", resp.StatusCode, resp.Status)
	}

	var result SingleResponse[BetaTester]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode create beta tester response: %w", err)
	}

	return &result.Data, nil
}

// DeleteBetaTester removes a beta tester from all groups and apps entirely.
func (c *Client) DeleteBetaTester(ctx context.Context, betaTesterID string) error {
	url := fmt.Sprintf("%s/v1/betaTesters/%s", BaseURL, betaTesterID)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create delete beta tester request: %w", err)
	}

	resp, err := c.Do(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to delete beta tester: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && (resp.StatusCode < 200 || resp.StatusCode >= 300) {
		return fmt.Errorf("failed to delete beta tester: HTTP %d %s", resp.StatusCode, resp.Status)
	}

	return nil
}

// FindBetaTesterByEmail looks up a beta tester by their email address. It
// returns nil without an error when no tester matches.
func (c *Client) FindBetaTesterByEmail(ctx context.Context, email string) (*BetaTester, error) {
	url := BaseURL + "/v1/betaTesters"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create find beta tester request: %w", err)
	}

	q := req.URL.Query()
	q.Set("filter[email]", email)
	q.Set("limit", "1")
	req.URL.RawQuery = q.Encode()

	resp, err := c.Do(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to find beta tester: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("failed to find beta tester: HTTP %d %s", resp.StatusCode, resp.Status)
	}

	var result Response[BetaTester]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode find beta tester response: %w", err)
	}

	if len(result.Data) == 0 {
		return nil, nil
	}

	return &result.Data[0], nil
}
