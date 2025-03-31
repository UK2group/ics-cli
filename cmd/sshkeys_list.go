package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var sshkeylistCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "Get a list of all SSH keys in your account",
	Run: func(cmd *cobra.Command, args []string) {
		// Make API call to get SSH keys
		sshKeys, err := getSSHKeys()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return
		}

		// If no ssh keys returned
		if len(sshKeys) == 0 {
			fmt.Println("No SSH Keys found in your account.")
			return
		}

		headerFmt := color.New(color.FgBlue).SprintfFunc()
		columnFmt := color.New(color.FgYellow).SprintfFunc()

		// Build table header
		tbl := table.New("Key Name", "Created At", "Updated At", "Assigned To Servers")
		tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

		// Build each server in a row
		for _, sshKey := range sshKeys {

			tbl.AddRow(sshKey.Label,
				time.Unix(sshKey.CreatedAt, 0).Format("2006-01-02 15:04:05"),
				time.Unix(sshKey.UpdatedAt, 0).Format("2006-01-02 15:04:05"),
				len(sshKey.AssignedServers))

		}

		tbl.Print()
	},
}

func init() {
	sshkeysCmd.AddCommand(sshkeylistCmd)
}
