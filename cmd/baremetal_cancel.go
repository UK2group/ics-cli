package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

// bmCancelCmd represents the cancel command for hourly billed servers
var bmCancelCmd = &cobra.Command{
	Use:   "cancel [serviceID]",
	Short: "Cancel an hourly billed Baremetal Server",
	Long: `Cancel an hourly billed Baremetal Server. This command will immediately terminate
the server and stop hourly billing. This action is irreversible.

Note: This command only works for servers with hourly billing enabled.
Monthly billed servers cannot be cancelled through this command.`,
	Example: `  
  # Cancel an hourly billed server with service ID 123456
  ics-cli baremetal cancel 123456

  # Cancel without confirmation prompt
  ics-cli baremetal cancel 123456 --force`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		// Get service ID from arguments
		serviceID := args[0]

		// Validate service ID is a number
		if _, err := strconv.Atoi(serviceID); err != nil {
			fmt.Fprintln(os.Stderr, "Error: Service ID must be a number")
			return
		}

		// Get the server ID from service ID
		serverID, err := getServerIDFromServiceID(serviceID)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

		// Get server details to check if it's hourly billed
		serverDetails, err := getServerDetails(serverID)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error getting server details:", err)
			return
		}

		// Check if server is hourly billed
		if !serverDetails.HourlyBilling {
			fmt.Fprintln(os.Stderr, "Error: This server is not hourly billed. Only hourly billed servers can be cancelled.")
			return
		}

		// Check if user wants to proceed
		force, _ := cmd.Flags().GetBool("force")
		if !force {
			fmt.Printf("Server Details:\n")
			fmt.Printf("  Service ID: %s\n", serviceID)
			fmt.Printf("  Hostname: %s\n", serverDetails.Hostname)
			fmt.Printf("  Friendly Name: %s\n", serverDetails.FriendlyName)
			fmt.Printf("  Billing Type: Hourly\n")
			fmt.Printf("  Datacenter: %s\n", serverDetails.DatacenterName)
			
			fmt.Print("\nAre you sure you want to cancel this hourly billed server? Re-enter the Service ID to confirm: ")
			var response string
			fmt.Scanln(&response)
			if response != serviceID {
				fmt.Println("Service ID does not match. Aborting.")
				return
			}
			// Extra confirmation for destructive actions
			fmt.Print("This action is irreversible and will immediately terminate the server and stop billing. Are you sure? (y/n): ")
			var response2 string
			fmt.Scanln(&response2)
			if response2 != "y" {
				fmt.Println("No user confirmation. Aborting.")
				return
			}
		}

		// Cancel the server
		cancelled, err := cancelHourlyServer(serverID)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error cancelling server:", err)
			return
		}

		if cancelled {
			fmt.Printf("%s\n", GreenText("\nSuccessfully cancelled the hourly billed server. Billing has been stopped."))
		} else {
			fmt.Printf("%s\n", RedText("\nFailed to cancel the server. Please try again or contact support."))
		}
	},
}

func init() {
	baremetalCmd.AddCommand(bmCancelCmd)

	bmCancelCmd.Flags().BoolP("force", "f", false, "Don't prompt for confirmation")
}