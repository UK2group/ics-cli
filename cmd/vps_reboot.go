package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

// vpsRebootCmd represents the reboot command
var vpsRebootCmd = &cobra.Command{
	Use:   "reboot [serviceID]",
	Short: "Reboot a VPS server",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		// Get service ID from arguments
		serviceID := args[0]

		// Validate service ID is a number
		if _, err := strconv.Atoi(serviceID); err != nil {
			fmt.Fprintln(os.Stderr, "Error: Service ID must be a number")
			return
		}

		// Step 1: Get the VPS ID and location ID from service ID
		vpsID, locationID, err := getVPSIDAndLocationFromServiceID(serviceID)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

		// Check if user wants to proceed
		dontPrompt, _ := cmd.Flags().GetBool("dont")
		if !dontPrompt {
			fmt.Print("Are you sure you want to reboot the VPS server? (y/n): ")
			var response string
			fmt.Scanln(&response)
			if response != "y" {
				return
			}
		}

		// Reboot the VPS
		rebootVPS, err := setVPSReboot(vpsID, locationID)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error sending reboot command:", err)
			return
		}

		if rebootVPS {
			fmt.Printf("%s\n", BlueHeading("Successfully rebooted the VPS server."))
		} else {
			fmt.Printf("%s\n", RedText("Failed to reboot the VPS server. Please try again."))
		}

	},
}

func init() {
	vpsCmd.AddCommand(vpsRebootCmd)

	vpsRebootCmd.Flags().BoolP("dont", "d", false, "Don't prompt for confirmation")
}
