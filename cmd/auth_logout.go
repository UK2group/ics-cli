package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// logoutCmd represents the logout command
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout of the Ingenuity Cloud Services API",
	Run: func(cmd *cobra.Command, args []string) {
		// Check if there is an API key to delete
		if viper.IsSet("api_key") {
			// Remove API key from configuration
			viper.Set("api_key", "")

			// Save the updated configuration
			if err := viper.WriteConfig(); err != nil {
				fmt.Println("Error removing API key from configuration:", err)
				return
			}
			fmt.Println("Successfully logged out from Ingenuity Cloud Services API")
		} else {
			fmt.Println("You are not currently logged in")
		}
	},
}

func init() {
	authCmd.AddCommand(logoutCmd)
}
