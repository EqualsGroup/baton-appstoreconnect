# `baton-appstoreconnect`

`baton-appstoreconnect` is a connector for Apple App Store Connect built using the [Baton SDK](https://github.com/conductorone/baton-sdk). It communicates with the [App Store Connect API](https://developer.apple.com/documentation/appstoreconnectapi) to sync data about users, apps, and role assignments.

Check out [Baton](https://github.com/conductorone/baton) to learn more about the project in general.

## Prerequisites

To use this connector you need an **App Store Connect API key**. To create one:

1. Sign in to [App Store Connect](https://appstoreconnect.apple.com/)
2. Go to **Users and Access** → **Integrations** → **App Store Connect API**
3. Click **Generate API Key** (requires Admin or Account Holder role)
4. Note your **Issuer ID** (shown at the top of the page)
5. Note the **Key ID** of the generated key
6. Download the `.p8` private key file (you can only download it once)

## Getting Started

### source

```bash
go install github.com/conductorone/baton/cmd/baton@main
go install github.com/EqualsGroup/baton-appstoreconnect/cmd/baton-appstoreconnect@main

baton-appstoreconnect \
  --issuer-id "your-issuer-id" \
  --key-id "your-key-id" \
  --private-key-path "/path/to/AuthKey_XXXXXXXX.p8"

baton resources
```

### Environment variables

All flags can be set via environment variables:

```bash
export BATON_ISSUER_ID="your-issuer-id"
export BATON_KEY_ID="your-key-id"
export BATON_PRIVATE_KEY_PATH="/path/to/AuthKey_XXXXXXXX.p8"

baton-appstoreconnect
baton resources
```

## Data Model

`baton-appstoreconnect` syncs the following resources:

| Resource Type | Description |
|--------------|-------------|
| User | App Store Connect users with their email, name, and assigned roles. |
| App | Apps in the account, with name, bundle ID, and SKU. Each app exposes an `access` entitlement. |
| Role | A single resource representing the App Store Connect account. Exposes one entitlement per role (Admin, Developer, App Manager, etc.). |

## Provisioning

| Action | Resource | Description |
|--------|----------|-------------|
| Delete | User | Remove a user from App Store Connect. |
| Grant | App Access | Users with `allAppsVisible` or explicit app visibility get app access grants. |
| Grant | Role | Each user's role assignments are synced as grants on the Role resource. |

## Contributing, Support and Issues

We welcome contributions and ideas. If you have questions, problems, or ideas: please open a GitHub Issue!

See [CONTRIBUTING.md](https://github.com/ConductorOne/baton/blob/main/CONTRIBUTING.md) for more details.

## Command Line Usage

```
baton-appstoreconnect

Usage:
  baton-appstoreconnect [flags]
  baton-appstoreconnect [command]

Available Commands:
  capabilities       Get connector capabilities
  completion         Generate the autocompletion script for the specified shell
  help               Help about any command

Flags:
      --issuer-id string           required: App Store Connect API Issuer ID ($BATON_ISSUER_ID)
      --key-id string              required: App Store Connect API Key ID ($BATON_KEY_ID)
      --private-key-path string    required: Path to the .p8 private key file ($BATON_PRIVATE_KEY_PATH)
      --client-id string           The client ID used to authenticate with ConductorOne ($BATON_CLIENT_ID)
      --client-secret string       The client secret used to authenticate with ConductorOne ($BATON_CLIENT_SECRET)
  -f, --file string                The path to the c1z file to sync with ($BATON_FILE) (default "sync.c1z")
  -h, --help                       help for baton-appstoreconnect
      --log-format string          The output format for logs: json, console ($BATON_LOG_FORMAT) (default "json")
      --log-level string           The log level: debug, info, warn, error ($BATON_LOG_LEVEL) (default "info")
  -p, --provisioning               This must be set in order for provisioning actions to be enabled ($BATON_PROVISIONING)
  -v, --version                    version for baton-appstoreconnect
```
