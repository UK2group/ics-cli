package cmd

import (
	"github.com/spf13/cobra"
)

// baremetalDeployCmd represents the baremetal deploy command
var baremetalDeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy a new Bare Metal server",
}

func init() {
	baremetalCmd.AddCommand(baremetalDeployCmd)
}
