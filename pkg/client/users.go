package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// ListUsers lists all users in App Store Connect with pagination support.
func (c *Client) ListUsers(ctx context.Context, cursor string) ([]User, *PaginationLinks, error) {
	url := BaseURL + "/v1/users"
	if cursor != "" {
		url = cursor // App Store Connect pagination uses full URLs in links.next
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create list users request: %w", err)
	}

	// Request 100 users per page (max allowed).
	if cursor == "" {
		q := req.URL.Query()
		q.Set("limit", "100")
		req.URL.RawQuery = q.Encode()
	}

	resp, err := c.Do(ctx, req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list users: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, nil, fmt.Errorf("failed to list users: HTTP %d %s", resp.StatusCode, resp.Status)
	}

	var result Response[User]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, nil, fmt.Errorf("failed to decode list users response: %w", err)
	}

	return result.Data, &result.Links, nil
}

// GetUser retrieves a single user by ID.
func (c *Client) GetUser(ctx context.Context, userID string) (*User, error) {
	url := fmt.Sprintf("%s/v1/users/%s", BaseURL, userID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get user request: %w", err)
	}

	resp, err := c.Do(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("failed to get user: HTTP %d %s", resp.StatusCode, resp.Status)
	}

	var result SingleResponse[User]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode get user response: %w", err)
	}

	return &result.Data, nil
}

// DeleteUser removes a user from App Store Connect.
func (c *Client) DeleteUser(ctx context.Context, userID string) error {
	url := fmt.Sprintf("%s/v1/users/%s", BaseURL, userID)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create delete user request: %w", err)
	}

	resp, err := c.Do(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && (resp.StatusCode < 200 || resp.StatusCode >= 300) {
		return fmt.Errorf("failed to delete user: HTTP %d %s", resp.StatusCode, resp.Status)
	}

	return nil
}

// InviteUser creates a user invitation in App Store Connect.
func (c *Client) InviteUser(ctx context.Context, email, firstName, lastName string, roles []string, allAppsVisible bool, visibleAppIDs []string) error {
	invitation := UserInvitationCreateRequest{
		Data: UserInvitationCreateData{
			Type: "userInvitations",
			Attributes: UserInvitationCreateAttributes{
				Email:          email,
				FirstName:      firstName,
				LastName:       lastName,
				Roles:          roles,
				AllAppsVisible: allAppsVisible,
			},
		},
	}

	// Add visible apps relationship if specific apps are specified.
	if !allAppsVisible && len(visibleAppIDs) > 0 {
		relData := make([]RelationshipData, 0, len(visibleAppIDs))
		for _, appID := range visibleAppIDs {
			relData = append(relData, RelationshipData{
				Type: "apps",
				ID:   appID,
			})
		}
		invitation.Data.Relationships = &UserInvitationRelationships{
			VisibleApps: &VisibleAppsRelationship{
				Data: relData,
			},
		}
	}

	body, err := json.Marshal(invitation)
	if err != nil {
		return fmt.Errorf("failed to marshal user invitation: %w", err)
	}

	url := BaseURL + "/v1/userInvitations"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create invite user request: %w", err)
	}

	resp, err := c.Do(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to invite user: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("failed to invite user: HTTP %d %s", resp.StatusCode, resp.Status)
	}

	return nil
}

// ListUserVisibleApps lists the apps visible to a specific user.
func (c *Client) ListUserVisibleApps(ctx context.Context, userID, cursor string) ([]App, *PaginationLinks, error) {
	url := fmt.Sprintf("%s/v1/users/%s/visibleApps", BaseURL, userID)
	if cursor != "" {
		url = cursor
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create list user visible apps request: %w", err)
	}

	if cursor == "" {
		q := req.URL.Query()
		q.Set("limit", "100")
		req.URL.RawQuery = q.Encode()
	}

	resp, err := c.Do(ctx, req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list user visible apps: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, nil, fmt.Errorf("failed to list user visible apps: HTTP %d %s", resp.StatusCode, resp.Status)
	}

	var result UserVisibleAppsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, nil, fmt.Errorf("failed to decode user visible apps response: %w", err)
	}

	return result.Data, &result.Links, nil
}
