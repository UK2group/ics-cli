package cmd

import (
	"github.com/spf13/cobra"
)

// baremetalCmd represents the baremetal command
var baremetalCmd = &cobra.Command{
	Use:     "baremetal",
	Aliases: []string{"bm", "metal"},
	Short:   "Create, Destroy and Manage Baremetal Servers",
}

func init() {
	rootCmd.AddCommand(baremetalCmd)
}
