package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var sshKeyDeleteCmd = &cobra.Command{
	Use:   "delete [sshKeyName]",
	Short: "Delete an existing SSH Key",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Get the SSH Key Name from the arguments
		sshKeyName := args[0]

		// Get the SSH Key from the Label
		sshKey, err := getSSHKeyFromLabel(sshKeyName)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

		// Update SSH Key Label
		deleteKey, err := deleteSSHKey(sshKey.ID)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error deleting SSH Key:", err)
			return
		}

		if deleteKey {
			fmt.Printf("%s\n", BlueHeading("Successfully deleted SSH Key"))
		}
	},
}

func init() {
	sshkeysCmd.AddCommand(sshKeyDeleteCmd)
}
