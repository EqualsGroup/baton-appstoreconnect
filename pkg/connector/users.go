package connector

import (
	"context"
	"fmt"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	resourceSdk "github.com/conductorone/baton-sdk/pkg/types/resource"
	"github.com/conductorone/baton-appstoreconnect/pkg/client"
)

type userBuilder struct {
	client *client.Client
}

func (u *userBuilder) ResourceType(_ context.Context) *v2.ResourceType {
	return userResourceType
}

func newUserResource(user client.User) (*v2.Resource, error) {
	displayName := fmt.Sprintf("%s %s", user.Attributes.FirstName, user.Attributes.LastName)
	if displayName == " " {
		displayName = user.Attributes.Email
	}

	profile := map[string]interface{}{
		"username":        user.Attributes.Username,
		"first_name":      user.Attributes.FirstName,
		"last_name":       user.Attributes.LastName,
		"roles":           user.Attributes.Roles,
		"all_apps_visible": user.Attributes.AllAppsVisible,
	}

	return resourceSdk.NewUserResource(
		displayName,
		userResourceType,
		user.ID,
		[]resourceSdk.UserTraitOption{
			resourceSdk.WithEmail(user.Attributes.Email, true),
			resourceSdk.WithUserProfile(profile),
		},
	)
}

// List returns all users from App Store Connect.
func (u *userBuilder) List(ctx context.Context, _ *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	var cursor string
	if pToken != nil {
		cursor = pToken.Token
	}

	users, links, err := u.client.ListUsers(ctx, cursor)
	if err != nil {
		return nil, "", nil, fmt.Errorf("baton-appstoreconnect: failed to list users: %w", err)
	}

	ret := make([]*v2.Resource, 0, len(users))
	for _, user := range users {
		resource, err := newUserResource(user)
		if err != nil {
			return nil, "", nil, fmt.Errorf("baton-appstoreconnect: failed to create user resource for %s: %w", user.ID, err)
		}
		ret = append(ret, resource)
	}

	nextCursor := ""
	if links != nil && links.Next != "" {
		nextCursor = links.Next
	}

	return ret, nextCursor, nil, nil
}

// Entitlements returns an empty slice for users (users don't own entitlements).
func (u *userBuilder) Entitlements(_ context.Context, _ *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

// Grants returns an empty slice for users.
func (u *userBuilder) Grants(_ context.Context, _ *v2.Resource, _ *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

// Delete removes a user from App Store Connect.
func (u *userBuilder) Delete(ctx context.Context, resourceId *v2.ResourceId) (annotations.Annotations, error) {
	err := u.client.DeleteUser(ctx, resourceId.Resource)
	if err != nil {
		return nil, fmt.Errorf("baton-appstoreconnect: failed to delete user %s: %w", resourceId.Resource, err)
	}

	return nil, nil
}

func newUserBuilder(client *client.Client) *userBuilder {
	return &userBuilder{
		client: client,
	}
}
