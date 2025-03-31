package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

// bmReinstall represents the reinstall
var bmReinstallCmd = &cobra.Command{
	Use:   "reinstall [serviceID]",
	Short: "Reinstall the operation system on a Baremetal Server",
	Example: `  
  # Reinstall Ubuntu 24.04 on a Baremetal Server with service ID 123456
  ics-cli baremetal reinstall 123456 --operatingsystem ubuntu-24-04 --reason "Reinstalling OS"`,
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

		imageID, _ := cmd.Flags().GetString("os")
		if imageID == "" {
			fmt.Fprintln(os.Stderr, "Error: Operating System ID is required (use --os or -o)")
			return
		}

		reason, _ := cmd.Flags().GetString("reason")

		// Check if user wants to proceed
		dontPrompt, _ := cmd.Flags().GetBool("dont")
		if !dontPrompt {
			fmt.Print("Are you sure you want to perform an OS reinstall on the server? Re-enter the Service ID to confirm: ")
			var response string
			fmt.Scanln(&response)
			if response != serviceID {
				fmt.Println("Service ID does not match. Aborting.")
				return
			}
			// Extra confirmation for destructive actions
			fmt.Print("This action is irreversible and will erase all data on the target server. Are you sure? (y/n): ")
			var response2 string
			fmt.Scanln(&response2)
			if response2 != "y" {
				fmt.Println("No user confirmation. Aborting.")
				return
			}
		}

		// Reinstall the server
		reinstall, err := setReinstallOS(serverID, reason, imageID)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error sending reinstall command:", err)
			return
		}

		if reinstall {
			fmt.Printf("%s\n", BlueHeading("\nSuccessfully started a reinstall on the server. Please allow 10-15 minutes for the server to be reinstalled."))
		} else {
			fmt.Printf("%s\n", RedText("\nFailed to reinstall the server. Please try again."))
		}
	},
}

func init() {
	baremetalCmd.AddCommand(bmReinstallCmd)

	bmReinstallCmd.Flags().StringP("os", "o", "", "Operaring System ID to reinstall (Use the oslist command to get the image ID)")
	bmReinstallCmd.Flags().StringP("reason", "r", "", "(Optional) Reason for reinstalling the OS")
	bmReinstallCmd.Flags().BoolP("dont", "d", false, "Don't prompt for confirmation")
	bmReinstallCmd.MarkFlagRequired("image")
}
