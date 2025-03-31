package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all baremetal servers in your account",
	Run: func(cmd *cobra.Command, args []string) {

		// Make API call to fetch server list
		servers, err := getServerList()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error fetching server list:", err)
			return
		}

		// Get flags
		displayBySite, _ := cmd.Flags().GetBool("display")
		filterBySite, _ := cmd.Flags().GetString("site")

		// If no servers returned
		if len(servers) == 0 {
			fmt.Println("No servers found in your account.")
			return
		}

		headerFmt := color.New(color.FgBlue).SprintfFunc()
		columnFmt := color.New(color.FgYellow).SprintfFunc()

		// Create a map to group servers by datacenter if filtering by site
		if displayBySite {
			datacenterMap := make(map[string][]Server)
			for _, server := range servers {
				datacenterMap[server.DatacenterName] = append(datacenterMap[server.DatacenterName], server)
			}

			for datacenter, servers := range datacenterMap {
				fmt.Printf("\nDatacenter: %s\n", datacenter)

				// Build table header
				tbl := table.New("Service ID", "Hostname", "Primary IP", "Friendly Name")
				tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

				for _, server := range servers {

					tbl.AddRow(server.ServiceID,
						server.Hostname,
						server.PublicIP,
						server.FriendlyName)

				}

				tbl.Print()

			}
		} else {
			// Build table header
			tbl := table.New("Service ID", "Hostname", "Primary IP", "Datacenter", "Friendly Name")
			tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

			// Build each server in a row
			for _, server := range servers {
				// Use "contains" logic instead of exact matching
				if filterBySite != "" && !strings.Contains(strings.ToLower(server.DatacenterName), strings.ToLower(filterBySite)) {
					continue
				}

				tbl.AddRow(server.ServiceID,
					server.Hostname,
					server.PublicIP,
					server.DatacenterName,
					server.FriendlyName)

			}

			tbl.Print()
		}
	},
}

func init() {
	baremetalCmd.AddCommand(listCmd)

	listCmd.Flags().BoolP("display", "d", false, "Display servers grouped per site")
	listCmd.Flags().StringP("site", "s", "", "Filter servers by site")
}
