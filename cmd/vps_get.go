package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

// vpsGetCmd represents the get command
var vpsGetCmd = &cobra.Command{
	Use:     "get [serviceID]",
	Aliases: []string{"status"},
	Short:   "Get the information of a VPS Server by its Service ID",
	Args:    cobra.ExactArgs(1),
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

		// Step 2: Get VPS details
		vps, err := getVPSDetails(vpsID, locationID)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

		printVPSDetails(cmd, vps)
	},
}

func init() {
	vpsCmd.AddCommand(vpsGetCmd)

	// Boolean flag to display root password
	vpsGetCmd.Flags().BoolP("password", "p", false, "Display root password")
	vpsGetCmd.Flags().BoolP("power", "s", false, "Ignore obtaining power status")
}
