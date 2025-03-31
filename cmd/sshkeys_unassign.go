package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var sshKeyUnassignCmd = &cobra.Command{
	Use:   "unassign [sshKeyName]",
	Short: "Unassign an existing SSH Key from a Baremetal Server",
	Example: `
  # Unassign an existing SSH Key from a Baremetal Server with Service ID 123456
  ics-cli sshkeys unassign my-ssh-key --server 123456`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		// Get the SSH Key Name from the arguments
		sshKeyName := args[0]

		// Check if Server ID is provided
		serviceID, _ := cmd.Flags().GetString("server")
		if serviceID == "" {
			fmt.Fprintln(os.Stderr, "Error: Service ID must be provided with --server or -s flag")
			return
		}

		// Get the Server ID from Service ID
		serverID, err := getServerIDFromServiceID(serviceID)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

		// Get the SSH Key from the Label
		sshKey, err := getSSHKeyFromLabel(sshKeyName)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

		// Unassign the SSH Key to the Server
		unassignKey, err := unassignSSHKey(serverID, sshKey.ID)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error unassigning SSH Key:", err)
			return
		}

		if unassignKey {
			fmt.Printf("%s\n", BlueHeading("Successfully unassigned SSH Key from server. The SSH Key will be removed after the next reinstall."))
		}
	},
}

func init() {
	sshkeysCmd.AddCommand(sshKeyUnassignCmd)

	sshKeyUnassignCmd.Flags().StringP("server", "s", "", "Service ID to unassign the SSH Key from")
	sshKeyUnassignCmd.MarkFlagRequired("server")

}
