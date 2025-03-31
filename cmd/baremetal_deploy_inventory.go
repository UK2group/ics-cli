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
Shows the total quantity available, location, server type, and price.

You can filter the results using various flags:
  --datacenter NYC1     Show only inventory in a specific datacenter
  --sku c1.small       Show only a specific server type
  --min-price 100       Show servers with price >= $100
  --max-price 300       Show servers with price <= $300`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get filter flags
		datacenter, _ := cmd.Flags().GetString("datacenter")
		sku, _ := cmd.Flags().GetString("sku")
		minPriceStr, _ := cmd.Flags().GetString("min-price")
		maxPriceStr, _ := cmd.Flags().GetString("max-price")

		// Parse price filters
		var minPrice, maxPrice float64
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

		// Get the inventory from API
		inventory, err := getInventory()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting inventory: %v\n", err)
			return
		}

		// Apply filters
		filteredInventory := filterInventory(inventory, datacenter, sku, minPrice, maxPrice)

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
		tbl := table.New("Location", "Server Type", "Price (USD)", "Available Quantity")
		tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

		// Add each item to the table
		for _, item := range groupedInventory {

			// Format the price to include currency symbol if not present
			price := item.Price
			if !strings.HasPrefix(price, "$") {
				price = "$" + price
			}

			tbl.AddRow(
				item.LocationCode,
				item.SkuProductName,
				price,
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
	bmdInventoryCmd.Flags().String("min-price", "", "Filter by minimum price")
	bmdInventoryCmd.Flags().String("max-price", "", "Filter by maximum price")
}
