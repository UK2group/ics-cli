package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// vpsCancelCmd represents the cancel command
var vpsCancelCmd = &cobra.Command{
	Use:   "cancel [serviceID]",
	Short: "Cancel a VPS server",
	Long: `Cancel a VPS server by its Service ID.

WARNING: This is irreversible and all data will be permanently deleted.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serviceID := args[0]
		force, _ := cmd.Flags().GetBool("force")

		// Get the VPS ID and location ID from service ID
		vpsID, locationID, err := getVPSIDAndLocationFromServiceID(serviceID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return
		}

		if !force {
			fmt.Println(YellowText("\nWARNING: This will permanently delete the VPS and all its data!"))
			fmt.Print("Are you sure you want to cancel this VPS? (y/N): ")

			var response string
			fmt.Scanln(&response)

			if response != "y" && response != "Y" && response != "yes" && response != "YES" {
				fmt.Println("Cancellation aborted")
				return
			}
		}

		// Cancel the VPS
		response, err := cancelVPS(vpsID, locationID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error canceling VPS: %v\n", err)
			return
		}

		if response.StatusCode == 200 {
			fmt.Printf("\n%s\n", GreenText("✓ VPS cancellation successful!"))
			fmt.Printf("%s\n", WhiteText(response.Message))
		} else {
			fmt.Printf("\n%s\n", YellowText("VPS cancellation may have issues"))
			fmt.Printf("%s\n", WhiteText(response.Message))
		}
	},
}

func init() {
	vpsCmd.AddCommand(vpsCancelCmd)

	// Add flags
	vpsCancelCmd.Flags().Bool("force", false, "Skip confirmation prompt")
}
