package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

// ikvmCmd represents the ikvm command
var ikvmCmd = &cobra.Command{
	Use:   "ikvm [serviceID]",
	Short: "Generate an IPMI IKVM Console URL",
	Args:  cobra.ExactArgs(1),
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

		// Get the iKVM access link
		solLink, err := getIkvmResponse(serverID)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error getting iKVM access:", err)
			return
		}

		// Print the SOL access link
		fmt.Printf("%s %s\n", BlueHeading("iKVM Access Link:"), WhiteText(solLink))

		// If --open flag is specified, try to open the link in a browser
		openInBrowser, _ := cmd.Flags().GetBool("browser")
		if !openInBrowser {
			fmt.Println("Opening iKVM link in your default browser...")
			err := openBrowser(solLink)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error opening browser:", err)
			}
		}
	},
}

func init() {
	baremetalCmd.AddCommand(ikvmCmd)

	// Add flag to open iKVM in browser
	ikvmCmd.Flags().BoolP("browser", "d", false, "Don't open iKVM access in default browser")
}
