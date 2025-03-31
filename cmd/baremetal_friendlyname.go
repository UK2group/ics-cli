package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

// bmNameCmd represents the name command
var bmNameCmd = &cobra.Command{
	Use:   "friendlyname [serviceID]",
	Short: "Set a friendly name for a server",
	Example: `  
	# Set a friendly name for a server
	ics-cli baremetal friendlyname 123456 --name "Production Server 1"`,
	Args: cobra.ExactArgs(1),
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

		// Get the friendly name
		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			fmt.Fprintln(os.Stderr, "Error: Friendly name is required with --name or -n flag")
			return
		}

		// Set the friendly name
		customPXE, err := setFriendlyName(serverID, name)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error setting Friendly name:", err)
			return
		}

		if customPXE {
			fmt.Printf("%s\n", BlueHeading("Successfully updated the friendly name of the server."))
		}
	},
}

func init() {
	baremetalCmd.AddCommand(bmNameCmd)

	// Add flag to open SOL in browser
	bmNameCmd.Flags().StringP("name", "n", "", "Friendly name for the server")
}
