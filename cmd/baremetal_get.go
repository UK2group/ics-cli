package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

// bmstatusCmd represents the bare metal status command
var bmstatusCmd = &cobra.Command{
	Use:     "get [serviceID]",
	Aliases: []string{"status"},
	Short:   "Get the information of a Baremetal Server by its Service ID",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		// Get service ID from arguments
		serviceID := args[0]

		// Validate service ID is a number
		if _, err := strconv.Atoi(serviceID); err != nil {
			fmt.Fprintln(os.Stderr, "Error: Service ID must be a number")
			return
		}

		// Step 1: Get the server ID from service ID
		serverID, err := getServerIDFromServiceID(serviceID)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

		// Step 2: Get server details
		server, err := getServerDetails(serverID)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

		printServerDetails(cmd, server)
	},
}

func init() {
	baremetalCmd.AddCommand(bmstatusCmd)

	// Boolean flag to display root password
	bmstatusCmd.Flags().BoolP("password", "p", false, "Display root password")
	bmstatusCmd.Flags().BoolP("power", "s", false, "Ignore obtaining power status")
}
