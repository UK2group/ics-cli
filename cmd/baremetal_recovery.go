package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

// recoveryCmd represents the reboot command
var recoveryCmd = &cobra.Command{
	Use:   "recovery [serviceID]",
	Short: "Boot a baremetal server into a recovery image",
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
			fmt.Print("Are you sure you want to boot the server into a recovery image? (y/n): ")
			var response string
			fmt.Scanln(&response)
			if response != "y" {
				return
			}
		}

		// Reboot the server
		rebootServer, err := setRecovery(serverID)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error sending recovery command:", err)
			return
		}

		if rebootServer {
			fmt.Printf("%s\n", BlueHeading("Successfully booted the server into a recovery image."))
		} else {
			fmt.Printf("%s\n", RedText("Failed to boot the server into a recovery image. Please try again."))
		}

	},
}

func init() {
	baremetalCmd.AddCommand(recoveryCmd)

	recoveryCmd.Flags().BoolP("dont", "d", false, "Don't prompt for confirmation")
}
