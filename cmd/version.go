package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Version   = "1.0.1"
	BuildTime = "unknown"
	Commit    = "unknown"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of ICS CLI",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s-%s (%s)", Version, Commit, BuildTime)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
