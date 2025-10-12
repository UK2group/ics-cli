package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

// getVPSIDFromServiceID finds the VPS ID corresponding to a service ID
func getVPSIDFromServiceID(serviceID string) (int, error) {
	// Use makeAPIRequest to get the VPS list
	var vpsList VPSResponse

	err := makeAPIRequest(
		"GET",
		30,
		fmt.Sprintf("%s/rest-api/ssd-vps/servers", getAPIBaseURL()),
		nil,
		&vpsList,
	)

	if err != nil {
		return 0, fmt.Errorf("failed to get VPS list: %w", err)
	}

	// Find the VPS with the matching service ID
	for _, vps := range vpsList.Data {
		if strconv.Itoa(vps.ServiceID) == serviceID {
			return vps.ID, nil
		}
	}

	return 0, fmt.Errorf("VPS with Service ID %s not found", serviceID)
}

// getVPSIDAndLocationFromServiceID finds both the VPS ID and location ID from a service ID
func getVPSIDAndLocationFromServiceID(serviceID string) (int, int, error) {
	// Use makeAPIRequest to get the VPS list
	var vpsList VPSResponse

	err := makeAPIRequest(
		"GET",
		30,
		fmt.Sprintf("%s/rest-api/ssd-vps/servers", getAPIBaseURL()),
		nil,
		&vpsList,
	)

	if err != nil {
		return 0, 0, fmt.Errorf("failed to get VPS list: %w", err)
	}

	// Find the VPS with the matching service ID
	for _, vps := range vpsList.Data {
		if strconv.Itoa(vps.ServiceID) == serviceID {
			return vps.ID, vps.Location.ID, nil
		}
	}

	return 0, 0, fmt.Errorf("VPS with Service ID %s not found", serviceID)
}

// getVPSList retrieves the list of all VPS servers
func getVPSList() ([]VPSBasic, error) {
	var vpsResp VPSResponse

	err := makeAPIRequest(
		"GET",
		30,
		fmt.Sprintf("%s/rest-api/ssd-vps/servers", getAPIBaseURL()),
		nil,
		&vpsResp,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get VPS list: %w", err)
	}

	return vpsResp.Data, nil
}

// getVPSDetails gets detailed information about a VPS
func getVPSDetails(vpsID int, locationID int) (*VPSDetail, error) {
	var response VPSDetailResponse

	err := makeAPIRequest(
		"GET",
		30,
		fmt.Sprintf("%s/rest-api/ssd-vps/locations/%d/servers/%d", getAPIBaseURL(), locationID, vpsID),
		nil,
		&response,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get VPS details: %w", err)
	}

	return &response.Data, nil
}

// getVPSPowerStatus gets the power status of a VPS
func getVPSPowerStatus(vpsID string) (bool, error) {
	var response VPSPowerStatusResponse

	err := makeAPIRequest(
		"GET",
		30,
		fmt.Sprintf("%s/rest-api/vps/%s/power/status", getAPIBaseURL(), vpsID),
		nil,
		&response,
	)

	if err != nil {
		return false, fmt.Errorf("failed to get power status: %w", err)
	}

	return response.Data.IsPoweredOn, nil
}

// setVPSPowerOn powers on a VPS
func setVPSPowerOn(vpsID int, locationID int) (bool, error) {
	var response VPSEmptyDataResponse

	endpoint := fmt.Sprintf("%s/rest-api/ssd-vps/locations/%d/servers/%d/power/on", getAPIBaseURL(), locationID, vpsID)

	err := makeAPIRequest(
		"POST",
		60,
		endpoint,
		nil,
		&response,
	)

	if err != nil {
		return false, fmt.Errorf("%w", err)
	}

	// The API returns an empty data array on success, so we check the status code
	return response.StatusCode == 200, nil
}

// setVPSPowerOff powers off a VPS
func setVPSPowerOff(vpsID int, locationID int) (bool, error) {
	var response VPSEmptyDataResponse

	endpoint := fmt.Sprintf("%s/rest-api/ssd-vps/locations/%d/servers/%d/power/off", getAPIBaseURL(), locationID, vpsID)

	err := makeAPIRequest(
		"POST",
		60,
		endpoint,
		nil,
		&response,
	)

	if err != nil {
		return false, fmt.Errorf("%w", err)
	}

	// The API returns an empty data array on success, so we check the status code
	return response.StatusCode == 200, nil
}

// setVPSReboot reboots a VPS
func setVPSReboot(vpsID int, locationID int) (bool, error) {
	var response VPSEmptyDataResponse

	endpoint := fmt.Sprintf("%s/rest-api/ssd-vps/locations/%d/servers/%d/power/reboot", getAPIBaseURL(), locationID, vpsID)

	err := makeAPIRequest(
		"POST",
		60,
		endpoint,
		nil,
		&response,
	)

	if err != nil {
		return false, fmt.Errorf("%w", err)
	}

	// The API returns an empty data array on success, so we check the status code and message
	return response.StatusCode == 200, nil
}

// cancelVPS cancels a VPS server
func cancelVPS(vpsID int, locationID int) (VPSEmptyDataResponse, error) {
	var response VPSEmptyDataResponse

	err := makeAPIRequest(
		"DELETE",
		60,
		fmt.Sprintf("%s/rest-api/ssd-vps/locations/%d/servers/%d", getAPIBaseURL(), locationID, vpsID),
		nil,
		&response,
	)

	if err != nil {
		return VPSEmptyDataResponse{}, err
	}

	return response, nil
}

// getVPSPlans retrieves available VPS plans
func getVPSPlans() ([]VPSPlan, error) {
	var response VPSPlanResponse

	err := makeAPIRequest(
		"GET",
		60,
		fmt.Sprintf("%s/rest-api/ssd-vps/plans", getAPIBaseURL()),
		nil,
		&response,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get VPS plans: %w", err)
	}

	return response.Data, nil
}

// getVPSOSList retrieves available operating systems for VPS in a location
func getVPSOSList(locationID int) ([]VPSOS, error) {
	var response VPSOSResponse

	err := makeAPIRequest(
		"GET",
		60,
		fmt.Sprintf("%s/rest-api/ssd-vps/locations/%d/operating-systems", getAPIBaseURL(), locationID),
		nil,
		&response,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get VPS OS list: %w", err)
	}

	return response.Data, nil
}

// getVPSLocations retrieves available VPS locations
func getVPSLocations() ([]VPSLocation, error) {
	var response VPSLocationResponse

	err := makeAPIRequest(
		"GET",
		60,
		fmt.Sprintf("%s/rest-api/ssd-vps/locations", getAPIBaseURL()),
		nil,
		&response,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get VPS locations: %w", err)
	}

	return response.Data, nil
}

// getVPSLocationIDByName finds a location ID by its name
func getVPSLocationIDByName(locationName string) (int, error) {
	locations, err := getVPSLocations()
	if err != nil {
		return 0, err
	}

	for _, location := range locations {
		if location.Name == locationName {
			return location.ID, nil
		}
	}

	return 0, fmt.Errorf("location '%s' not found", locationName)
}

// getVPSOSComponentCodeByName finds an OS component code by its name in a specific location
func getVPSOSComponentCodeByName(osName string, locationID int) (string, error) {
	osList, err := getVPSOSList(locationID)
	if err != nil {
		return "", err
	}

	for _, os := range osList {
		if os.Name == osName {
			return os.ComponentCode, nil
		}
	}

	return "", fmt.Errorf("operating system '%s' not found in this location", osName)
}

// deployVPS deploys a new VPS server
func deployVPS(request VPSDeployRequest, locationID int) (VPSDeployResponse, error) {
	var response VPSDeployResponse

	// Create request body
	requestBody, err := json.Marshal(request)
	if err != nil {
		return VPSDeployResponse{}, fmt.Errorf("error creating request body: %w", err)
	}

	// Make API request
	err = makeAPIRequest(
		"POST",
		60,
		fmt.Sprintf("%s/rest-api/ssd-vps/locations/%d/servers", getAPIBaseURL(), locationID),
		bytes.NewBuffer(requestBody),
		&response,
	)

	if err != nil {
		return VPSDeployResponse{}, err
	}

	return response, nil
}

// printVPSDetails prints detailed VPS information
func printVPSDetails(cmd *cobra.Command, vps *VPSDetail) {
	// General Information Section
	fmt.Printf("%s %s\n", BlueHeading("Hostname:"), WhiteText(vps.Hostname))
	fmt.Printf("%s %s\n", BlueHeading("Friendly Name:"), WhiteText(vps.Label))

	// Operating System
	osInfo := vps.OperatingSystem
	if vps.OperatingSystemDistro != "" {
		osInfo = fmt.Sprintf("%s (%s)", vps.OperatingSystemDistro, vps.OperatingSystem)
	}
	fmt.Printf("%s %s\n", BlueHeading("Operating System:"), WhiteText(osInfo))

	// Status Information
	statusText := "Active"
	if vps.Suspended {
		statusText = "Suspended"
	}
	fmt.Printf("%s %s\n", BlueHeading("Status:"), WhiteText(statusText))

	// Billing Information
	if vps.HourlyBilling == 1 {
		fmt.Printf("%s %s\n", BlueHeading("Hourly Billing:"), WhiteText("Yes"))
	}

	if vps.Backups == 1 {
		fmt.Printf("%s %s\n", BlueHeading("Backups:"), WhiteText("Enabled"))
	}

	// Resource Information
	fmt.Println(BlueHeading("\n=== RESOURCES ==="))
	fmt.Printf("%s %d\n", BlueHeading("CPU Cores:"), vps.CPUs)
	fmt.Printf("%s %d MB\n", BlueHeading("RAM:"), vps.Memory)
	fmt.Printf("%s %d GB\n", BlueHeading("Disk:"), vps.TotalDiskSize)

	// Network Information Section
	fmt.Println(BlueHeading("\n=== NETWORK INFORMATION ==="))
	if len(vps.IPAddresses) > 0 {
		primaryIP := vps.IPAddresses[0].IPAddress.Address
		fmt.Printf("%s %s\n", BlueHeading("Primary IP Address:"), WhiteText(primaryIP))

		// Additional IPs
		for i := 1; i < len(vps.IPAddresses); i++ {
			fmt.Printf("%s %s\n", BlueHeading("Secondary IP Address:"), WhiteText(vps.IPAddresses[i].IPAddress.Address))
		}
	}

	// Credentials Section
	fmt.Println(BlueHeading("\n=== CREDENTIALS ==="))
	fmt.Printf("%s %s\n", BlueHeading("Username:"), WhiteText("root"))

	// Redact the root password by default
	showRootPassword, _ := cmd.Flags().GetBool("password")
	password := vps.InitialRootPassword
	if !showRootPassword {
		password = "********"
	}
	fmt.Printf("%s %s\n", BlueHeading("Password:"), WhiteText(password))

	// Power Status Section
	hidePowerStatus, _ := cmd.Flags().GetBool("power")
	if !hidePowerStatus {
		fmt.Println(BlueHeading("\n=== POWER STATUS ==="))

		if vps.Booted {
			fmt.Printf("%s %s\n", BlueHeading("Power State:"), GreenText("Powered On"))
		} else {
			fmt.Printf("%s %s\n", BlueHeading("Power State:"), RedText("Powered Off"))
		}
	}
}
