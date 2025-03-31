package cmd

import (
	"github.com/spf13/cobra"
)

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Login / Logout to your Ingenuity Cloud Services Acccount",
}

func init() {
	rootCmd.AddCommand(authCmd)
}
