# ICS CLI

A powerful command-line interface for managing Ingenuity Cloud Services resources.

[![Release](https://img.shields.io/github/v/release/UK2Group/ics-cli)](https://github.com/UK2Group/ics-cli/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/UK2Group/ics-cli)](https://goreportcard.com/report/github.com/UK2Group/ics-cli)
[![License](https://img.shields.io/github/license/UK2Group/ics-cli)](https://github.com/UK2Group/ics-cli/blob/main/LICENSE)

## Overview

ICS CLI provides a convenient way to manage your Ingenuity Cloud Services resources through a simple command-line interface. It allows you to provision and manage baremetal servers, VPS instances, and more without needing to use the web interface.

## Features

- **Baremetal Server Management**: List, deploy, power control, and manage your dedicated servers
- **VPS Management**: Deploy, manage, and control virtual private servers across multiple datacenters
- **Hourly Billing Support**: Deploy servers with hourly billing and cancel hourly-billed servers
- **SSH Key Management**: Create, list, and manage SSH keys for server access
- **Inventory Browsing**: View available hardware in different datacenters with hourly pricing options
- **Add-on Management**: Browse and select operating systems, software licenses, and support levels

## Installation

### macOS

#### Homebrew (Recommended)

```bash
brew tap UK2Group/ics-cli
brew install ics-cli
```

#### Manual Installation

```bash
# For Intel Macs
curl -L https://github.com/UK2Group/ics-cli/releases/latest/download/ics-cli-macos-amd64 -o /usr/local/bin/ics-cli
chmod +x /usr/local/bin/ics-cli

# For Apple Silicon (M1/M2) Macs
curl -L https://github.com/UK2Group/ics-cli/releases/latest/download/ics-cli-macos-arm64 -o /usr/local/bin/ics-cli
chmod +x /usr/local/bin/ics-cli
```

### Linux

```bash
curl -L https://github.com/UK2Group/ics-cli/releases/latest/download/ics-cli-linux-amd64 -o /usr/local/bin/ics-cli
chmod +x /usr/local/bin/ics-cli
```

### Windows

1. Download the [latest release](https://github.com/UK2Group/ics-cli/releases/latest/download/ics-cli-windows-amd64.exe)
2. Rename to `ics-cli.exe` if desired
3. Add to a directory in your PATH or run directly

## Getting Started

### Authentication

Before using the ICS CLI, you need to authenticate:

```bash
# Set API credentials
ics-cli auth login
```

### First Steps

```bash
# View available commands
ics-cli --help

# List your baremetal servers
ics-cli baremetal list

# List your VPS servers
ics-cli vps list

# View version information
ics-cli version
```

## Usage Examples

### VPS Management

```bash
# List all VPS servers
ics-cli vps list

# List VPS servers grouped by datacenter
ics-cli vps list --display datacenter

# Get detailed information about a VPS
ics-cli vps get [ServiceID]

# Show root password when getting VPS details
ics-cli vps get [ServiceID] --password

# Power control
ics-cli vps poweron [ServiceID]
ics-cli vps poweroff [ServiceID]
ics-cli vps reboot [ServiceID]

# Cancel a VPS (IRREVERSIBLE - deletes all data)
ics-cli vps cancel [ServiceID]
```

### Deploying a New VPS

```bash
# View available VPS plans
ics-cli vps deploy list-plans

# View available datacenters
ics-cli vps deploy list-locations

# View available operating systems for a location
ics-cli vps deploy list-os "Chicago - G, US"

# Deploy a VPS with monthly billing
ics-cli vps deploy create --plan V2 --location "Chicago - G, US" --os "Ubuntu 22.04" --hostname web1 --friendlyname "Web Server 1"

# Deploy a VPS with hourly billing and backups
ics-cli vps deploy create --plan V4 --location "New York City - A, US" --os "Debian 11" --hostname db1 --bill-hourly --backups --friendlyname "Database Server"
```

### Baremetal Server Management

```bash
# List all servers
ics-cli baremetal list

# List servers at a specific site
ics-cli baremetal list --site NYC1

# Get detailed information about a server
ics-cli baremetal get [ServiceID]

# Power control
ics-cli baremetal poweron [ServiceID]
ics-cli baremetal poweroff [ServiceID]
ics-cli baremetal reboot [ServiceID]

# Access remote console
ics-cli baremetal ikvm [ServiceID]
ics-cli baremetal sol [ServiceID]
```

### Deploying a New Server

```bash
# View available inventory
ics-cli baremetal deploy list-inventory --datacenter NYC1

# Check available add-ons for a specific server type
ics-cli baremetal deploy list-addons --sku c1i.small --datacenter NYC1

# Order a new server (monthly billing)
ics-cli baremetal deploy create --sku c1i.small --datacenter NYC1 --os DEBIAN_11 --ssh-keys "My Key"

# Order a server with hourly billing
ics-cli baremetal deploy create --sku c1i.small --datacenter NYC1 --os DEBIAN_11 --bill-hourly

# Order multiple servers with hourly billing and additional options
ics-cli baremetal deploy create --sku c1.large --datacenter LHR1 --os UBUNTU_22_04 --quantity 2 --bill-hourly --ssh-keys "Production Key" --support BASICSUP
```

### SSH Key Management

```bash
# List all SSH keys
ics-cli sshkeys list

# Add a new SSH key
ics-cli sshkeys add --name "My New Key" --key "ssh-rsa AAAA..."

# Delete an SSH key
ics-cli sshkeys delete --name "My Old Key"

# Assign an SSH key to a server
ics-cli sshkeys assign --name "My Key" --server SERVER_ID
```

### Hourly Billing Management

```bash
# View inventory with hourly pricing information
ics-cli baremetal deploy list-inventory --hourly-enabled

# Filter by hourly price range
ics-cli baremetal deploy list-inventory --min-hourly-price 0.50 --max-hourly-price 2.00

# Cancel an hourly-billed server (IRREVERSIBLE)
ics-cli baremetal cancel [ServiceID]

# Cancel with force flag (no confirmation prompts)
ics-cli baremetal cancel [ServiceID] --force
```

## Environment Variables

The CLI supports the following environment variables:

| Variable | Description |
|----------|-------------|
| `ICS_API_KEY` | Your Ingenuity Cloud Services API key |
| `ICS_CONFIG_FILE` | Custom path to config file |

## Contributing

We welcome contributions! Please see our [Contributing Guidelines](CONTRIBUTING.md) for details.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

If you encounter any problems or have suggestions, please:

1. Check the [documentation](https://www.ingenuitycloudservices.com)
2. Search for [existing issues](https://github.com/UK2Group/ics-cli/issues)
3. Open a [new issue](https://github.com/UK2Group/ics-cli/issues/new) if needed

For paid support, please contact Ingenuity Cloud Services support through your account portal.

---

**Built with ❤️ by the Ingenuity Cloud Services Team**