package connector

import (
	"context"
	"fmt"

	"github.com/conductorone/baton-appstoreconnect/pkg/client"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	"github.com/conductorone/baton-sdk/pkg/types/entitlement"
	"github.com/conductorone/baton-sdk/pkg/types/grant"
	resourceSdk "github.com/conductorone/baton-sdk/pkg/types/resource"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
)

// betaGroupMember is the entitlement granted to a beta tester who belongs to a
// TestFlight beta group.
const betaGroupMember = "member"

type betaGroupBuilder struct {
	client *client.Client
}

func (b *betaGroupBuilder) ResourceType(_ context.Context) *v2.ResourceType {
	return betaGroupResourceType
}

func newBetaGroupResource(group client.BetaGroup) (*v2.Resource, error) {
	profile := map[string]interface{}{
		"app_id":                   group.AppID(),
		"is_internal_group":        group.Attributes.IsInternalGroup,
		"has_access_to_all_builds": group.Attributes.HasAccessToAllBuilds,
		"public_link_enabled":      group.Attributes.PublicLinkEnabled,
		"public_link":              group.Attributes.PublicLink,
		"feedback_enabled":         group.Attributes.FeedbackEnabled,
		"created_date":             group.Attributes.CreatedDate,
	}

	return resourceSdk.NewGroupResource(
		group.Attributes.Name,
		betaGroupResourceType,
		group.ID,
		[]resourceSdk.GroupTraitOption{
			resourceSdk.WithGroupProfile(profile),
		},
	)
}

// List returns all TestFlight beta groups across every app in the account.
func (b *betaGroupBuilder) List(ctx context.Context, _ *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	var cursor string
	if pToken != nil {
		cursor = pToken.Token
	}

	groups, links, err := b.client.ListBetaGroups(ctx, cursor)
	if err != nil {
		return nil, "", nil, fmt.Errorf("baton-appstoreconnect: failed to list beta groups: %w", err)
	}

	ret := make([]*v2.Resource, 0, len(groups))
	for _, group := range groups {
		resource, err := newBetaGroupResource(group)
		if err != nil {
			return nil, "", nil, fmt.Errorf("baton-appstoreconnect: failed to create beta group resource for %s: %w", group.ID, err)
		}
		ret = append(ret, resource)
	}

	nextCursor := ""
	if links != nil && links.Next != "" {
		nextCursor = links.Next
	}

	return ret, nextCursor, nil, nil
}

// Entitlements returns a "member" entitlement for each beta group, grantable to
// beta testers.
func (b *betaGroupBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	return []*v2.Entitlement{
		entitlement.NewAssignmentEntitlement(
			resource,
			betaGroupMember,
			entitlement.WithDescription(fmt.Sprintf("Member of the %s TestFlight beta group", resource.DisplayName)),
			entitlement.WithDisplayName(fmt.Sprintf("%s Beta Group Member", resource.DisplayName)),
			entitlement.WithGrantableTo(betaTesterResourceType),
		),
	}, "", nil, nil
}

// Grants returns a grant for each beta tester in this beta group.
func (b *betaGroupBuilder) Grants(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	var cursor string
	if pToken != nil {
		cursor = pToken.Token
	}

	testers, links, err := b.client.ListBetaGroupTesters(ctx, resource.Id.Resource, cursor)
	if err != nil {
		return nil, "", nil, fmt.Errorf("baton-appstoreconnect: failed to list testers for beta group %s: %w", resource.Id.Resource, err)
	}

	ret := make([]*v2.Grant, 0, len(testers))
	for _, tester := range testers {
		resourceId, err := resourceSdk.NewResourceID(betaTesterResourceType, tester.ID)
		if err != nil {
			return nil, "", nil, fmt.Errorf("baton-appstoreconnect: failed to create resource ID for beta tester %s: %w", tester.ID, err)
		}
		ret = append(ret, grant.NewGrant(resource, betaGroupMember, resourceId))
	}

	nextCursor := ""
	if links != nil && links.Next != "" {
		nextCursor = links.Next
	}

	return ret, nextCursor, nil, nil
}

// Grant adds a beta tester to a TestFlight beta group. Adding a tester to a
// group is what sends them a TestFlight invitation.
func (b *betaGroupBuilder) Grant(ctx context.Context, principal *v2.Resource, entitlementToGrant *v2.Entitlement) (annotations.Annotations, error) {
	l := ctxzap.Extract(ctx)

	if principal.Id.ResourceType != betaTesterResourceType.Id {
		return nil, fmt.Errorf("baton-appstoreconnect: only beta testers can be added to beta groups, got %s", principal.Id.ResourceType)
	}

	groupID := entitlementToGrant.Resource.Id.Resource

	err := b.client.AddBetaTesterToGroup(ctx, groupID, principal.Id.Resource)
	if err != nil {
		return nil, fmt.Errorf("baton-appstoreconnect: failed to add beta tester %s to beta group %s: %w", principal.Id.Resource, groupID, err)
	}

	l.Info("baton-appstoreconnect: added beta tester to beta group",
		zap.String("beta_tester_id", principal.Id.Resource),
		zap.String("beta_group_id", groupID),
	)

	return nil, nil
}

// Revoke removes a beta tester from a TestFlight beta group. The tester keeps
// their other group memberships.
func (b *betaGroupBuilder) Revoke(ctx context.Context, grantToRevoke *v2.Grant) (annotations.Annotations, error) {
	l := ctxzap.Extract(ctx)

	principal := grantToRevoke.Principal
	if principal.Id.ResourceType != betaTesterResourceType.Id {
		return nil, fmt.Errorf("baton-appstoreconnect: only beta testers can be removed from beta groups, got %s", principal.Id.ResourceType)
	}

	groupID := grantToRevoke.Entitlement.Resource.Id.Resource

	err := b.client.RemoveBetaTesterFromGroup(ctx, groupID, principal.Id.Resource)
	if err != nil {
		return nil, fmt.Errorf("baton-appstoreconnect: failed to remove beta tester %s from beta group %s: %w", principal.Id.Resource, groupID, err)
	}

	l.Info("baton-appstoreconnect: removed beta tester from beta group",
		zap.String("beta_tester_id", principal.Id.Resource),
		zap.String("beta_group_id", groupID),
	)

	return nil, nil
}

func newBetaGroupBuilder(client *client.Client) *betaGroupBuilder {
	return &betaGroupBuilder{
		client: client,
	}
}
