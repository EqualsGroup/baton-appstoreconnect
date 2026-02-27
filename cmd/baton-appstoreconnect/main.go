package main

import (
	"context"
	"fmt"
	"os"

	"github.com/conductorone/baton-sdk/pkg/config"
	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"
	"github.com/conductorone/baton-sdk/pkg/field"
	"github.com/conductorone/baton-sdk/pkg/types"
	cfg "github.com/conductorone/baton-appstoreconnect/pkg/config"
	"github.com/conductorone/baton-appstoreconnect/pkg/connector"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
)

var version = "dev"

func main() {
	ctx := context.Background()

	_, cmd, err := config.DefineConfiguration(
		ctx,
		"baton-appstoreconnect",
		getConnector[*cfg.AppStoreConnect],
		cfg.Config,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	cmd.Version = version

	err = cmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func getConnector[T field.Configurable](ctx context.Context, c T) (types.ConnectorServer, error) {
	l := ctxzap.Extract(ctx)
	if err := field.Validate(cfg.Config, c); err != nil {
		return nil, err
	}

	issuerID := c.GetString(cfg.IssuerID.FieldName)
	keyID := c.GetString(cfg.KeyID.FieldName)
	privateKeyPath := c.GetString(cfg.PrivateKeyPath.FieldName)

	cb, err := connector.New(ctx, issuerID, keyID, privateKeyPath)
	if err != nil {
		l.Error("error creating connector", zap.Error(err))
		return nil, err
	}

	srv, err := connectorbuilder.NewConnector(ctx, cb)
	if err != nil {
		l.Error("error creating connector server", zap.Error(err))
		return nil, err
	}

	return srv, nil
}
