# ICS CLI

A powerful command-line interface for managing Ingenuity Cloud Services resources.

[![Release](https://img.shields.io/github/v/release/UK2Group/ics-cli)](https://github.com/UK2Group/ics-cli/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/UK2Group/ics-cli)](https://goreportcard.com/report/github.com/UK2Group/ics-cli)
[![License](https://img.shields.io/github/license/UK2Group/ics-cli)](https://github.com/UK2Group/ics-cli/blob/main/LICENSE)

## Overview

ICS CLI provides a convenient way to manage your Ingenuity Cloud Services resources through a simple command-line interface. It allows you to provision and manage baremetal servers and more without needing to use the web interface.

## Features

- **Baremetal Server Management**: List, deploy, power control, and manage your dedicated servers
- **SSH Key Management**: Create, list, and manage SSH keys for server access
- **Inventory Browsing**: View available hardware in different datacenters
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

# View version information
ics-cli version
```

## Usage Examples

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
ics-cli baremetal list-inventory --datacenter NYC1

# Check available add-ons for a specific server type
ics-cli baremetal list-addons --sku c1i.small --datacenter NYC1

# Order a new server
ics-cli baremetal order --sku c1i.small --datacenter NYC1 --os DEBIAN_11 --ssh-keys "My Key"
```

### SSH Key Management

```bash
# List all SSH keys
ics-cli sshkey list

# Add a new SSH key
ics-cli sshkey add --name "My New Key" --key "ssh-rsa AAAA..."

# Delete an SSH key
ics-cli sshkey delete --name "My Old Key"

# Assign an SSH key to a server
ics-cli sshkey assign --name "My Key" --server SERVER_ID
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