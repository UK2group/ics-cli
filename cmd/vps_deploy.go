package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
)

// vpsDeployCmd represents the deploy command group
var vpsDeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy and manage VPS server deployments",
	Long:  `Commands for deploying new VPS servers, viewing available plans, and checking deployment options.`,
}

// vpsDeployListPlansCmd lists available VPS plans
var vpsDeployListPlansCmd = &cobra.Command{
	Use:   "list-plans",
	Short: "List available VPS plans",
	Long:  `List all available VPS plans with their specifications and pricing.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get VPS plans from API
		plans, err := getVPSPlans()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting VPS plans: %v\n", err)
			return
		}

		if len(plans) == 0 {
			fmt.Println("No VPS plans available.")
			return
		}

		headerFmt := color.New(color.FgBlue, color.Bold).SprintfFunc()
		columnFmt := color.New(color.FgYellow).SprintfFunc()

		// Create table with headers
		tbl := table.New("Plan", "Description", "Monthly Price")
		tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

		// Add each plan to the table
		for _, plan := range plans {
			price := plan.Price
			if !strings.HasPrefix(price, "$") {
				price = "$" + price
			}

			tbl.AddRow(
				plan.ProductName,
				plan.Label,
				price,
			)
		}

		tbl.Print()
	},
}

// vpsDeployListOSCmd lists available operating systems
var vpsDeployListOSCmd = &cobra.Command{
	Use:   "list-os [location]",
	Short: "List available operating systems for VPS",
	Long:  `List all available operating system templates for VPS deployment in a specific location.`,
	Example: `  # List operating systems for a location
  ics-cli vps deploy list-os "Salt Lake City - K, US"`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Get location name from argument
		locationName := args[0]

		// Get location ID by name
		locationID, err := getVPSLocationIDByName(locationName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			fmt.Fprintln(os.Stderr, "Use 'ics-cli vps deploy list-locations' to see available locations")
			return
		}

		// Get operating systems from API
		osList, err := getVPSOSList(locationID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting operating systems: %v\n", err)
			return
		}

		if len(osList) == 0 {
			fmt.Println("No operating systems available for this location.")
			return
		}

		headerFmt := color.New(color.FgBlue, color.Bold).SprintfFunc()
		columnFmt := color.New(color.FgYellow).SprintfFunc()

		// Create table with headers
		tbl := table.New("Name", "Architecture", "Price")
		tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

		// Add each OS to the table
		for _, os := range osList {
			// Format price
			priceStr := "Free"
			if os.Price > 0 {
				priceStr = fmt.Sprintf("$%.2f /mo", os.Price)
			}

			tbl.AddRow(
				os.Name,
				os.OperatingSystemArch,
				priceStr,
			)
		}

		tbl.Print()
	},
}

// vpsDeployCreateCmd creates a new VPS
var vpsDeployCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Deploy a new VPS server",
	Long:  `Deploy a new VPS server with the specified plan, location, and operating system.`,
	Example: `  # Deploy a VPS with basic configuration
  ics-cli vps deploy create --plan V1 --location "Chicago - G, US" --os "Debian 9.0" --hostname myserver

  # Deploy a VPS with hourly billing and backups
  ics-cli vps deploy create --plan V2 --location "Chicago - G, US" --os "Ubuntu 22.04" --hostname web1 --bill-hourly --backups --friendlyname "Web Server 1"`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get flags
		plan, _ := cmd.Flags().GetString("plan")
		location, _ := cmd.Flags().GetString("location")
		osName, _ := cmd.Flags().GetString("os")
		hostname, _ := cmd.Flags().GetString("hostname")
		friendlyName, _ := cmd.Flags().GetString("friendlyname")
		billHourly, _ := cmd.Flags().GetBool("bill-hourly")
		backups, _ := cmd.Flags().GetBool("backups")

		// Validate required fields
		if plan == "" || location == "" || osName == "" || hostname == "" {
			fmt.Fprintln(os.Stderr, "Error: --plan, --location, --os, and --hostname are required")
			return
		}

		// Get location ID from location name
		locationID, err := getVPSLocationIDByName(location)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			fmt.Fprintln(os.Stderr, "Use 'ics-cli vps deploy list-locations' to see available locations")
			return
		}

		// Get OS component code from OS name
		osComponentCode, err := getVPSOSComponentCodeByName(osName, locationID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			fmt.Fprintf(os.Stderr, "Use 'ics-cli vps deploy list-os \"%s\"' to see available operating systems\n", location)
			return
		}

		// Use hostname as friendly name if not provided
		if friendlyName == "" {
			friendlyName = hostname
		}

		// Display order details
		fmt.Println(BlueHeading("=== VPS Deployment Details ==="))
		fmt.Printf("%s %s\n", BlueHeading("Plan:"), WhiteText(plan))
		fmt.Printf("%s %s\n", BlueHeading("Location:"), WhiteText(location))
		fmt.Printf("%s %s\n", BlueHeading("Operating System:"), WhiteText(osName))
		fmt.Printf("%s %s\n", BlueHeading("Hostname:"), WhiteText(hostname))
		fmt.Printf("%s %s\n", BlueHeading("Friendly Name:"), WhiteText(friendlyName))

		billingType := "Monthly"
		if billHourly {
			billingType = "Hourly"
		}
		fmt.Printf("%s %s\n", BlueHeading("Billing Type:"), WhiteText(billingType))

		if backups {
			fmt.Printf("%s %s\n", BlueHeading("Backups:"), WhiteText("Enabled"))
		}

		fmt.Println() // Add blank line before deployment

		// Create deployment request
		request := VPSDeployRequest{
			ProductName:     plan,
			Label:           friendlyName,
			Hostname:        hostname,
			Backups:         backups,
			BillHourly:      billHourly,
			OSComponentCode: osComponentCode,
		}

		// Deploy VPS
		response, err := deployVPS(request, locationID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error deploying VPS: %v\n", err)
			return
		}

		// Display success message
		if response.StatusCode == 201 {
			fmt.Printf("\n%s\n", GreenText("✓ VPS deployment successful!"))
			fmt.Printf("%s\n", WhiteText(response.Message))
			fmt.Println("\nUse 'ics-cli vps list' to view your VPS servers")
		} else {
			fmt.Printf("\n%s\n", YellowText("VPS deployment may have issues"))
			fmt.Printf("%s\n", WhiteText(response.Message))
		}
	},
}

// vpsDeployListLocationsCmd lists available VPS locations
var vpsDeployListLocationsCmd = &cobra.Command{
	Use:   "list-locations",
	Short: "List available VPS datacenter locations",
	Long:  `List all available datacenter locations where VPS servers can be deployed.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get VPS locations from API
		locations, err := getVPSLocations()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting VPS locations: %v\n", err)
			return
		}

		if len(locations) == 0 {
			fmt.Println("No VPS locations available.")
			return
		}

		headerFmt := color.New(color.FgBlue, color.Bold).SprintfFunc()
		columnFmt := color.New(color.FgYellow).SprintfFunc()

		// Create table with headers
		tbl := table.New("Location")
		tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

		// Add each location to the table
		for _, location := range locations {
			tbl.AddRow(location.Name)
		}

		tbl.Print()
	},
}

func init() {
	vpsCmd.AddCommand(vpsDeployCmd)
	vpsDeployCmd.AddCommand(vpsDeployListPlansCmd)
	vpsDeployCmd.AddCommand(vpsDeployListOSCmd)
	vpsDeployCmd.AddCommand(vpsDeployListLocationsCmd)
	vpsDeployCmd.AddCommand(vpsDeployCreateCmd)

	// Flags for create
	vpsDeployCreateCmd.Flags().String("plan", "", "VPS plan (required, e.g., V1, V2)")
	vpsDeployCreateCmd.Flags().String("location", "", "Datacenter location (required, e.g., \"Chicago - G, US\")")
	vpsDeployCreateCmd.Flags().String("os", "", "Operating system (required, e.g., \"Debian 9.0\")")
	vpsDeployCreateCmd.Flags().String("hostname", "", "Hostname for the VPS - letters, numbers, hyphens only (required, e.g., web1)")
	vpsDeployCreateCmd.Flags().String("friendlyname", "", "Friendly name/label for the VPS (defaults to hostname)")
	vpsDeployCreateCmd.Flags().Bool("bill-hourly", false, "Enable hourly billing")
	vpsDeployCreateCmd.Flags().Bool("backups", false, "Enable backups")

	// Mark required flags
	vpsDeployCreateCmd.MarkFlagRequired("plan")
	vpsDeployCreateCmd.MarkFlagRequired("location")
	vpsDeployCreateCmd.MarkFlagRequired("os")
	vpsDeployCreateCmd.MarkFlagRequired("hostname")
}
