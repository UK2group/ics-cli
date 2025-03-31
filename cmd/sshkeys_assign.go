package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var sshKeyAssignCmd = &cobra.Command{
	Use:   "assign [sshKeyName]",
	Short: "Assign an existing SSH Key to a Baremetal Server",
	Example: `
  # Assign an existing SSH Key to a Baremetal Server with Service ID 123456
  ics-cli sshkeys assign my-ssh-key --server 123456`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		// Get the SSH Key Name from the arguments
		sshKeyName := args[0]

		// Check if Service ID is provided
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

		// Get the SSH Key ID from the Label
		sshKey, err := getSSHKeyFromLabel(sshKeyName)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

		// Get existing SSH keys assigned to the server
		// Convert serverID to int, ICS API returns a string
		serverIDInt, err := strconv.Atoi(serverID)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error converting server ID to int:", err)
			return
		}
		existingKeys, _ := getServerSSHKeys(serverIDInt)

		// Convert serverIDInt to a list of IDs
		var sshKeyIDs []int
		for _, key := range existingKeys {
			sshKeyIDs = append(sshKeyIDs, key.ID)
		}

		// Append the new SSH Key ID to the existing keys
		sshKeyIDs = append(sshKeyIDs, sshKey.ID)

		// Assign the SSH Key to the Server
		assignKey, err := assignSSHKeys(serviceID, sshKeyIDs)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error assigning SSH Key:", err)
			return
		}

		if assignKey {
			fmt.Printf("%s\n", BlueHeading("Successfully assigned SSH Key to server. The SSH Key will be available after the next reinstall."))
		}
	},
}

func init() {
	sshkeysCmd.AddCommand(sshKeyAssignCmd)

	sshKeyAssignCmd.Flags().StringP("server", "s", "", "Service ID to assign the SSH Key to")
	sshKeyAssignCmd.MarkFlagRequired("server")

}
