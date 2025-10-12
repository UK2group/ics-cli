package cmd

import (
	"encoding/json"
	"testing"
)

func TestVPSResponseUnmarshal(t *testing.T) {
	jsonData := `{
		"statusCode": 200,
		"message": "Success",
		"data": [
			{
				"id": 123,
				"service_id": 1001,
				"label": "Web Server",
				"hostname": "web1",
				"ip_address": "192.0.2.1",
				"backups": 1,
				"hourly_billing": 1,
				"plan_id": 5,
				"location": {
					"id": 93,
					"name": "Chicago - G, US"
				}
			}
		]
	}`

	var response VPSResponse
	err := json.Unmarshal([]byte(jsonData), &response)

	if err != nil {
		t.Fatalf("Failed to unmarshal VPSResponse: %v", err)
	}

	if response.StatusCode != 200 {
		t.Errorf("StatusCode = %d, want 200", response.StatusCode)
	}

	if len(response.Data) != 1 {
		t.Fatalf("Data length = %d, want 1", len(response.Data))
	}

	vps := response.Data[0]
	if vps.ID != 123 {
		t.Errorf("VPS ID = %d, want 123", vps.ID)
	}

	if vps.ServiceID != 1001 {
		t.Errorf("ServiceID = %d, want 1001", vps.ServiceID)
	}

	if vps.Hostname != "web1" {
		t.Errorf("Hostname = %s, want web1", vps.Hostname)
	}

	if vps.Label != "Web Server" {
		t.Errorf("Label = %s, want Web Server", vps.Label)
	}

	if vps.HourlyBilling != 1 {
		t.Errorf("HourlyBilling = %d, want 1", vps.HourlyBilling)
	}

	if vps.Location.ID != 93 {
		t.Errorf("Location.ID = %d, want 93", vps.Location.ID)
	}

	if vps.Location.Name != "Chicago - G, US" {
		t.Errorf("Location.Name = %s, want Chicago - G, US", vps.Location.Name)
	}
}

func TestVPSDetailResponseUnmarshal(t *testing.T) {
	jsonData := `{
		"statusCode": 200,
		"message": "Success",
		"data": {
			"id": 456,
			"hostname": "web2",
			"memory": 2048,
			"cpus": 2,
			"created_at": "2024-01-01T00:00:00Z",
			"built": true,
			"booted": true,
			"label": "Production Server",
			"operating_system": "linux",
			"operating_system_distro": "Ubuntu 22.04",
			"note": "Production web server",
			"suspended": false,
			"domain": "example.com",
			"total_disk_size": 50,
			"initial_root_password": "secretpass123",
			"hourly_billing": 1,
			"backups": 0,
			"ip_addresses": [
				{
					"ip_address": {
						"address": "192.0.2.10",
						"broadcast": "192.0.2.255",
						"network_address": "192.0.2.0",
						"gateway": "192.0.2.1",
						"created_at": "2024-01-01T00:00:00Z",
						"rdnsEntry": "web2.example.com"
					}
				}
			]
		}
	}`

	var response VPSDetailResponse
	err := json.Unmarshal([]byte(jsonData), &response)

	if err != nil {
		t.Fatalf("Failed to unmarshal VPSDetailResponse: %v", err)
	}

	if response.Data.ID != 456 {
		t.Errorf("ID = %d, want 456", response.Data.ID)
	}

	if response.Data.Hostname != "web2" {
		t.Errorf("Hostname = %s, want web2", response.Data.Hostname)
	}

	if response.Data.Memory != 2048 {
		t.Errorf("Memory = %d, want 2048", response.Data.Memory)
	}

	if response.Data.CPUs != 2 {
		t.Errorf("CPUs = %d, want 2", response.Data.CPUs)
	}

	if !response.Data.Built {
		t.Error("Built = false, want true")
	}

	if !response.Data.Booted {
		t.Error("Booted = false, want true")
	}

	if response.Data.OperatingSystemDistro != "Ubuntu 22.04" {
		t.Errorf("OperatingSystemDistro = %s, want Ubuntu 22.04", response.Data.OperatingSystemDistro)
	}

	if response.Data.TotalDiskSize != 50 {
		t.Errorf("TotalDiskSize = %d, want 50", response.Data.TotalDiskSize)
	}

	if response.Data.HourlyBilling != 1 {
		t.Errorf("HourlyBilling = %d, want 1", response.Data.HourlyBilling)
	}

	if len(response.Data.IPAddresses) != 1 {
		t.Fatalf("IPAddresses length = %d, want 1", len(response.Data.IPAddresses))
	}

	if response.Data.IPAddresses[0].IPAddress.Address != "192.0.2.10" {
		t.Errorf("IP Address = %s, want 192.0.2.10", response.Data.IPAddresses[0].IPAddress.Address)
	}
}

func TestVPSPowerStatusResponseUnmarshal(t *testing.T) {
	tests := []struct {
		name        string
		jsonData    string
		wantPowered bool
	}{
		{
			name: "powered on",
			jsonData: `{
				"statusCode": 200,
				"message": "Success",
				"data": {
					"is_powered_on": true
				}
			}`,
			wantPowered: true,
		},
		{
			name: "powered off",
			jsonData: `{
				"statusCode": 200,
				"message": "Success",
				"data": {
					"is_powered_on": false
				}
			}`,
			wantPowered: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var response VPSPowerStatusResponse
			err := json.Unmarshal([]byte(tt.jsonData), &response)

			if err != nil {
				t.Fatalf("Failed to unmarshal VPSPowerStatusResponse: %v", err)
			}

			if response.Data.IsPoweredOn != tt.wantPowered {
				t.Errorf("IsPoweredOn = %v, want %v", response.Data.IsPoweredOn, tt.wantPowered)
			}
		})
	}
}

func TestVPSEmptyDataResponseUnmarshal(t *testing.T) {
	jsonData := `{
		"status_code": 200,
		"message": "Virtual Machine Rebooted",
		"data": []
	}`

	var response VPSEmptyDataResponse
	err := json.Unmarshal([]byte(jsonData), &response)

	if err != nil {
		t.Fatalf("Failed to unmarshal VPSEmptyDataResponse: %v", err)
	}

	if response.StatusCode != 200 {
		t.Errorf("StatusCode = %d, want 200", response.StatusCode)
	}

	if response.Message != "Virtual Machine Rebooted" {
		t.Errorf("Message = %s, want Virtual Machine Rebooted", response.Message)
	}

	if len(response.Data) != 0 {
		t.Errorf("Data length = %d, want 0", len(response.Data))
	}
}

func TestVPSPlanResponseUnmarshal(t *testing.T) {
	jsonData := `{
		"statusCode": 200,
		"message": "Success",
		"data": [
			{
				"id": 1,
				"label": "2 GB RAM, 1 CPU Core, 50 GB SSD",
				"price": "10.00",
				"product_name": "V1"
			},
			{
				"id": 2,
				"label": "4 GB RAM, 2 CPU Cores, 80 GB SSD",
				"price": "20.00",
				"product_name": "V2"
			}
		]
	}`

	var response VPSPlanResponse
	err := json.Unmarshal([]byte(jsonData), &response)

	if err != nil {
		t.Fatalf("Failed to unmarshal VPSPlanResponse: %v", err)
	}

	if len(response.Data) != 2 {
		t.Fatalf("Data length = %d, want 2", len(response.Data))
	}

	plan := response.Data[0]
	if plan.ProductName != "V1" {
		t.Errorf("ProductName = %s, want V1", plan.ProductName)
	}

	if plan.Price != "10.00" {
		t.Errorf("Price = %s, want 10.00", plan.Price)
	}

	if plan.Label != "2 GB RAM, 1 CPU Core, 50 GB SSD" {
		t.Errorf("Label = %s, want 2 GB RAM, 1 CPU Core, 50 GB SSD", plan.Label)
	}
}

func TestVPSOSResponseUnmarshal(t *testing.T) {
	jsonData := `{
		"statusCode": 200,
		"message": "Success",
		"data": [
			{
				"name": "Ubuntu 22.04",
				"component_code": "SSDVPSUBUNTU22",
				"price": 0.0,
				"price_hourly": 0.0,
				"currency": "USD",
				"min_disk_size": "20",
				"operating_system_arch": "x64",
				"min_memory_size": "1024"
			},
			{
				"name": "Debian 11",
				"component_code": "SSDVPSDEBIAN11",
				"price": 5.00,
				"price_hourly": 0.01,
				"currency": "USD",
				"min_disk_size": "20",
				"operating_system_arch": "x64",
				"min_memory_size": "1024"
			}
		]
	}`

	var response VPSOSResponse
	err := json.Unmarshal([]byte(jsonData), &response)

	if err != nil {
		t.Fatalf("Failed to unmarshal VPSOSResponse: %v", err)
	}

	if len(response.Data) != 2 {
		t.Fatalf("Data length = %d, want 2", len(response.Data))
	}

	os := response.Data[0]
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

func TestVPSLocationResponseUnmarshal(t *testing.T) {
	jsonData := `{
		"statusCode": 200,
		"message": "Success",
		"data": [
			{
				"id": 93,
				"name": "Chicago - G, US"
			},
			{
				"id": 94,
				"name": "New York City - A, US"
			}
		]
	}`

	var response VPSLocationResponse
	err := json.Unmarshal([]byte(jsonData), &response)

	if err != nil {
		t.Fatalf("Failed to unmarshal VPSLocationResponse: %v", err)
	}

	if len(response.Data) != 2 {
		t.Fatalf("Data length = %d, want 2", len(response.Data))
	}

	loc := response.Data[0]
	if loc.ID != 93 {
		t.Errorf("ID = %d, want 93", loc.ID)
	}

	if loc.Name != "Chicago - G, US" {
		t.Errorf("Name = %s, want Chicago - G, US", loc.Name)
	}
}

func TestVPSDeployRequestMarshal(t *testing.T) {
	request := VPSDeployRequest{
		ProductName:     "V2",
		Label:           "Web Server",
		Hostname:        "web1",
		Backups:         true,
		BillHourly:      true,
		OSComponentCode: "SSDVPSUBUNTU22",
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("Failed to marshal VPSDeployRequest: %v", err)
	}

	var unmarshaled VPSDeployRequest
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal VPSDeployRequest: %v", err)
	}

	if unmarshaled.ProductName != "V2" {
		t.Errorf("ProductName = %s, want V2", unmarshaled.ProductName)
	}

	if unmarshaled.Hostname != "web1" {
		t.Errorf("Hostname = %s, want web1", unmarshaled.Hostname)
	}

	if !unmarshaled.Backups {
		t.Error("Backups = false, want true")
	}

	if !unmarshaled.BillHourly {
		t.Error("BillHourly = false, want true")
	}
}

func TestVPSDeployResponseUnmarshal(t *testing.T) {
	jsonData := `{
		"status_code": 201,
		"message": "Order Created Successfully, the server is being provisioned",
		"data": []
	}`

	var response VPSDeployResponse
	err := json.Unmarshal([]byte(jsonData), &response)

	if err != nil {
		t.Fatalf("Failed to unmarshal VPSDeployResponse: %v", err)
	}

	if response.StatusCode != 201 {
		t.Errorf("StatusCode = %d, want 201", response.StatusCode)
	}

	if response.Message != "Order Created Successfully, the server is being provisioned" {
		t.Errorf("Message = %s, want Order Created Successfully, the server is being provisioned", response.Message)
	}

	if len(response.Data) != 0 {
		t.Errorf("Data length = %d, want 0", len(response.Data))
	}
}
