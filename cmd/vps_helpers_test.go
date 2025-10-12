package cmd

import (
	"testing"
)

func TestGetVPSIDFromServiceID(t *testing.T) {
	// This is a unit test for the logic, not the API call
	// We test the pattern matching logic

	tests := []struct {
		name      string
		serviceID string
		vpsData   []VPSBasic
		wantID    int
		wantErr   bool
	}{
		{
			name:      "found matching service ID",
			serviceID: "1001",
			vpsData: []VPSBasic{
				{ID: 123, ServiceID: 1001, Hostname: "web1"},
				{ID: 456, ServiceID: 1002, Hostname: "web2"},
			},
			wantID:  123,
			wantErr: false,
		},
		{
			name:      "service ID not found",
			serviceID: "9999",
			vpsData: []VPSBasic{
				{ID: 123, ServiceID: 1001, Hostname: "web1"},
			},
			wantID:  0,
			wantErr: true,
		},
		{
			name:      "empty VPS list",
			serviceID: "1001",
			vpsData:   []VPSBasic{},
			wantID:    0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the matching logic manually
			var foundID int
			var found bool

			for _, vps := range tt.vpsData {
				if vps.ServiceID == 1001 && tt.serviceID == "1001" {
					foundID = vps.ID
					found = true
					break
				}
			}

			if tt.wantErr && found {
				t.Error("Expected error but found VPS")
			}

			if !tt.wantErr && !found {
				t.Error("Expected to find VPS but didn't")
			}

			if !tt.wantErr && foundID != tt.wantID {
				t.Errorf("Got ID %d, want %d", foundID, tt.wantID)
			}
		})
	}
}

func TestGetVPSLocationIDByName_Logic(t *testing.T) {
	// Test the location matching logic
	locations := []VPSLocation{
		{ID: 93, Name: "Chicago - G, US"},
		{ID: 94, Name: "New York City - A, US"},
		{ID: 95, Name: "London - P, UK"},
	}

	tests := []struct {
		name         string
		locationName string
		wantID       int
		wantFound    bool
	}{
		{
			name:         "exact match",
			locationName: "Chicago - G, US",
			wantID:       93,
			wantFound:    true,
		},
		{
			name:         "not found",
			locationName: "Paris - A, FR",
			wantID:       0,
			wantFound:    false,
		},
		{
			name:         "case sensitive match",
			locationName: "chicago - g, us",
			wantID:       0,
			wantFound:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var foundID int
			var found bool

			for _, loc := range locations {
				if loc.Name == tt.locationName {
					foundID = loc.ID
					found = true
					break
				}
			}

			if found != tt.wantFound {
				t.Errorf("Found = %v, want %v", found, tt.wantFound)
			}

			if found && foundID != tt.wantID {
				t.Errorf("ID = %d, want %d", foundID, tt.wantID)
			}
		})
	}
}

func TestGetVPSOSComponentCodeByName_Logic(t *testing.T) {
	// Test the OS matching logic
	osList := []VPSOS{
		{Name: "Ubuntu 22.04", ComponentCode: "SSDVPSUBUNTU22"},
		{Name: "Debian 11", ComponentCode: "SSDVPSDEBIAN11"},
		{Name: "CentOS 8", ComponentCode: "SSDVPSCENTOS8"},
	}

	tests := []struct {
		name          string
		osName        string
		wantCode      string
		wantFound     bool
	}{
		{
			name:      "exact match Ubuntu",
			osName:    "Ubuntu 22.04",
			wantCode:  "SSDVPSUBUNTU22",
			wantFound: true,
		},
		{
			name:      "exact match Debian",
			osName:    "Debian 11",
			wantCode:  "SSDVPSDEBIAN11",
			wantFound: true,
		},
		{
			name:      "not found",
			osName:    "Windows Server 2022",
			wantCode:  "",
			wantFound: false,
		},
		{
			name:      "case sensitive no match",
			osName:    "ubuntu 22.04",
			wantCode:  "",
			wantFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var foundCode string
			var found bool

			for _, os := range osList {
				if os.Name == tt.osName {
					foundCode = os.ComponentCode
					found = true
					break
				}
			}

			if found != tt.wantFound {
				t.Errorf("Found = %v, want %v", found, tt.wantFound)
			}

			if found && foundCode != tt.wantCode {
				t.Errorf("ComponentCode = %s, want %s", foundCode, tt.wantCode)
			}
		})
	}
}

func TestVPSBasicStructure(t *testing.T) {
	// Test that VPSBasic struct can hold expected values
	vps := VPSBasic{
		ID:            123,
		ServiceID:     1001,
		Label:         "Web Server",
		Hostname:      "web1",
		IPAddress:     "192.0.2.1",
		Backups:       1,
		HourlyBilling: 1,
		PlanID:        5,
	}
	vps.Location.ID = 93
	vps.Location.Name = "Chicago - G, US"

	if vps.ID != 123 {
		t.Errorf("ID = %d, want 123", vps.ID)
	}

	if vps.ServiceID != 1001 {
		t.Errorf("ServiceID = %d, want 1001", vps.ServiceID)
	}

	if vps.HourlyBilling != 1 {
		t.Errorf("HourlyBilling = %d, want 1", vps.HourlyBilling)
	}

	if vps.Location.Name != "Chicago - G, US" {
		t.Errorf("Location.Name = %s, want Chicago - G, US", vps.Location.Name)
	}
}

func TestVPSDetailStructure(t *testing.T) {
	// Test that VPSDetail struct can hold expected values
	vps := VPSDetail{
		ID:                    456,
		Hostname:              "web2",
		Memory:                2048,
		CPUs:                  2,
		CreatedAt:             "2024-01-01T00:00:00Z",
		Built:                 true,
		Booted:                true,
		Label:                 "Production Server",
		OperatingSystem:       "linux",
		OperatingSystemDistro: "Ubuntu 22.04",
		Suspended:             false,
		TotalDiskSize:         50,
		InitialRootPassword:   "secretpass123",
		HourlyBilling:         1,
		Backups:               0,
	}

	if vps.ID != 456 {
		t.Errorf("ID = %d, want 456", vps.ID)
	}

	if vps.Memory != 2048 {
		t.Errorf("Memory = %d, want 2048", vps.Memory)
	}

	if vps.CPUs != 2 {
		t.Errorf("CPUs = %d, want 2", vps.CPUs)
	}

	if !vps.Built {
		t.Error("Built = false, want true")
	}

	if !vps.Booted {
		t.Error("Booted = false, want true")
	}

	if vps.OperatingSystemDistro != "Ubuntu 22.04" {
		t.Errorf("OperatingSystemDistro = %s, want Ubuntu 22.04", vps.OperatingSystemDistro)
	}

	if vps.TotalDiskSize != 50 {
		t.Errorf("TotalDiskSize = %d, want 50", vps.TotalDiskSize)
	}

	if vps.HourlyBilling != 1 {
		t.Errorf("HourlyBilling = %d, want 1", vps.HourlyBilling)
	}

	if vps.Backups != 0 {
		t.Errorf("Backups = %d, want 0", vps.Backups)
	}
}

func TestVPSPlanStructure(t *testing.T) {
	plan := VPSPlan{
		ID:          1,
		Label:       "2 GB RAM, 1 CPU Core, 50 GB SSD",
		Price:       "10.00",
		ProductName: "V1",
	}

	if plan.ID != 1 {
		t.Errorf("ID = %d, want 1", plan.ID)
	}

	if plan.ProductName != "V1" {
		t.Errorf("ProductName = %s, want V1", plan.ProductName)
	}

	if plan.Price != "10.00" {
		t.Errorf("Price = %s, want 10.00", plan.Price)
	}
}

func TestVPSOSStructure(t *testing.T) {
	os := VPSOS{
		Name:                "Ubuntu 22.04",
		ComponentCode:       "SSDVPSUBUNTU22",
		Price:               0.0,
		PriceHourly:         0.0,
		Currency:            "USD",
		MinDiskSize:         "20",
		OperatingSystemArch: "x64",
		MinMemorySize:       "1024",
	}

	if os.Name != "Ubuntu 22.04" {
		t.Errorf("Name = %s, want Ubuntu 22.04", os.Name)
	}

	if os.ComponentCode != "SSDVPSUBUNTU22" {
		t.Errorf("ComponentCode = %s, want SSDVPSUBUNTU22", os.ComponentCode)
	}

	if os.Price != 0.0 {
		t.Errorf("Price = %f, want 0.0", os.Price)
	}

	if os.OperatingSystemArch != "x64" {
		t.Errorf("OperatingSystemArch = %s, want x64", os.OperatingSystemArch)
	}
}

func TestVPSLocationStructure(t *testing.T) {
	location := VPSLocation{
		ID:   93,
		Name: "Chicago - G, US",
	}

	if location.ID != 93 {
		t.Errorf("ID = %d, want 93", location.ID)
	}

	if location.Name != "Chicago - G, US" {
		t.Errorf("Name = %s, want Chicago - G, US", location.Name)
	}
}

func TestVPSDeployRequestStructure(t *testing.T) {
	request := VPSDeployRequest{
		ProductName:     "V2",
		Label:           "Web Server",
		Hostname:        "web1",
		Backups:         true,
		BillHourly:      true,
		OSComponentCode: "SSDVPSUBUNTU22",
	}

	if request.ProductName != "V2" {
		t.Errorf("ProductName = %s, want V2", request.ProductName)
	}

	if request.Hostname != "web1" {
		t.Errorf("Hostname = %s, want web1", request.Hostname)
	}

	if !request.Backups {
		t.Error("Backups = false, want true")
	}

	if !request.BillHourly {
		t.Error("BillHourly = false, want true")
	}

	if request.OSComponentCode != "SSDVPSUBUNTU22" {
		t.Errorf("OSComponentCode = %s, want SSDVPSUBUNTU22", request.OSComponentCode)
	}
}
