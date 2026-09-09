package connector

import (
	"context"
	"fmt"
	"strings"

	"github.com/conductorone/baton-appstoreconnect/pkg/client"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	resourceSdk "github.com/conductorone/baton-sdk/pkg/types/resource"
)

type betaTesterBuilder struct {
	client *client.Client
}

func (b *betaTesterBuilder) ResourceType(_ context.Context) *v2.ResourceType {
	return betaTesterResourceType
}

func newBetaTesterResource(tester client.BetaTester) (*v2.Resource, error) {
	displayName := strings.TrimSpace(fmt.Sprintf("%s %s", tester.Attributes.FirstName, tester.Attributes.LastName))
	if displayName == "" {
		displayName = tester.Attributes.Email
	}

	profile := map[string]interface{}{
		"first_name":  tester.Attributes.FirstName,
		"last_name":   tester.Attributes.LastName,
		"invite_type": tester.Attributes.InviteType,
		"state":       tester.Attributes.State,
	}

	opts := []resourceSdk.UserTraitOption{
		resourceSdk.WithUserProfile(profile),
	}
	if tester.Attributes.Email != "" {
		opts = append(opts, resourceSdk.WithEmail(tester.Attributes.Email, true))
	}

	return resourceSdk.NewUserResource(
		displayName,
		betaTesterResourceType,
		tester.ID,
		opts,
	)
}

// List returns all TestFlight beta testers in the account.
func (b *betaTesterBuilder) List(ctx context.Context, _ *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	var cursor string
	if pToken != nil {
		cursor = pToken.Token
	}

	testers, links, err := b.client.ListBetaTesters(ctx, cursor)
	if err != nil {
		return nil, "", nil, fmt.Errorf("baton-appstoreconnect: failed to list beta testers: %w", err)
	}

	ret := make([]*v2.Resource, 0, len(testers))
	for _, tester := range testers {
		resource, err := newBetaTesterResource(tester)
		if err != nil {
			return nil, "", nil, fmt.Errorf("baton-appstoreconnect: failed to create beta tester resource for %s: %w", tester.ID, err)
		}
		ret = append(ret, resource)
	}

	nextCursor := ""
	if links != nil && links.Next != "" {
		nextCursor = links.Next
	}

	return ret, nextCursor, nil, nil
}

// Entitlements returns an empty slice for beta testers (they don't own entitlements).
func (b *betaTesterBuilder) Entitlements(_ context.Context, _ *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

// Grants returns an empty slice for beta testers. Beta group membership is
// synced as grants on the beta group resource.
func (b *betaTesterBuilder) Grants(_ context.Context, _ *v2.Resource, _ *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

// Delete removes a beta tester from TestFlight entirely, across every group and app.
func (b *betaTesterBuilder) Delete(ctx context.Context, resourceId *v2.ResourceId) (annotations.Annotations, error) {
	err := b.client.DeleteBetaTester(ctx, resourceId.Resource)
	if err != nil {
		return nil, fmt.Errorf("baton-appstoreconnect: failed to delete beta tester %s: %w", resourceId.Resource, err)
	}

	return nil, nil
}

func newBetaTesterBuilder(client *client.Client) *betaTesterBuilder {
	return &betaTesterBuilder{
		client: client,
	}
}
