package cmd

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getsshkeysCmd = &cobra.Command{
	Use:   "get [sshKeyName]",
	Short: "Get SSH key details by Name",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Get the SSH Key Name from the arguments
		sshKeyName := args[0]

		// Get the SSH Key from the Label
		sshKey, err := getSSHKeyFromLabel(sshKeyName)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

		// Display SSH Key details
		fmt.Printf("%s %s\n", BlueHeading("Name:"), WhiteText(sshKey.Label))
		fmt.Printf("%s %s\n", BlueHeading("Created At:"), WhiteText(time.Unix(sshKey.CreatedAt, 0).Format("2006-01-02 15:04:05")))
		fmt.Printf("%s %s\n", BlueHeading("Updated At:"), WhiteText(time.Unix(sshKey.UpdatedAt, 0).Format("2006-01-02 15:04:05")))
		fmt.Printf("%s %s\n", BlueHeading("Assigned to servers:"), WhiteText(strconv.Itoa(len(sshKey.AssignedServers))))
		fmt.Printf("%s \n%s\n", BlueHeading("SSH Key:"), WhiteText(sshKey.Key))

	},
}

func init() {
	sshkeysCmd.AddCommand(getsshkeysCmd)
	getsshkeysCmd.Flags().StringP("name", "n", "", "Name of the SSK Key")
}
