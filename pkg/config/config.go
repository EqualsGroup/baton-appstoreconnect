package config

import (
	"github.com/conductorone/baton-sdk/pkg/field"
)

var (
	IssuerID = field.StringField(
		"issuer-id",
		field.WithDisplayName("Issuer ID"),
		field.WithDescription("App Store Connect API Issuer ID"),
		field.WithRequired(true),
	)

	KeyID = field.StringField(
		"key-id",
		field.WithDisplayName("Key ID"),
		field.WithDescription("App Store Connect API Key ID"),
		field.WithRequired(true),
	)

	PrivateKeyPath = field.StringField(
		"private-key-path",
		field.WithDisplayName("Private Key Path"),
		field.WithDescription("Path to the .p8 private key file for App Store Connect API"),
		field.WithRequired(true),
	)

	ConfigurationFields = []field.SchemaField{IssuerID, KeyID, PrivateKeyPath}

	FieldRelationships = []field.SchemaFieldRelationship{}
)

var Config = field.NewConfiguration(
	ConfigurationFields,
	field.WithConstraints(FieldRelationships...),
	field.WithConnectorDisplayName("App Store Connect"),
)
