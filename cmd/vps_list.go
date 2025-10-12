package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
)

// vpsListCmd represents the list command
var vpsListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all VPS servers in your account",
	Run: func(cmd *cobra.Command, args []string) {

		// Make API call to fetch VPS list
		vpsList, err := getVPSList()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error fetching VPS list:", err)
			return
		}

		// Get flags
		displayBySite, _ := cmd.Flags().GetBool("display")
		filterBySite, _ := cmd.Flags().GetString("site")

		// If no VPS servers returned
		if len(vpsList) == 0 {
			fmt.Println("No VPS servers found in your account.")
			return
		}

		headerFmt := color.New(color.FgBlue).SprintfFunc()
		columnFmt := color.New(color.FgYellow).SprintfFunc()

		// Create a map to group VPS by datacenter if filtering by site
		if displayBySite {
			datacenterMap := make(map[string][]VPSBasic)
			for _, vps := range vpsList {
				datacenterMap[vps.Location.Name] = append(datacenterMap[vps.Location.Name], vps)
			}

			for datacenter, vpsServers := range datacenterMap {
				fmt.Printf("\nDatacenter: %s\n", datacenter)

				// Build table header
				tbl := table.New("Service ID", "Hostname", "Primary IP", "Hourly Billing", "Friendly Name")
				tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

				for _, vps := range vpsServers {
					hourlyBilling := "No"
					if vps.HourlyBilling == 1 {
						hourlyBilling = "Yes"
					}

					tbl.AddRow(vps.ServiceID,
						vps.Hostname,
						vps.IPAddress,
						hourlyBilling,
						vps.Label)
				}

				tbl.Print()
			}
		} else {
			// Build table header
			tbl := table.New("Service ID", "Hostname", "Primary IP", "Datacenter", "Hourly Billing", "Friendly Name")
			tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

			// Build each VPS in a row
			for _, vps := range vpsList {
				// Use "contains" logic instead of exact matching
				if filterBySite != "" && !strings.Contains(strings.ToLower(vps.Location.Name), strings.ToLower(filterBySite)) {
					continue
				}

				hourlyBilling := "No"
				if vps.HourlyBilling == 1 {
					hourlyBilling = "Yes"
				}

				tbl.AddRow(vps.ServiceID,
					vps.Hostname,
					vps.IPAddress,
					vps.Location.Name,
					hourlyBilling,
					vps.Label)
			}

			tbl.Print()
		}
	},
}

func init() {
	vpsCmd.AddCommand(vpsListCmd)

	vpsListCmd.Flags().BoolP("display", "d", false, "Display VPS servers grouped per site")
	vpsListCmd.Flags().StringP("site", "s", "", "Filter VPS servers by site")
}
