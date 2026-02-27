package connector

import (
	"context"
	"fmt"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	"github.com/conductorone/baton-sdk/pkg/types/entitlement"
	"github.com/conductorone/baton-sdk/pkg/types/grant"
	resourceSdk "github.com/conductorone/baton-sdk/pkg/types/resource"
	"github.com/conductorone/baton-appstoreconnect/pkg/client"
)

const appAccess = "access"

type appBuilder struct {
	client *client.Client
}

func (a *appBuilder) ResourceType(_ context.Context) *v2.ResourceType {
	return appResourceType
}

func newAppResource(app client.App) (*v2.Resource, error) {
	profile := map[string]interface{}{
		"bundle_id": app.Attributes.BundleID,
		"sku":       app.Attributes.SKU,
	}

	return resourceSdk.NewAppResource(
		app.Attributes.Name,
		appResourceType,
		app.ID,
		[]resourceSdk.AppTraitOption{
			resourceSdk.WithAppProfile(profile),
		},
	)
}

// List returns all apps from App Store Connect.
func (a *appBuilder) List(ctx context.Context, _ *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	var cursor string
	if pToken != nil {
		cursor = pToken.Token
	}

	apps, links, err := a.client.ListApps(ctx, cursor)
	if err != nil {
		return nil, "", nil, fmt.Errorf("baton-appstoreconnect: failed to list apps: %w", err)
	}

	ret := make([]*v2.Resource, 0, len(apps))
	for _, app := range apps {
		resource, err := newAppResource(app)
		if err != nil {
			return nil, "", nil, fmt.Errorf("baton-appstoreconnect: failed to create app resource for %s: %w", app.ID, err)
		}
		ret = append(ret, resource)
	}

	nextCursor := ""
	if links != nil && links.Next != "" {
		nextCursor = links.Next
	}

	return ret, nextCursor, nil, nil
}

// Entitlements returns an "access" entitlement for each app, grantable to users.
// This represents a user having access to a specific app.
func (a *appBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	return []*v2.Entitlement{
		entitlement.NewAssignmentEntitlement(
			resource,
			appAccess,
			entitlement.WithDescription(fmt.Sprintf("Access to %s app", resource.DisplayName)),
			entitlement.WithDisplayName(fmt.Sprintf("%s Access", resource.DisplayName)),
			entitlement.WithGrantableTo(userResourceType),
		),
	}, "", nil, nil
}

// Grants returns grants for users who have access to this specific app.
// We iterate all users and check which ones have this app in their visible apps list.
func (a *appBuilder) Grants(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	// We use a compound cursor: first paginate through users, then for each user
	// check if they have access to this app. Since the API doesn't support
	// listing users per app directly, we need to iterate users.
	var cursor string
	if pToken != nil {
		cursor = pToken.Token
	}

	users, links, err := a.client.ListUsers(ctx, cursor)
	if err != nil {
		return nil, "", nil, fmt.Errorf("baton-appstoreconnect: failed to list users for app grants: %w", err)
	}

	appID := resource.Id.Resource
	var ret []*v2.Grant

	for _, user := range users {
		// Users with allAppsVisible have access to all apps.
		if user.Attributes.AllAppsVisible {
			resourceId, err := resourceSdk.NewResourceID(userResourceType, user.ID)
			if err != nil {
				return nil, "", nil, fmt.Errorf("baton-appstoreconnect: failed to create resource ID for user %s: %w", user.ID, err)
			}
			ret = append(ret, grant.NewGrant(resource, appAccess, resourceId))
			continue
		}

		// For users without allAppsVisible, check their visible apps.
		hasAccess, err := a.userHasAppAccess(ctx, user.ID, appID)
		if err != nil {
			return nil, "", nil, fmt.Errorf("baton-appstoreconnect: failed to check app access for user %s: %w", user.ID, err)
		}
		if hasAccess {
			resourceId, err := resourceSdk.NewResourceID(userResourceType, user.ID)
			if err != nil {
				return nil, "", nil, fmt.Errorf("baton-appstoreconnect: failed to create resource ID for user %s: %w", user.ID, err)
			}
			ret = append(ret, grant.NewGrant(resource, appAccess, resourceId))
		}
	}

	nextCursor := ""
	if links != nil && links.Next != "" {
		nextCursor = links.Next
	}

	return ret, nextCursor, nil, nil
}

// userHasAppAccess checks if a user has access to a specific app by paginating
// through their visible apps.
func (a *appBuilder) userHasAppAccess(ctx context.Context, userID, appID string) (bool, error) {
	cursor := ""
	for {
		apps, links, err := a.client.ListUserVisibleApps(ctx, userID, cursor)
		if err != nil {
			return false, err
		}

		for _, app := range apps {
			if app.ID == appID {
				return true, nil
			}
		}

		if links == nil || links.Next == "" {
			break
		}
		cursor = links.Next
	}

	return false, nil
}

func newAppBuilder(client *client.Client) *appBuilder {
	return &appBuilder{
		client: client,
	}
}
