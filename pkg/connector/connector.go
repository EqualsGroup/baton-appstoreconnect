package connector

import (
	"context"
	"io"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"
	"github.com/conductorone/baton-appstoreconnect/pkg/client"
)

// Connector implements the baton connector interface for App Store Connect.
type Connector struct {
	client *client.Client
}

// ResourceSyncers returns a ResourceSyncer for each resource type that should
// be synced from App Store Connect.
func (c *Connector) ResourceSyncers(_ context.Context) []connectorbuilder.ResourceSyncer {
	return []connectorbuilder.ResourceSyncer{
		newUserBuilder(c.client),
		newAppBuilder(c.client),
		newRoleBuilder(c.client),
	}
}

// Asset takes an input AssetRef and attempts to fetch it using the connector's
// authenticated HTTP client. Not implemented for App Store Connect.
func (c *Connector) Asset(_ context.Context, _ *v2.AssetRef) (string, io.ReadCloser, error) {
	return "", nil, nil
}

// Metadata returns metadata about the connector.
func (c *Connector) Metadata(_ context.Context) (*v2.ConnectorMetadata, error) {
	return &v2.ConnectorMetadata{
		DisplayName: "App Store Connect",
		Description: "Connector for Apple App Store Connect. Syncs users, apps, roles, and app-specific access.",
	}, nil
}

// Validate is called to ensure that the connector is properly configured.
// It exercises the API credentials by listing users with a limit of 1.
func (c *Connector) Validate(ctx context.Context) (annotations.Annotations, error) {
	_, _, err := c.client.ListUsers(ctx, "")
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// New returns a new instance of the App Store Connect connector.
func New(ctx context.Context, issuerID, keyID, privateKeyPath string) (*Connector, error) {
	cl, err := client.New(ctx, issuerID, keyID, privateKeyPath)
	if err != nil {
		return nil, err
	}

	return &Connector{
		client: cl,
	}, nil
}
