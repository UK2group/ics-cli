package cmd

import (
	"github.com/spf13/cobra"
)

// sshkeysCmd represents the sshkeys command
var sshkeysCmd = &cobra.Command{
	Use:   "sshkeys",
	Short: "Manage SSH Keys for your account",
}

func init() {
	rootCmd.AddCommand(sshkeysCmd)
}
