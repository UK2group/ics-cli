# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

ICS CLI is a command-line interface for managing Ingenuity Cloud Services (ICS) resources. It's built in Go using the Cobra framework for CLI commands and Viper for configuration management. The tool provides commands for managing bare metal servers, SSH keys, and authentication.

## Development Commands

### Building
```bash
# Build for current platform
go build -v ./...

# Build with version information (example)
VERSION="1.0.0"
BUILDTIME=$(date -u '+%Y-%m-%d %H:%M:%S')
COMMIT=$(git rev-parse HEAD)
go build -ldflags="-X 'main.Version=$VERSION' -X 'main.BuildTime=$BUILDTIME' -X 'main.Commit=$COMMIT'"

# Build for specific platforms (used by release workflow)
GOOS=darwin GOARCH=amd64 go build -o builds/ics-cli-macos-amd64
GOOS=darwin GOARCH=arm64 go build -o builds/ics-cli-macos-arm64
GOOS=linux GOARCH=amd64 go build -o builds/ics-cli-linux-amd64
GOOS=windows GOARCH=amd64 go build -o builds/ics-cli-windows-amd64.exe
```

### Testing and Linting
```bash
# Run tests
go test -v ./...

# Run go vet
go vet ./...

# Download and verify dependencies
go mod download
go mod verify
```

### Running Locally
```bash
# Run without building
go run main.go [command]

# Run built binary
./ics-cli [command]

# Common commands for testing
./ics-cli auth login
./ics-cli baremetal list
./ics-cli sshkeys list
```

## Architecture

### Project Structure

```
ics-cli/
├── main.go                    # Entry point - calls cmd.Execute()
├── cmd/                       # All CLI commands live here
│   ├── root.go               # Root command and config initialization
│   ├── helpers.go            # Shared helpers (API requests, colors, browser opening)
│   ├── auth.go               # Auth command group
│   ├── auth_*.go            # Auth subcommands (login, logout, check)
│   ├── baremetal.go         # Baremetal command group
│   ├── baremetal_*.go       # Baremetal subcommands
│   ├── baremetal_api_response.go  # API response structs
│   ├── baremetal_helpers.go      # Baremetal-specific helpers
│   ├── baremetal_deploy.go       # Deploy command group
│   ├── baremetal_deploy_*.go     # Deploy subcommands
│   ├── sshkeys.go           # SSH keys command group
│   ├── sshkeys_*.go         # SSH key subcommands
│   ├── sshkeys_api_response.go   # SSH keys API response structs
│   └── sshkeys_helpers.go        # SSH keys helpers
```

### Command Pattern

All commands follow Cobra's structure:
1. Each command group (auth, baremetal, sshkeys) has a main file defining the parent command
2. Subcommands are in separate files following the naming pattern `{group}_{subcommand}.go`
3. Commands are registered with their parent in `init()` functions
4. The `baremetal` command has aliases: `bm`, `metal`

### API Integration

**Base API URL**: Controlled by `getAPIBaseURL()` in cmd/helpers.go:
- Default: `https://api.ingenuitycloudservices.com`
- Override: Set `ICS_API_URL` environment variable (used for staging testing)

**Authentication**:
- API key stored in `~/.ics-cli.yaml` config file
- All API requests use `X-Api-Token` header
- Managed by Viper configuration library

**API Request Pattern**: The `makeAPIRequest()` function in cmd/helpers.go handles all API calls:
- Takes: HTTP method, timeout, URL, request body, result interface
- Handles: Authentication, error responses, JSON parsing
- Returns: Unmarshaled response or error

### API Response Types

All API response structures are defined in:
- `cmd/baremetal_api_response.go` - Server, inventory, deployment responses
- `cmd/sshkeys_api_response.go` - SSH key responses

**FlexibleID Type**: Custom JSON unmarshaler that handles API fields that may be either string or numeric (see baremetal_api_response.go:8-41)

### Helper Functions

**cmd/helpers.go** provides shared utilities:
- `makeAPIRequest()` - Generic API request handler
- `openBrowser()` - Cross-platform browser opening
- Color formatters: `BlueHeading()`, `WhiteText()`, `GreenText()`, `RedText()`, `YellowText()`

**cmd/baremetal_helpers.go** provides baremetal-specific functions:
- Server ID lookups and conversions
- Power management (on/off/reboot/recovery)
- Remote access (IKVM, SOL)
- Inventory filtering and grouping
- Order placement and confirmation
- Server cancellation for hourly-billed servers

**cmd/sshkeys_helpers.go** provides SSH key operations

### Configuration

Viper manages configuration with these priorities (highest to lowest):
1. Command-line flag: `--config`
2. Environment variable: `ICS_CONFIG_FILE`
3. Default: `~/.ics-cli.yaml`

Auto-loaded environment variables with `viper.AutomaticEnv()` (e.g., `ICS_API_KEY`)

## Important Development Notes

### API Service ID vs Server ID
- The API returns both `service_id` (integer) and `id` (can be string or numeric)
- `service_id` is the user-facing identifier shown in commands
- Internal `id` is used for API calls
- Use `getServerIDFromServiceID()` to convert between them

### Hourly Billing
- Servers can be billed monthly or hourly
- Only hourly-billed servers can be cancelled via CLI (`baremetal cancel` command)
- The `HourlyBilling` and `HourlyPrice` fields are part of server details
- Inventory API includes `hourly_enabled` and `price_hourly` fields

### Color Output
All user-facing output uses the fatih/color library through helper functions. This ensures consistent formatting across commands and proper handling of terminal color support.

### Browser Integration
Commands that open URLs (like `baremetal ikvm` and `baremetal sol`) use `openBrowser()` which handles platform differences (macOS, Linux, Windows).

## Release Process

Releases are automated via GitHub Actions:
- Triggered by pushing a tag matching `v*` (e.g., `v1.2.3`)
- Builds binaries for macOS (amd64, arm64), Linux (amd64), Windows (amd64)
- Generates SHA256 checksums
- Creates GitHub release with auto-generated changelog
- See `.github/workflows/release.yaml` for details

## Environment Variables

- `ICS_API_KEY` - API authentication key (alternative to config file)
- `ICS_API_URL` - Override API base URL (for staging/testing)
- `ICS_CONFIG_FILE` - Custom config file path

## Adding New Commands

1. Create file: `cmd/{group}_{command}.go`
2. Define cobra.Command with Use, Short, Run fields
3. Add command to parent in `init()` function
4. Add flags using `cmd.Flags()` or `cmd.PersistentFlags()`
5. Use `makeAPIRequest()` for API calls
6. Define response structs in appropriate `*_api_response.go` file
7. Use color helpers for output formatting
