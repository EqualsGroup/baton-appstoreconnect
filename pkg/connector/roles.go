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

const (
	// roleResourceID is the static ID for the App Store Connect role resource.
	roleResourceID = "appstoreconnect"
)

// roleDescription maps role slugs to human-readable descriptions.
var roleDescriptions = map[string]string{
	"ADMIN":                          "Full access to manage all aspects of App Store Connect",
	"FINANCE":                        "Manage financial information including reports and tax forms",
	"ACCOUNT_HOLDER":                 "The account holder has full access and is responsible for the account",
	"SALES":                          "View sales and download reports",
	"DEVELOPER":                      "Upload builds, manage app metadata, and submit apps for review",
	"APP_MANAGER":                    "Manage app information, pricing, and availability",
	"CUSTOMER_SUPPORT":               "Respond to customer reviews and manage customer issues",
	"MARKETING":                      "Manage marketing assets and promotional content",
	"CREATE_APPS":                    "Create new apps in App Store Connect",
	"CLOUD_MANAGED_DEVELOPER_ID":     "Manage cloud-managed Developer ID signing",
	"CLOUD_MANAGED_APP_DISTRIBUTION": "Manage cloud-managed app distribution",
	"ACCESS_TO_REPORTS":              "Access to financial and sales reports",
	"GENERATE_INDIVIDUAL_KEYS":       "Generate individual API keys",
	"IMAGE_MANAGER":                  "Manage app images and screenshots",
}

type roleBuilder struct {
	client *client.Client
}

func (r *roleBuilder) ResourceType(_ context.Context) *v2.ResourceType {
	return roleResourceType
}

// List returns a single static resource representing the App Store Connect account.
// Roles are modeled as entitlements on this resource.
func (r *roleBuilder) List(_ context.Context, _ *v2.ResourceId, _ *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	profile := map[string]interface{}{
		"type": "App Store Connect Account",
	}

	resource, err := resourceSdk.NewGroupResource(
		"App Store Connect",
		roleResourceType,
		roleResourceID,
		[]resourceSdk.GroupTraitOption{
			resourceSdk.WithGroupProfile(profile),
		},
	)
	if err != nil {
		return nil, "", nil, fmt.Errorf("baton-appstoreconnect: failed to create role resource: %w", err)
	}

	return []*v2.Resource{resource}, "", nil, nil
}

// Entitlements returns one permission entitlement per App Store Connect role.
func (r *roleBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	entitlements := make([]*v2.Entitlement, 0, len(client.AppStoreConnectRoles))

	for _, role := range client.AppStoreConnectRoles {
		desc, ok := roleDescriptions[role]
		if !ok {
			desc = fmt.Sprintf("%s role in App Store Connect", role)
		}

		entitlements = append(entitlements, entitlement.NewPermissionEntitlement(
			resource,
			role,
			entitlement.WithDescription(desc),
			entitlement.WithDisplayName(fmt.Sprintf("App Store Connect %s", role)),
			entitlement.WithGrantableTo(userResourceType),
		))
	}

	return entitlements, "", nil, nil
}

// Grants returns a grant for each user's assigned roles.
func (r *roleBuilder) Grants(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	var cursor string
	if pToken != nil {
		cursor = pToken.Token
	}

	users, links, err := r.client.ListUsers(ctx, cursor)
	if err != nil {
		return nil, "", nil, fmt.Errorf("baton-appstoreconnect: failed to list users for role grants: %w", err)
	}

	var ret []*v2.Grant
	for _, user := range users {
		for _, role := range user.Attributes.Roles {
			resourceId, err := resourceSdk.NewResourceID(userResourceType, user.ID)
			if err != nil {
				return nil, "", nil, fmt.Errorf("baton-appstoreconnect: failed to create resource ID for user %s: %w", user.ID, err)
			}
			ret = append(ret, grant.NewGrant(resource, role, resourceId))
		}
	}

	nextCursor := ""
	if links != nil && links.Next != "" {
		nextCursor = links.Next
	}

	return ret, nextCursor, nil, nil
}

func newRoleBuilder(client *client.Client) *roleBuilder {
	return &roleBuilder{
		client: client,
	}
}
