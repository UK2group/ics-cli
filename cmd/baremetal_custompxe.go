package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

// custompxeCmd represents the custompxe command
var custompxeCmd = &cobra.Command{
	Use:   "custompxe [serviceID]",
	Short: "Specify a custom PXE boot URL to load your preferred network boot environment",
	Long: `Specify a custom PXE boot URL to load your preferred network boot environment.
This allows you to deploy custom operating systems, recovery tools, or provisioning scripts tailored to your needs.
Ensure the URL points to a valid PXE boot server with the necessary configurations.
Your server will be rebooted to apply the custom PXE boot.
Incorrect configurations may result in boot failures. Ensure your PXE setup is properly configured before proceeding.`,
	Args:    cobra.ExactArgs(1),
	Aliases: []string{"pxe"},
	Run: func(cmd *cobra.Command, args []string) {

		// Get service ID from arguments
		serviceID := args[0]

		// Validate service ID is a number
		if _, err := strconv.Atoi(serviceID); err != nil {
			fmt.Fprintln(os.Stderr, "Error: Service ID must be a number")
			return
		}

		// Check if URL is provided
		url, _ := cmd.Flags().GetString("url")
		if url == "" {
			fmt.Fprintln(os.Stderr, "Error: URL is required with --url or -u flag")
			return
		}

		// Step 1: Get the server ID from service ID
		serverID, err := getServerIDFromServiceID(serviceID)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

		// Set the custom PXE URL
		customPXE, err := setPXEUrl(serverID, url)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error setting PXE URL:", err)
			return
		}

		if customPXE {
			fmt.Printf("%s\n", BlueHeading("Successfully updated the server with custom PXE URL and requested a reboot"))
		}

	},
}

func init() {
	baremetalCmd.AddCommand(custompxeCmd)

	custompxeCmd.Flags().StringP("url", "u", "", "URL to set as the custom PXE boot")
	sshKeysAddCmd.MarkFlagRequired("url")
}
