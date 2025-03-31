package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
var sshKeysAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new SSH Key to the account",
	Example: `  
  # Add an SSH key from a file with a simple name
  ics-cli sshkeys add --name MyKey --file ~/.ssh/id_rsa.pub
  
  # Add an SSH key directly with a multi-word name
  ics-cli sshkeys add --name "My Production Key" --key "ssh-rsa AAAAB3..."`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get command flags
		name, _ := cmd.Flags().GetString("name")
		key, _ := cmd.Flags().GetString("key")
		filePath, _ := cmd.Flags().GetString("file")

		// Validate required parameters
		if name == "" {
			fmt.Fprintln(os.Stderr, "Error: SSH Key name is required (--name)")
			return
		}

		// Get the SSH key from either file or direct input
		var sshKey string
		if filePath != "" {
			// Read from file
			keyData, err := os.ReadFile(filePath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading SSH key file: %s\n", err)
				return
			}
			sshKey = string(keyData)
		} else if key != "" {
			// Use direct input
			sshKey = key
		} else {
			fmt.Fprintln(os.Stderr, "Error: Either SSH key (--key) or key file (--file) is required")
			return
		}

		// Clean up the key (remove extra whitespace, comments, etc.)
		sshKey = cleanSSHKey(sshKey)

		// Create request body
		requestData := map[string]string{
			"label":      name,
			"public_key": sshKey,
		}

		// Add the SSH key
		id, err := addSSHKey(requestData)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error adding SSH key: %s\n", err)
			return
		}

		if id {
			fmt.Printf("%s\n", BlueHeading("Successfully added SSH Key"))
		}
	},
}

func init() {
	sshkeysCmd.AddCommand(sshKeysAddCmd)

	sshKeysAddCmd.Flags().StringP("name", "n", "", "Name of the SSH Key")
	sshKeysAddCmd.Flags().StringP("key", "k", "", "SSH Key")
	sshKeysAddCmd.Flags().StringP("file", "f", "", "File containing the SSH Key")

	sshKeysAddCmd.MarkFlagRequired("name")
}
