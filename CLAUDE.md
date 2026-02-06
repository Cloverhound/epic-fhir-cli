# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

fhir-cli is a Go CLI application for querying Epic's FHIR R4 APIs. It uses OAuth2 backend authentication with JWT assertions.

## Build Commands

```bash
make build         # Build the binary
make install       # Install to GOPATH/bin
make test          # Run tests
make lint          # Run golangci-lint
make fmt           # Format code
make tidy          # go mod tidy
make build-all     # Cross-compile for all platforms
```

## Architecture

```
fhir-cli/
├── main.go                 # Entry point
├── cmd/                    # Cobra CLI commands
│   ├── root.go            # Root command, global flags
│   ├── config.go          # Configuration management
│   ├── auth.go            # OAuth2 token commands
│   ├── resource.go        # Generic resource operations
│   ├── patient.go         # Patient resource
│   ├── observation.go     # Observation resource
│   ├── condition.go       # Condition resource
│   ├── medication.go      # MedicationRequest resource
│   └── ...                # Other resource commands
├── internal/
│   ├── config/            # Configuration and profile management
│   │   └── config.go      # Config file handling, token cache
│   ├── auth/              # OAuth2 authentication
│   │   ├── jwt.go         # JWT generation with RS384 signing
│   │   └── oauth.go       # Token exchange, caching
│   ├── fhir/              # FHIR client
│   │   ├── client.go      # HTTP client, API operations
│   │   └── resources.go   # Resource type definitions
│   └── output/            # Output formatting
│       └── formatter.go   # JSON, YAML, table formatters
└── Makefile
```

## Key Patterns

### Adding a New Resource Command

1. Create `cmd/<resource>.go`
2. Define cobra commands for get, search, list operations
3. Use `fhir.NewClient(config.CurrentProfile)` for API calls
4. Use `output.Print(result, output.ParseFormat(GetOutput()))` for output

Example pattern:
```go
var resourceCmd = &cobra.Command{
    Use:   "resource",
    Short: "Resource operations",
}

var resourceGetCmd = &cobra.Command{
    Use:   "get <id>",
    Args:  cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        client := fhir.NewClient(config.CurrentProfile)
        result, err := client.Get("ResourceType", args[0])
        if err != nil {
            return err
        }
        return output.Print(result, output.ParseFormat(GetOutput()))
    },
}
```

### Authentication Flow

1. Load profile from `~/.fhir-cli/config.yaml`
2. Generate JWT with RS384 using private key
3. Exchange JWT for access token at token endpoint
4. Cache token in `~/.fhir-cli/token_<profile>.yaml`
5. Reuse cached token until 5 min before expiry

### FHIR Client

The `fhir.Client` handles all HTTP operations:
- `Get(resourceType, id)` - Fetch single resource
- `Search(resourceType, params)` - Search with query params
- `Create(resourceType, resource)` - Create resource
- `Update(resourceType, id, resource)` - Update resource
- `Operation(resourceType, id, operation, params)` - FHIR operations

### Output Formatting

Three formats supported via `-o` flag:
- `json` (default): Pretty-printed JSON
- `yaml`: YAML format
- `table`: Formatted ASCII tables with resource-specific layouts

## Configuration

Config stored in `~/.fhir-cli/config.yaml`:
- Multiple profiles supported
- Each profile has: client_id, private_key path, token_url, fhir_base_url, scopes
- Token cache per profile in `~/.fhir-cli/token_<profile>.yaml`

## Testing with Epic Sandbox

Use Epic's public sandbox:
- Base URL: `https://fhir.epic.com/interconnect-fhir-oauth/api/FHIR/R4`
- Token URL: `https://fhir.epic.com/interconnect-fhir-oauth/oauth2/token`

Register app at https://fhir.epic.com to get client credentials.
