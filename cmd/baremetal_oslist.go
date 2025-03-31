package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/fatih/color"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
)

// bmNameCmd represents the name command
var bmOSList = &cobra.Command{
	Use:   "oslist [serviceID]",
	Short: "Get a list of operating systems available for a server",
	Example: `  
	# Get a list of operating systems available for a server with service ID 123456
	ics-cli baremetal oslist 123456`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		// Get service ID from arguments
		serviceID := args[0]

		// Validate service ID is a number
		if _, err := strconv.Atoi(serviceID); err != nil {
			fmt.Fprintln(os.Stderr, "Error: Service ID must be a number")
			return
		}

		// Step 1: Get the server ID from service ID
		serverID, err := getServerIDFromServiceID(serviceID)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

		// Get the OS List
		osList, err := getOSList(serverID)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error getting OS List:", err)
			return
		}

		// If no ssh keys returned
		if len(osList) == 0 {
			fmt.Println("No available Operating Systems found.")
			return
		}

		headerFmt := color.New(color.FgBlue).SprintfFunc()
		columnFmt := color.New(color.FgYellow).SprintfFunc()

		// Build table header
		tbl := table.New("ID", "Name", "License Required")
		tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

		// Build each OS in a row
		for _, os := range osList {
			// Determine if license is required
			licenseRequired := "No"
			if len(os.Licenses) > 0 && os.Licenses[0] != "" {
				licenseRequired = "Yes"
			}

			tbl.AddRow(os.ID, os.Name, licenseRequired)
		}

		tbl.Print()
	},
}

func init() {
	baremetalCmd.AddCommand(bmOSList)
}
