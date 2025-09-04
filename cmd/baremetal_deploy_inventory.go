package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
)

// inventoryCmd represents the inventory command
var bmdInventoryCmd = &cobra.Command{
	Use:   "list-inventory",
	Short: "List available Baremetal Server inventory",
	Long: `Display available Baremetal Server inventory grouped by location and server type.
Shows the total quantity available, location, server type, monthly pricing, hourly billing availability, and hourly pricing.

You can filter the results using various flags:
  --datacenter NYC1         Show only inventory in a specific datacenter
  --sku c1.small           Show only a specific server type
  --min-price 100           Show servers with monthly price >= $100
  --max-price 300           Show servers with monthly price <= $300
  --min-hourly-price 1.5    Show servers with hourly price >= $1.5
  --max-hourly-price 5.0    Show servers with hourly price <= $5.0
  --hourly-enabled          Show only servers with hourly billing enabled`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get filter flags
		datacenter, _ := cmd.Flags().GetString("datacenter")
		sku, _ := cmd.Flags().GetString("sku")
		minPriceStr, _ := cmd.Flags().GetString("min-price")
		maxPriceStr, _ := cmd.Flags().GetString("max-price")
		minHourlyPriceStr, _ := cmd.Flags().GetString("min-hourly-price")
		maxHourlyPriceStr, _ := cmd.Flags().GetString("max-hourly-price")
		hourlyEnabled, _ := cmd.Flags().GetBool("hourly-enabled")

		// Parse price filters
		var minPrice, maxPrice, minHourlyPrice, maxHourlyPrice float64
		var err error

		if minPriceStr != "" {
			minPrice, err = strconv.ParseFloat(minPriceStr, 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: Invalid min-price value: %v\n", err)
				return
			}
		}

		if maxPriceStr != "" {
			maxPrice, err = strconv.ParseFloat(maxPriceStr, 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: Invalid max-price value: %v\n", err)
				return
			}
		}

		if minHourlyPriceStr != "" {
			minHourlyPrice, err = strconv.ParseFloat(minHourlyPriceStr, 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: Invalid min-hourly-price value: %v\n", err)
				return
			}
		}

		if maxHourlyPriceStr != "" {
			maxHourlyPrice, err = strconv.ParseFloat(maxHourlyPriceStr, 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: Invalid max-hourly-price value: %v\n", err)
				return
			}
		}

		// Get the inventory from API
		inventory, err := getInventory()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting inventory: %v\n", err)
			return
		}

		// Apply filters
		filteredInventory := filterInventory(inventory, datacenter, sku, minPrice, maxPrice, minHourlyPrice, maxHourlyPrice, hourlyEnabled)

		// Group the filtered inventory by location and SKU
		groupedInventory := groupInventory(filteredInventory)

		// Display the grouped inventory
		if len(groupedInventory) == 0 {
			fmt.Println("No inventory available.")
			return
		}

		headerFmt := color.New(color.FgBlue, color.Bold).SprintfFunc()
		columnFmt := color.New(color.FgYellow).SprintfFunc()

		// Create table with headers
		tbl := table.New("Location", "Server Type", "Monthly Price", "Hourly Enabled", "Hourly Price", "Available Quantity")
		tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

		// Add each item to the table
		for _, item := range groupedInventory {

			// Format the monthly price to include currency symbol if not present
			monthlyPrice := item.Price
			if !strings.HasPrefix(monthlyPrice, "$") {
				monthlyPrice = "$" + monthlyPrice
			}

			// Format hourly enabled status
			var hourlyEnabledText string
			if item.HourlyEnabled {
				hourlyEnabledText = "Yes"
			} else {
				hourlyEnabledText = "No"
			}

			// Format hourly price if available
			var hourlyPrice string
			if item.HourlyEnabled {
				hourlyPrice = item.HourlyPrice
				if !strings.HasPrefix(hourlyPrice, "$") {
					hourlyPrice = "$" + hourlyPrice
				}
			} else {
				hourlyPrice = "N/A"
			}

			tbl.AddRow(
				item.LocationCode,
				item.SkuProductName,
				monthlyPrice,
				hourlyEnabledText,
				hourlyPrice,
				strconv.Itoa(item.TotalQuantity),
			)
		}

		tbl.Print()
	},
}

func init() {
	baremetalDeployCmd.AddCommand(bmdInventoryCmd)

	// Add filter flags
	bmdInventoryCmd.Flags().String("datacenter", "", "Filter by datacenter (e.g., NYC1)")
	bmdInventoryCmd.Flags().String("sku", "", "Filter by server type (e.g., c1.small)")
	bmdInventoryCmd.Flags().String("min-price", "", "Filter by minimum monthly price")
	bmdInventoryCmd.Flags().String("max-price", "", "Filter by maximum monthly price")
	bmdInventoryCmd.Flags().String("min-hourly-price", "", "Filter by minimum hourly price")
	bmdInventoryCmd.Flags().String("max-hourly-price", "", "Filter by maximum hourly price")
	bmdInventoryCmd.Flags().Bool("hourly-enabled", false, "Show only servers with hourly billing enabled")
}
