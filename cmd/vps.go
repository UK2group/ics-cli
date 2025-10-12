package cmd

import (
	"github.com/spf13/cobra"
)

// vpsCmd represents the vps command
var vpsCmd = &cobra.Command{
	Use:   "vps",
	Short: "Create, Destroy and Manage VPS Servers",
}

func init() {
	rootCmd.AddCommand(vpsCmd)
}
