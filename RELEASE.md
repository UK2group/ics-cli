# Release Process

This document describes how to create a new release of ICS CLI.

## Prerequisites

- All changes merged to `main` branch
- All tests passing
- Version bumped in `cmd/version.go`

## Release Steps

### 1. Update Version

Edit `cmd/version.go` and update the version number:

```go
var (
    Version   = "1.0.3"  // <- Update this
    BuildTime = "unknown"
    Commit    = "unknown"
)
```

### 2. Commit and Push Version Bump

```bash
git add cmd/version.go
git commit -m "Bump version to 1.0.3"
git push origin main
```

### 3. Create and Push Git Tag

**IMPORTANT:** The tag must match the version in `cmd/version.go`

```bash
# Create tag (must start with 'v')
git tag v1.0.3

# Push tag to trigger release workflow
git push origin v1.0.3
```

This will automatically:
- Build binaries for macOS (amd64/arm64), Linux (amd64), and Windows (amd64)
- Generate SHA256 checksums
- Create a GitHub release with all binaries

### 4. Get SHA256 Checksums

Once the release is published, download the checksums:

```bash
# Option 1: From GitHub release page
# Go to: https://github.com/UK2Group/ics-cli/releases/latest
# Download SHA256SUMS.txt

# Option 2: Using curl
curl -L https://github.com/UK2Group/ics-cli/releases/download/v1.0.3/SHA256SUMS.txt
```

You'll need the checksums for:
- `ics-cli-macos-amd64`
- `ics-cli-macos-arm64`

### 5. Update Homebrew Formula

Navigate to the Homebrew tap repository:

```bash
cd /path/to/homebrew-ics-cli
# Or clone it:
# git clone https://github.com/UK2Group/homebrew-ics-cli
```

Edit the formula file (usually `Formula/ics-cli.rb` or `ics-cli.rb`):

```ruby
class IcsCli < Formula
  desc "CLI for Ingenuity Cloud Services"
  homepage "https://github.com/UK2Group/ics-cli"
  version "1.0.3"  # <- Update this

  if Hardware::CPU.arm?
    url "https://github.com/UK2Group/ics-cli/releases/download/v1.0.3/ics-cli-macos-arm64"
    sha256 "PASTE_ARM64_SHA256_HERE"  # <- Update this
  else
    url "https://github.com/UK2Group/ics-cli/releases/download/v1.0.3/ics-cli-macos-amd64"
    sha256 "PASTE_AMD64_SHA256_HERE"  # <- Update this
  end

  def install
    bin.install "ics-cli-macos-#{Hardware::CPU.arch}" => "ics-cli"
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/ics-cli version")
  end
end
```

### 6. Commit and Push Homebrew Formula

```bash
cd /path/to/homebrew-ics-cli
git add Formula/ics-cli.rb  # or ics-cli.rb
git commit -m "Update ics-cli to v1.0.3"
git push origin main
```

### 7. Test the Update

```bash
# Update Homebrew
brew update

# Upgrade ics-cli
brew upgrade ics-cli

# Verify version
ics-cli version
# Should output: 1.0.3-<commit> (YYYY-MM-DD HH:MM:SS)
```

## Quick Reference

```bash
# Full release workflow
git add cmd/version.go
git commit -m "Bump version to 1.0.3"
git push
git tag v1.0.3
git push origin v1.0.3

# Wait for GitHub Action to complete, then:
# 1. Get SHA256 from: https://github.com/UK2Group/ics-cli/releases/latest
# 2. Update homebrew-ics-cli repository with new version and SHA256
# 3. Test: brew update && brew upgrade ics-cli
```

## Troubleshooting

### Version shows "unknown"
- The ldflags in `.github/workflows/release.yaml` must use the full package path
- Should be: `github.com/UK2Group/ics-cli/cmd.Version`

### Homebrew doesn't see new version
- Make sure you pushed changes to the homebrew-ics-cli repository
- Run `brew update` to refresh Homebrew's formulae
- Check formula location: `brew formula ics-cli`

### Release workflow didn't trigger
- Tag must start with 'v' (e.g., `v1.0.3`)
- Check GitHub Actions tab for errors
- Verify tag was pushed: `git ls-remote --tags origin`
