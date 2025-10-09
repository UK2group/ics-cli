package cmd

import (
	"testing"
)

func TestCleanSSHKey(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "clean RSA key",
			input: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC user@host",
			want:  "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC user@host",
		},
		{
			name:  "RSA key with leading/trailing whitespace",
			input: "  ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC user@host  \n",
			want:  "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC user@host",
		},
		{
			name:  "ED25519 key",
			input: "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIAbCdEfGhIjKlMnOpQrStUvWxYz user@host",
			want:  "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIAbCdEfGhIjKlMnOpQrStUvWxYz user@host",
		},
		{
			name:  "ECDSA key",
			input: "ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTY user@host",
			want:  "ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTY user@host",
		},
		{
			name: "key with multiple lines and comments",
			input: `# This is a comment
ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC user@host
# Another comment`,
			want: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC user@host",
		},
		{
			name: "key with extra newlines",
			input: `

ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC user@host

`,
			want: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC user@host",
		},
		{
			name:  "DSS key",
			input: "ssh-dss AAAAB3NzaC1kc3MAAACBAP user@host",
			want:  "ssh-dss AAAAB3NzaC1kc3MAAACBAP user@host",
		},
		{
			name:  "invalid key format returns trimmed input",
			input: "  not-a-valid-ssh-key  ",
			want:  "not-a-valid-ssh-key",
		},
		{
			name:  "empty string",
			input: "   \n  \n  ",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cleanSSHKey(tt.input)
			if got != tt.want {
				t.Errorf("cleanSSHKey() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestCleanSSHKeyPreservesKeyTypes(t *testing.T) {
	keyTypes := []string{
		"ssh-rsa",
		"ssh-ed25519",
		"ssh-dss",
		"ecdsa-sha2-nistp256",
		"ecdsa-sha2-nistp384",
		"ecdsa-sha2-nistp521",
	}

	for _, keyType := range keyTypes {
		t.Run(keyType, func(t *testing.T) {
			input := keyType + " AAAAB3NzaC1 user@host"
			result := cleanSSHKey(input)

			if result != input {
				t.Errorf("cleanSSHKey() modified clean key: got %q, want %q", result, input)
			}
		})
	}
}
