package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
)

// bmdAddonsCmd represents the list-addons command
var bmdAddonsCmd = &cobra.Command{
	Use:   "list-addons",
	Short: "List available add-ons for a Baremetal Server",
	Long: `Display available add-ons for a Baremetal Server including:
- Operating Systems
- Software Licenses
- Support Levels

You must specify both the SKU (server type) and datacenter location.`,
	Example: `  # List all add-ons for a specific server type in a location
  ics-cli baremetal list-addons --sku c1.small --datacenter NYC1`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get required parameters
		sku, _ := cmd.Flags().GetString("sku")
		datacenter, _ := cmd.Flags().GetString("datacenter")

		// Validate required parameters
		if sku == "" {
			fmt.Fprintln(os.Stderr, "Error: --sku flag is required")
			return
		}

		if datacenter == "" {
			fmt.Fprintln(os.Stderr, "Error: --datacenter flag is required")
			return
		}

		// Get add-ons for the specified SKU and datacenter
		addons, err := getAddons(sku, datacenter)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting add-ons: %v\n", err)
			return
		}

		// Display the add-ons in tables
		headerFmt := color.New(color.FgBlue, color.Bold).SprintfFunc()
		columnFmt := color.New(color.FgYellow).SprintfFunc()
		headerText := color.New(color.FgGreen, color.Bold).SprintfFunc()

		// Display Operating Systems
		if len(addons.Data.OperatingSystems.Products) > 0 {
			fmt.Println(headerText("\n%s", addons.Data.OperatingSystems.Name))

			// Create table for Operating Systems
			osTbl := table.New("Name", "Type", "Product Code", "Price")
			osTbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

			// Add each OS to the table
			for _, os := range addons.Data.OperatingSystems.Products {
				// Format price
				priceStr := "Free"
				if os.Price > 0 {
					priceStr = "$" + fmt.Sprintf("%.2f", os.Price) + " /mo"
				}

				// Check if price is per core
				if os.PricePerCore != nil {
					if perCore, ok := os.PricePerCore.(float64); ok && perCore > 0 {
						priceStr = fmt.Sprintf("$%.2f per core", perCore) + " /mo"
					}
				}

				osTbl.AddRow(
					os.Name,
					os.OSType,
					os.ProductCode,
					priceStr,
				)
			}

			osTbl.Print()
		}

		// Display Licenses
		if len(addons.Data.Licenses.Products) > 0 {
			fmt.Println(headerText("\n%s", addons.Data.Licenses.Name))

			// Create table for Licenses
			licTbl := table.New("Name", "Product Code", "Price")
			licTbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

			// Add each license to the table
			for _, lic := range addons.Data.Licenses.Products {
				// Format price
				priceStr := "Free"
				if lic.Price > 0 {
					priceStr = "$" + fmt.Sprintf("%.2f", lic.Price) + " /mo"
				}

				licTbl.AddRow(
					lic.Name,
					lic.ProductCode,
					priceStr,
				)
			}

			licTbl.Print()
		}

		// Display Support Levels
		if len(addons.Data.SupportLevels.Products) > 0 {
			fmt.Println(headerText("\n%s", addons.Data.SupportLevels.Name))

			// Create table for Support Levels
			supTbl := table.New("Name", "Description", "Product Code", "Price")
			supTbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

			// Add each support level to the table
			for _, sup := range addons.Data.SupportLevels.Products {
				// Format price
				priceStr := "Free"
				if sup.Price > 0 {
					priceStr = "$" + fmt.Sprintf("%.2f", sup.Price) + " /mo"
				}

				supTbl.AddRow(
					sup.Name,
					sup.Description,
					sup.ProductCode,
					priceStr,
				)
			}

			supTbl.Print()
		}
	},
}

func init() {
	baremetalDeployCmd.AddCommand(bmdAddonsCmd)

	// Add required flags
	bmdAddonsCmd.Flags().String("sku", "", "Server type/SKU (required, e.g., c1.small)")
	bmdAddonsCmd.Flags().String("datacenter", "", "Datacenter location (required, e.g., NYC1)")

	// Mark flags as required
	bmdAddonsCmd.MarkFlagRequired("sku")
	bmdAddonsCmd.MarkFlagRequired("datacenter")
}
