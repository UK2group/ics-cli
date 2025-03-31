package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// bmdDeployCmd represents the order command
var bmdDeployCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new Baremetal Server Order",
	Long: `Create a new Baremetal Server order with the specified configuration.
    
Required parameters include:
--sku Server type (SKU)
--datacenter Datacenter location
--os Operating system

Optional parameters include:
--ssh-keys SSH keys
--license Software licenses
--bandwidth Additional bandwidth
--support Support level
--quantity Quantity (defaults to 1)`,
	Example: `  # Order a server with Debian 11
  ics-cli baremetal create --sku c1.small --datacenter NYC1 --os DEBIAN_11
  
  # Order a server with SSH keys and support
  ics-cli baremetal create --sku c1.small --datacenter NYC1 --os DEBIAN_11 --ssh-keys "My Key,Work Key" --support BASICSUP`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get required parameters
		sku, _ := cmd.Flags().GetString("sku")
		datacenter, _ := cmd.Flags().GetString("datacenter")
		osCode, _ := cmd.Flags().GetString("os")

		// Validate required parameters
		if sku == "" {
			fmt.Fprintln(os.Stderr, "Error: --sku flag is required")
			return
		}

		if datacenter == "" {
			fmt.Fprintln(os.Stderr, "Error: --datacenter flag is required")
			return
		}

		if osCode == "" {
			fmt.Fprintln(os.Stderr, "Error: --os flag is required")
			return
		}

		// Get optional parameters
		quantity, _ := cmd.Flags().GetInt("quantity")
		licenseCode, _ := cmd.Flags().GetString("license")
		bandwidthTB, _ := cmd.Flags().GetInt("bandwidth")
		supportCode, _ := cmd.Flags().GetString("support")
		sshKeysStr, _ := cmd.Flags().GetString("ssh-keys")

		// Default quantity to 1 if not specified
		if quantity <= 0 {
			quantity = 1
		}

		// Build the request
		orderRequest := OrderRequest{
			SKUProductName:             sku,
			Quantity:                   quantity,
			LocationCode:               datacenter,
			OperatingSystemProductCode: osCode,
		}

		// Add optional parameters if provided
		if licenseCode != "" {
			orderRequest.LicenseProductCode = licenseCode
		}

		if bandwidthTB > 0 {
			orderRequest.AdditionalBandwidthTB = bandwidthTB
		}

		if supportCode != "" {
			orderRequest.SupportLevelProductCode = supportCode
		}

		// Process SSH keys if provided
		if sshKeysStr != "" {
			keyNames := strings.Split(sshKeysStr, ",")
			keyIDs := []int{}

			for _, keyName := range keyNames {
				keyName = strings.TrimSpace(keyName)
				if keyName == "" {
					continue
				}

				// Look up the key ID from the label
				key, err := getSSHKeyFromLabel(keyName)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error finding SSH key '%s': %v\n", keyName, err)
					return
				}

				keyIDs = append(keyIDs, key.ID)
			}

			if len(keyIDs) > 0 {
				orderRequest.SSHKeyIDs = keyIDs
			}
		}

		// Confirm the order details with the user
		confirmOrder := confirmOrderDetails(orderRequest)
		if !confirmOrder {
			fmt.Printf("%s", RedText("Order cancelled, you have not been charged."))
			return
		}

		// Place the order
		orderResult, err := placeOrder(orderRequest)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error placing order: %v\n", err)
			return
		}

		// Display the order result
		fmt.Println(GreenText("\nOrder placed successfully"))
		fmt.Printf("Service IDs: %v", orderResult.Data.OrderServiceIDs)
		fmt.Println("\nServices in this order will be provisioned within 60 minutes.")
	},
}

func init() {
	baremetalDeployCmd.AddCommand(bmdDeployCmd)

	// Required flags
	bmdDeployCmd.Flags().String("sku", "", "Server type/SKU (required, e.g., c1i.small)")
	bmdDeployCmd.Flags().String("datacenter", "", "Datacenter location (required, e.g., NYC1)")
	bmdDeployCmd.Flags().String("os", "", "Operating system product code (required, e.g., DEBIAN_11)")

	// Optional flags
	bmdDeployCmd.Flags().Int("quantity", 1, "Number of servers to order (default: 1)")
	bmdDeployCmd.Flags().String("license", "", "License product code (e.g., CPANEL100)")
	bmdDeployCmd.Flags().Int("bandwidth", 0, "Additional bandwidth in TB")
	bmdDeployCmd.Flags().String("support", "", "Support level product code (e.g., BASICSUP)")
	bmdDeployCmd.Flags().String("ssh-keys", "", "Comma-separated list of SSH key names to assign")

	// Mark required flags
	bmdDeployCmd.MarkFlagRequired("sku")
	bmdDeployCmd.MarkFlagRequired("datacenter")
	bmdDeployCmd.MarkFlagRequired("os")
}
