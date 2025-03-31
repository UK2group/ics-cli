package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

// getServerIDFromServiceID finds the server ID corresponding to a service ID
func getServerIDFromServiceID(serviceID string) (string, error) {
	// Use makeAPIRequest to get the server list
	var serverList ServerResponse

	err := makeAPIRequest(
		"GET",
		30,
		"https://api.ingenuitycloudservices.com/rest-api/servers",
		nil,
		&serverList,
	)

	if err != nil {
		return "", fmt.Errorf("failed to get server list: %w", err)
	}

	// Find the server with the matching service ID
	for _, server := range serverList.Data {
		if strconv.Itoa(server.ServiceID) == serviceID {
			return server.ID, nil
		}
	}

	return "", fmt.Errorf("server with Service ID %s not found", serviceID)
}

// getSolResponse gets the SOL Link to a server
func getIkvmResponse(serverID string) (string, error) {
	var response RemoteAccessResponse

	err := makeAPIRequest(
		"POST",
		60,
		fmt.Sprintf("https://api.ingenuitycloudservices.com/rest-api/servers/%s/remote-access/ikvm", serverID),
		nil,
		&response,
	)

	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	return response.Data.Redirect, nil
}

// getSolResponse gets the SOL Link to a server
func setPXEUrl(serverID, url string) (bool, error) {
	var response RemoteAccessResponse

	// Pass url as pxe_script_url in the request body
	requestData := map[string]string{
		"pxe_script_url": url,
	}

	requestBody, err := json.Marshal(requestData)
	if err != nil {
		return false, fmt.Errorf("error creating request body: %w", err)
	}

	err = makeAPIRequest(
		"PUT",
		60,
		fmt.Sprintf("https://api.ingenuitycloudservices.com/rest-api/servers/%s/set-pxe", serverID),
		bytes.NewBuffer(requestBody),
		&response,
	)

	if err != nil {
		return false, fmt.Errorf("%w", err)
	}

	return true, nil
}

// setFriendlyName sets the frienly name of a server
func setFriendlyName(serverID, name string) (bool, error) {
	var response RemoteAccessResponse

	// Pass friendly name
	requestData := map[string]string{
		"friendly_name": name,
	}

	requestBody, err := json.Marshal(requestData)
	if err != nil {
		return false, fmt.Errorf("error creating request body: %w", err)
	}

	err = makeAPIRequest(
		"PUT",
		60,
		fmt.Sprintf("https://api.ingenuitycloudservices.com/rest-api/servers/%s/friendly-name", serverID),
		bytes.NewBuffer(requestBody),
		&response,
	)

	if err != nil {
		return false, fmt.Errorf("%w", err)
	}

	return true, nil
}

// getServerDetails gets detailed information about a server
func getServerDetails(serverID string) (*ServerDetail, error) {
	var response ServerDetailResponse

	err := makeAPIRequest(
		"GET",
		30,
		fmt.Sprintf("https://api.ingenuitycloudservices.com/rest-api/servers/%s", serverID),
		nil,
		&response,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get server details: %w", err)
	}

	return &response.Data, nil
}

// getPowerStatus gets the power status of a server
func getPowerStatus(serverID int) (bool, error) {
	var response PowerStatusResponse

	err := makeAPIRequest(
		"GET",
		30,
		fmt.Sprintf("https://api.ingenuitycloudservices.com/rest-api/servers/%d/power/status", serverID),
		nil,
		&response,
	)

	if err != nil {
		return false, fmt.Errorf("failed to get power status: %w", err)
	}

	return response.Data.IsPoweredOn, nil
}

// getServerSSHKeys gets the SSH keys assigned to a server
func getServerSSHKeys(serverID int) ([]AssignedSSHKey, error) {
	var keysResponse AssignedSSHKeysResponse

	err := makeAPIRequest(
		"GET",
		30,
		fmt.Sprintf("https://api.ingenuitycloudservices.com/rest-api/servers/%d/ssh-keys", serverID),
		nil,
		&keysResponse,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get SSH keys: %w", err)
	}

	return keysResponse.Data, nil
}

// Now update the printServerDetails function to include SSH keys
func printServerDetails(cmd *cobra.Command, server *ServerDetail) {

	// General Information Section
	//fmt.Println(blue("=== SERVER INFORMATION ==="))
	fmt.Printf("%s %s\n", BlueHeading("Server ID:"), WhiteText(strconv.Itoa(server.ServiceID)))
	fmt.Printf("%s %s\n", BlueHeading("Hostname:"), WhiteText(server.Hostname))
	fmt.Printf("%s %s\n", BlueHeading("Friendly Name:"), WhiteText(server.FriendlyName))
	fmt.Printf("%s %s\n", BlueHeading("Datacenter:"), WhiteText(server.DatacenterName))
	fmt.Printf("%s %s\n", BlueHeading("Operating System:"), WhiteText(server.OperatingSystemName))

	// Network Information Section
	fmt.Println(BlueHeading("\n=== NETWORK INFORMATION ==="))
	fmt.Printf("%s %s\n", BlueHeading("Primary IP Address:"), WhiteText(server.PublicIP))
	fmt.Printf("%s %s\n", BlueHeading("Primary MAC Address:"), WhiteText(server.MacAddress))
	for _, ip := range server.IPAddresses {
		if ip.IsPrimary {
			continue
		}
		fmt.Printf("%s %s\n", BlueHeading("Secondary IP Address:"), WhiteText(ip.IPAddress))
	}

	// Operating System Section
	fmt.Println(BlueHeading("\n=== CREDENTIALS ==="))
	fmt.Printf("%s %s\n", BlueHeading("Username:"), WhiteText(server.OperatingSystemUsername))

	// Redact the root password by default
	showRootPassword, _ := cmd.Flags().GetBool("password")
	if !showRootPassword {
		server.OperatingSystemPassword = "********"
	}
	fmt.Printf("%s %s\n", BlueHeading("Password:"), WhiteText(server.OperatingSystemPassword))

	// Add SSH Keys section
	sshKeys, err := getServerSSHKeys(server.ServerID)
	if err == nil && len(sshKeys) > 0 {
		// Create a slice to hold the key labels
		keyLabels := make([]string, len(sshKeys))
		for i, key := range sshKeys {
			keyLabels[i] = key.Label
		}

		// Join the key labels with commas
		fmt.Printf("%s %s\n", BlueHeading("Assigned SSH Keys:"), WhiteText(strings.Join(keyLabels, ", ")))
	} else if err != nil {
		// Optionally handle error silently or show a message
		fmt.Printf("%s %s\n", BlueHeading("Assigned SSH Keys:"), RedText("Unable to retrieve SSH keys"))
	} else {
		fmt.Printf("%s %s\n", BlueHeading("Assigned SSH Keys:"), WhiteText("None"))
	}

	hidePowerStatus, _ := cmd.Flags().GetBool("power")
	if !hidePowerStatus {
		// Power Status Section
		fmt.Println(BlueHeading("\n=== POWER STATUS ==="))

		isPoweredOn, err := getPowerStatus(server.ServerID)
		if err != nil {
			fmt.Printf("%s %s\n", BlueHeading("Power State:"), RedText("UNKNOWN"))
		} else {
			if isPoweredOn {
				fmt.Printf("%s %s\n", BlueHeading("Power State:"), GreenText("Powered On"))
			} else {
				fmt.Printf("%s %s\n", BlueHeading("Power State:"), RedText("Powered Off"))
			}
		}
	}
}

// Get server list
func getServerList() ([]Server, error) {
	var serverResp ServerResponse

	err := makeAPIRequest(
		"GET",
		30,
		"https://api.ingenuitycloudservices.com/rest-api/servers",
		nil, // no request body for GET
		&serverResp,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get server list: %w", err)
	}

	return serverResp.Data, nil
}

// setPowerOff sends a power off
func setPowerOff(serverID string) (bool, error) {
	var response GenericServerResponse

	err := makeAPIRequest(
		"POST",
		60,
		fmt.Sprintf("https://api.ingenuitycloudservices.com/rest-api/servers/%s/power/off", serverID),
		nil,
		&response,
	)

	if err != nil {
		return false, fmt.Errorf("%w", err)
	}

	// Check for success
	if !response.Data.Success {
		return false, nil
	}

	return true, nil
}

// setPowerOn sends a power on
func setPowerOn(serverID string) (bool, error) {
	var response GenericServerResponse

	err := makeAPIRequest(
		"POST",
		60,
		fmt.Sprintf("https://api.ingenuitycloudservices.com/rest-api/servers/%s/power/on", serverID),
		nil,
		&response,
	)

	if err != nil {
		return false, fmt.Errorf("%w", err)
	}

	// Check for success
	if !response.Data.Success {
		return false, nil
	}

	return true, nil
}

// setReboot sends a power off
func setReboot(serverID string) (bool, error) {
	var response GenericServerResponse

	err := makeAPIRequest(
		"POST",
		60,
		fmt.Sprintf("https://api.ingenuitycloudservices.com/rest-api/servers/%s/power/reboot", serverID),
		nil,
		&response,
	)

	if err != nil {
		return false, fmt.Errorf("%w", err)
	}

	// Check for success
	if !response.Data.Success {
		return false, nil
	}

	return true, nil
}

// setReboot sends a power off
func setRecovery(serverID string) (bool, error) {
	var response GenericServerResponse

	err := makeAPIRequest(
		"POST",
		60,
		fmt.Sprintf("https://api.ingenuitycloudservices.com/rest-api/servers/%s/power/off", serverID),
		nil,
		&response,
	)

	if err != nil {
		return false, fmt.Errorf("%w", err)
	}

	// Check for success
	if !response.Data.Success {
		return false, nil
	}

	return true, nil
}

// getSolResponse gets the SOL Link to a server
func getSolResponse(serverID string) (string, error) {
	var response RemoteAccessResponse

	err := makeAPIRequest(
		"POST",
		60,
		fmt.Sprintf("https://api.ingenuitycloudservices.com/rest-api/servers/%s/remote-access/sol", serverID),
		nil,
		&response,
	)

	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	return response.Data.Redirect, nil
}

// getOSList retrieves the list of available operating systems for a server
func getOSList(serverID string) ([]OS, error) {
	var response OSListResponse

	err := makeAPIRequest(
		"GET",
		60,
		fmt.Sprintf("https://api.ingenuitycloudservices.com/rest-api/servers/%s/provision/os-list", serverID),
		nil,
		&response,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get OS list: %w", err)
	}

	return response.Data.OSList, nil
}

// setPowerOff sends a power off
func setReinstallOS(serverID, reason, imageId string) (bool, error) {
	var response GenericServerResponse

	requestData := map[string]string{
		"os_image_id": imageId,
		"reason":      reason,
	}

	requestBody, err := json.Marshal(requestData)
	if err != nil {
		return false, fmt.Errorf("error creating request body: %w", err)
	}

	err = makeAPIRequest(
		"POST",
		60,
		fmt.Sprintf("https://api.ingenuitycloudservices.com/rest-api/servers/%s/provision/reload-os", serverID),
		bytes.NewBuffer(requestBody),
		&response,
	)

	if err != nil {
		return false, fmt.Errorf("%w", err)
	}

	// Check for success
	if !response.Data.Success {
		return false, nil
	}

	return true, nil
}

// getInventory retrieves server inventory from the API
func getInventory() ([]InventoryDetails, error) {
	var response InventoryResponse

	err := makeAPIRequest(
		"GET",
		60,
		"https://api.ingenuitycloudservices.com/rest-api/server-orders/inventory",
		nil,
		&response,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get inventory: %w", err)
	}

	return response.Data, nil
}

// filterInventory applies filters to the inventory data
func filterInventory(inventory []InventoryDetails, datacenter, sku string, minPrice, maxPrice float64) []InventoryDetails {
	filtered := make([]InventoryDetails, 0)

	for _, item := range inventory {
		// Filter by datacenter
		if datacenter != "" && !strings.EqualFold(item.LocationCode, datacenter) {
			continue
		}

		// Filter by SKU
		if sku != "" && !strings.EqualFold(item.SkuProductName, sku) {
			continue
		}

		// Filter by price
		itemPrice, err := strconv.ParseFloat(item.Price, 64)
		if err == nil {
			if minPrice > 0 && itemPrice < minPrice {
				continue
			}
			if maxPrice > 0 && itemPrice > maxPrice {
				continue
			}
		}

		// This item passed all filters, include it
		filtered = append(filtered, item)
	}

	return filtered
}

// groupInventory groups inventory by location and SKU
func groupInventory(inventory []InventoryDetails) []GroupedInventory {
	// Create a map to store grouped inventory
	groups := make(map[string]GroupedInventory)

	// Group the inventory
	for _, item := range inventory {
		// Create a key by combining location and SKU
		key := item.LocationCode + ":" + item.SkuProductName

		if group, exists := groups[key]; exists {
			// If the group exists, add to the quantity
			group.TotalQuantity += item.Quantity
			groups[key] = group
		} else {
			// Create a new group
			groups[key] = GroupedInventory{
				LocationCode:   item.LocationCode,
				SkuProductName: item.SkuProductName,
				Price:          item.Price,
				TotalQuantity:  item.Quantity,
			}
		}
	}

	// Convert the map to a slice for sorting
	result := make([]GroupedInventory, 0, len(groups))
	for _, group := range groups {
		result = append(result, group)
	}

	// Sort by location, then by SKU
	sort.Slice(result, func(i, j int) bool {
		if result[i].LocationCode != result[j].LocationCode {
			return result[i].LocationCode < result[j].LocationCode
		}
		return result[i].SkuProductName < result[j].SkuProductName
	})

	return result
}

// getAddons retrieves available add-ons for a server type in a datacenter
func getAddons(sku, datacenter string) (AddonsResponse, error) {
	var response AddonsResponse

	err := makeAPIRequest(
		"GET",
		60,
		fmt.Sprintf("https://api.ingenuitycloudservices.com/rest-api/server-orders/list-addons?sku_product_name=%s&location_code=%s", sku, datacenter),
		nil,
		&response,
	)

	if err != nil {
		return AddonsResponse{}, fmt.Errorf("failed to get add-ons: %w", err)
	}

	return response, nil
}

// placeOrder sends the order request to the API
func placeOrder(order OrderRequest) (OrderResponse, error) {
	var response OrderResponse

	// Create request body
	requestBody, err := json.Marshal(order)
	if err != nil {
		return OrderResponse{}, fmt.Errorf("error creating request body: %w", err)
	}

	// Make API request
	err = makeAPIRequest(
		"POST",
		60,
		"https://api.ingenuitycloudservices.com/rest-api/server-orders/order",
		bytes.NewBuffer(requestBody),
		&response,
	)

	if err != nil {
		return OrderResponse{}, fmt.Errorf("failed to place order: %w", err)
	}

	return response, nil
}

// confirmOrderDetails displays the order details and asks for confirmation
func confirmOrderDetails(order OrderRequest) bool {
	fmt.Println(BlueHeading("=== Order Details ==="))
	fmt.Printf("%s %s\n", BlueHeading("Server Type:"), WhiteText(order.SKUProductName))
	fmt.Printf("%s %s\n", BlueHeading("Datacenter:"), WhiteText(order.LocationCode))
	fmt.Printf("%s %s\n", BlueHeading("Operating System:"), WhiteText(order.OperatingSystemProductCode))
	fmt.Printf("%s %d\n", BlueHeading("Quantity:"), (order.Quantity))

	if order.LicenseProductCode != "" {
		fmt.Printf("%s %s\n", BlueHeading("License:"), WhiteText(order.LicenseProductCode))
	}

	if order.AdditionalBandwidthTB > 0 {
		fmt.Printf("%s %d TB\n", BlueHeading("Additional Bandwidth:"), (order.AdditionalBandwidthTB))
	}

	if order.SupportLevelProductCode != "" {
		fmt.Printf("%s %s\n", BlueHeading("Support Level:"), WhiteText(order.SupportLevelProductCode))
	}

	if len(order.SSHKeyIDs) > 0 {
		fmt.Printf("%s %d\n", BlueHeading("SSH Keys:"), (order.SSHKeyIDs))
	}

	fmt.Printf("%s", BlueHeading("\nAre you sure you want to place this order? (y/N):"))
	var response string
	fmt.Scanln(&response)

	return strings.ToLower(response) == "y" || strings.ToLower(response) == "yes"
}
