package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var sshKeyRenameCmd = &cobra.Command{
	Use:   "rename [sshKeyName]",
	Short: "Rename an existing SSH Key",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		// Get the SSH Key Name from the arguments
		sshKeyName := args[0]

		// Check if URL is provided
		newName, _ := cmd.Flags().GetString("name")
		if newName == "" {
			fmt.Fprintln(os.Stderr, "Error: New SSH Key Name required with --name or -n flag")
			return
		}

		// Get the SSH Key from the Label
		sshKey, err := getSSHKeyFromLabel(sshKeyName)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

		// Update SSH Key Label
		setLabel, err := setSSHKeyLabel(sshKey.ID, newName)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error renaming SSH Key:", err)
			return
		}

		if setLabel {
			fmt.Printf("%s\n", BlueHeading("Successfully renamed SSH Key"))
		}
	},
}

func init() {
	sshkeysCmd.AddCommand(sshKeyRenameCmd)

	sshKeyRenameCmd.Flags().StringP("name", "n", "", "New name for SSH Key")
	sshKeysAddCmd.MarkFlagRequired("name")

}
