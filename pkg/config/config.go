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
	)

	PrivateKey = field.StringField(
		"private-key",
		field.WithDisplayName("Private Key"),
		field.WithDescription("PEM-encoded .p8 private key contents for App Store Connect API"),
		field.WithIsSecret(true),
	)

	ConfigurationFields = []field.SchemaField{IssuerID, KeyID, PrivateKeyPath, PrivateKey}

	FieldRelationships = []field.SchemaFieldRelationship{
		field.FieldsAtLeastOneUsed(PrivateKeyPath, PrivateKey),
		field.FieldsMutuallyExclusive(PrivateKeyPath, PrivateKey),
	}
)

var Config = field.NewConfiguration(
	ConfigurationFields,
	field.WithConstraints(FieldRelationships...),
	field.WithConnectorDisplayName("App Store Connect"),
)
