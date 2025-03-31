package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

// poweroffCmd represents the poweroff command
var poweroffCmd = &cobra.Command{
	Use:   "poweroff [serviceID]",
	Short: "Power off a baremetal server",
	Args:  cobra.ExactArgs(1),
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

		// Check if user wants to proceed
		dontPrompt, _ := cmd.Flags().GetBool("dont")
		if !dontPrompt {
			fmt.Print("Are you sure you want to power off the server? (y/n): ")
			var response string
			fmt.Scanln(&response)
			if response != "y" {
				return
			}
		}

		// Power off the server
		powerOff, err := setPowerOff(serverID)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error sending power command:", err)
			return
		}

		if powerOff {
			fmt.Printf("%s\n", BlueHeading("Successfully powered off the server."))
		} else {
			fmt.Printf("%s\n", RedText("Failed to power off the server. Please try again."))
		}

	},
}

func init() {
	baremetalCmd.AddCommand(poweroffCmd)

	poweroffCmd.Flags().BoolP("dont", "d", false, "Don't prompt for confirmation")
}
