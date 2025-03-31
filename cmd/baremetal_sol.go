package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

// solCmd represents the sol command
var solCmd = &cobra.Command{
	Use:   "sol [serviceID]",
	Short: "Generate a Serial Over Lan (SOL) session",
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

		// Get the SOL access link
		solLink, err := getSolResponse(serverID)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error getting SOL access link:", err)
			return
		}

		// Print the SOL access link
		fmt.Printf("%s %s\n", BlueHeading("SOL Access Link:"), WhiteText(solLink))

		// If --open flag is specified, try to open the link in a browser
		openInBrowser, _ := cmd.Flags().GetBool("browser")
		if !openInBrowser {
			fmt.Println("Opening SOL access link in your default browser...")
			err := openBrowser(solLink)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error opening browser:", err)
			}
		}
	},
}

func init() {
	baremetalCmd.AddCommand(solCmd)

	// Add flag to open SOL in browser
	solCmd.Flags().BoolP("browser", "d", false, "Don't open SOL access in default browser")
}
