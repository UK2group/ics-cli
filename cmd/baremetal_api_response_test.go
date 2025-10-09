package cmd

import (
	"encoding/json"
	"testing"
)

func TestFlexibleIDUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:    "string ID",
			input:   `"abc123"`,
			want:    "abc123",
			wantErr: false,
		},
		{
			name:    "numeric ID",
			input:   `12345`,
			want:    "12345",
			wantErr: false,
		},
		{
			name:    "float ID",
			input:   `123.45`,
			want:    "123.45",
			wantErr: false,
		},
		{
			name:    "zero",
			input:   `0`,
			want:    "0",
			wantErr: false,
		},
		{
			name:    "negative number",
			input:   `-123`,
			want:    "-123",
			wantErr: false,
		},
		{
			name:    "large number",
			input:   `9999999999`,
			want:    "9999999999",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var id FlexibleID
			err := json.Unmarshal([]byte(tt.input), &id)

			if (err != nil) != tt.wantErr {
				t.Errorf("FlexibleID.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && id.String() != tt.want {
				t.Errorf("FlexibleID.String() = %v, want %v", id.String(), tt.want)
			}
		})
	}
}

func TestFlexibleIDInt(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		want    int
		wantErr bool
	}{
		{
			name:    "valid integer string",
			value:   "12345",
			want:    12345,
			wantErr: false,
		},
		{
			name:    "zero",
			value:   "0",
			want:    0,
			wantErr: false,
		},
		{
			name:    "negative",
			value:   "-100",
			want:    -100,
			wantErr: false,
		},
		{
			name:    "invalid non-numeric string",
			value:   "abc",
			want:    0,
			wantErr: true,
		},
		{
			name:    "float string",
			value:   "123.45",
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := FlexibleID{value: tt.value}
			got, err := id.Int()

			if (err != nil) != tt.wantErr {
				t.Errorf("FlexibleID.Int() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got != tt.want {
				t.Errorf("FlexibleID.Int() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFlexibleIDString(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{
			name:  "numeric string",
			value: "12345",
			want:  "12345",
		},
		{
			name:  "alphanumeric string",
			value: "abc123",
			want:  "abc123",
		},
		{
			name:  "empty string",
			value: "",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := FlexibleID{value: tt.value}
			if got := id.String(); got != tt.want {
				t.Errorf("FlexibleID.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServerResponseUnmarshal(t *testing.T) {
	jsonData := `{
		"statusCode": 200,
		"message": "Success",
		"data": [
			{
				"id": "server-123",
				"hostname": "server1.example.com",
				"mac_address": "00:11:22:33:44:55",
				"public_ip": "192.0.2.1",
				"service_id": 1001,
				"datacenter_name": "London",
				"friendly_name": "My Server",
				"server_type": "dedicated",
				"hourly_billing": true
			}
		]
	}`

	var response ServerResponse
	err := json.Unmarshal([]byte(jsonData), &response)

	if err != nil {
		t.Fatalf("Failed to unmarshal ServerResponse: %v", err)
	}

	if response.StatusCode != 200 {
		t.Errorf("StatusCode = %d, want 200", response.StatusCode)
	}

	if len(response.Data) != 1 {
		t.Fatalf("Data length = %d, want 1", len(response.Data))
	}

	server := response.Data[0]
	if server.ID != "server-123" {
		t.Errorf("Server ID = %s, want server-123", server.ID)
	}

	if server.Hostname != "server1.example.com" {
		t.Errorf("Hostname = %s, want server1.example.com", server.Hostname)
	}

	if server.ServiceID != 1001 {
		t.Errorf("ServiceID = %d, want 1001", server.ServiceID)
	}

	if !server.HourlyBilling {
		t.Error("HourlyBilling = false, want true")
	}
}

func TestServerDetailResponseUnmarshal(t *testing.T) {
	jsonData := `{
		"statusCode": 200,
		"message": "Success",
		"data": {
			"id": "server-456",
			"hostname": "server2.example.com",
			"mac_address": "AA:BB:CC:DD:EE:FF",
			"public_ip": "192.0.2.2",
			"service_id": 2002,
			"friendly_name": "Production Server",
			"operatingSystemId": "ubuntu-22.04",
			"operating_system": "Ubuntu 22.04 LTS",
			"operating_system_user": "root",
			"operating_system_password": "secretpass123",
			"datacenter": "New York",
			"hourly_billing": false,
			"hourly_price": 0.0,
			"ip_addresses": [],
			"network_port": [],
			"provisioning_status": {
				"isProvisioning": false,
				"statusMessage": "Ready"
			}
		}
	}`

	var response ServerDetailResponse
	err := json.Unmarshal([]byte(jsonData), &response)

	if err != nil {
		t.Fatalf("Failed to unmarshal ServerDetailResponse: %v", err)
	}

	if response.Data.ServerID.String() != "server-456" {
		t.Errorf("ServerID = %s, want server-456", response.Data.ServerID.String())
	}

	if response.Data.OperatingSystemName != "Ubuntu 22.04 LTS" {
		t.Errorf("OperatingSystemName = %s, want Ubuntu 22.04 LTS", response.Data.OperatingSystemName)
	}

	if response.Data.HourlyBilling {
		t.Error("HourlyBilling = true, want false")
	}
}

func TestServerDetailResponseWithNumericID(t *testing.T) {
	jsonData := `{
		"statusCode": 200,
		"message": "Success",
		"data": {
			"id": 12345,
			"hostname": "server3.example.com",
			"mac_address": "11:22:33:44:55:66",
			"public_ip": "192.0.2.3",
			"service_id": 3003,
			"friendly_name": "Test Server",
			"operatingSystemId": "centos-8",
			"operating_system": "CentOS 8",
			"operating_system_user": "admin",
			"operating_system_password": "password",
			"datacenter": "Amsterdam",
			"hourly_billing": true,
			"hourly_price": 1.25,
			"ip_addresses": [],
			"network_port": [],
			"provisioning_status": {
				"isProvisioning": true,
				"statusMessage": "Installing OS"
			}
		}
	}`

	var response ServerDetailResponse
	err := json.Unmarshal([]byte(jsonData), &response)

	if err != nil {
		t.Fatalf("Failed to unmarshal ServerDetailResponse with numeric ID: %v", err)
	}

	if response.Data.ServerID.String() != "12345" {
		t.Errorf("ServerID.String() = %s, want 12345", response.Data.ServerID.String())
	}

	serverIDInt, err := response.Data.ServerID.Int()
	if err != nil {
		t.Fatalf("ServerID.Int() error: %v", err)
	}

	if serverIDInt != 12345 {
		t.Errorf("ServerID.Int() = %d, want 12345", serverIDInt)
	}

	if response.Data.HourlyPrice != 1.25 {
		t.Errorf("HourlyPrice = %f, want 1.25", response.Data.HourlyPrice)
	}

	if !response.Data.ProvisioningStatus.IsProvisioning {
		t.Error("IsProvisioning = false, want true")
	}
}

func TestPowerStatusResponseUnmarshal(t *testing.T) {
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
			var response PowerStatusResponse
			err := json.Unmarshal([]byte(tt.jsonData), &response)

			if err != nil {
				t.Fatalf("Failed to unmarshal PowerStatusResponse: %v", err)
			}

			if response.Data.IsPoweredOn != tt.wantPowered {
				t.Errorf("IsPoweredOn = %v, want %v", response.Data.IsPoweredOn, tt.wantPowered)
			}
		})
	}
}

func TestGenericServerResponseUnmarshal(t *testing.T) {
	jsonData := `{
		"statusCode": 200,
		"message": "Operation completed",
		"data": {
			"success": true
		}
	}`

	var response GenericServerResponse
	err := json.Unmarshal([]byte(jsonData), &response)

	if err != nil {
		t.Fatalf("Failed to unmarshal GenericServerResponse: %v", err)
	}

	if !response.Data.Success {
		t.Error("Success = false, want true")
	}

	if response.Message != "Operation completed" {
		t.Errorf("Message = %s, want 'Operation completed'", response.Message)
	}
}

func TestOSListResponseUnmarshal(t *testing.T) {
	jsonData := `{
		"statusCode": 200,
		"message": "Success",
		"data": {
			"osList": [
				{
					"id": "ubuntu-22.04",
					"name": "Ubuntu",
					"version": "22.04 LTS",
					"licenses": ["license1"]
				},
				{
					"id": "centos-8",
					"name": "CentOS",
					"version": "8",
					"licenses": []
				}
			]
		}
	}`

	var response OSListResponse
	err := json.Unmarshal([]byte(jsonData), &response)

	if err != nil {
		t.Fatalf("Failed to unmarshal OSListResponse: %v", err)
	}

	if len(response.Data.OSList) != 2 {
		t.Fatalf("OSList length = %d, want 2", len(response.Data.OSList))
	}

	ubuntu := response.Data.OSList[0]
	if ubuntu.ID != "ubuntu-22.04" {
		t.Errorf("OS ID = %s, want ubuntu-22.04", ubuntu.ID)
	}

	if ubuntu.Name != "Ubuntu" {
		t.Errorf("OS Name = %s, want Ubuntu", ubuntu.Name)
	}

	if len(ubuntu.Licenses) != 1 {
		t.Errorf("Licenses length = %d, want 1", len(ubuntu.Licenses))
	}
}

func TestInventoryResponseUnmarshal(t *testing.T) {
	jsonData := `{
		"statusCode": 200,
		"message": "Success",
		"data": [
			{
				"sku_id": 100,
				"quantity": 5,
				"auto_provision_quantity": 2,
				"datacenter_id": 1,
				"region_id": 10,
				"location_code": "LHR",
				"cpu_brand": "Intel",
				"cpu_model": "Xeon E-2136",
				"cpu_clock_speed_ghz": 3.3,
				"cpu_cores": 6,
				"cpu_count": 1,
				"total_ssd_size_gb": 480,
				"total_hdd_size_gb": 0,
				"total_nvme_size_gb": 0,
				"raid_enabled": true,
				"total_ram_gb": 32,
				"nic_speed_mbps": 1000,
				"qt_product_id": 500,
				"status": "available",
				"metadata": [],
				"currency_code": "USD",
				"sku_product_name": "E-2136",
				"price": "150.00",
				"hourly_enabled": true,
				"price_hourly": "0.25"
			}
		]
	}`

	var response InventoryResponse
	err := json.Unmarshal([]byte(jsonData), &response)

	if err != nil {
		t.Fatalf("Failed to unmarshal InventoryResponse: %v", err)
	}

	if len(response.Data) != 1 {
		t.Fatalf("Data length = %d, want 1", len(response.Data))
	}

	item := response.Data[0]
	if item.SkuProductName != "E-2136" {
		t.Errorf("SkuProductName = %s, want E-2136", item.SkuProductName)
	}

	if item.LocationCode != "LHR" {
		t.Errorf("LocationCode = %s, want LHR", item.LocationCode)
	}

	if item.Quantity != 5 {
		t.Errorf("Quantity = %d, want 5", item.Quantity)
	}

	if !item.HourlyEnabled {
		t.Error("HourlyEnabled = false, want true")
	}
}
