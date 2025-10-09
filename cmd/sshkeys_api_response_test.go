package cmd

import (
	"encoding/json"
	"testing"
)

func TestSSHKeysResponseUnmarshal(t *testing.T) {
	jsonData := `{
		"statusCode": 200,
		"message": "Success",
		"data": [
			{
				"id": 1,
				"label": "My SSH Key",
				"key": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC user@host",
				"created_at": 1640000000,
				"updated_at": 1640000100,
				"assigned_servers": [
					{
						"server_id": "server-123",
						"service_id": 1001,
						"domain": "example.com",
						"hostname": "server1.example.com",
						"datacenter_name": "London"
					}
				]
			}
		]
	}`

	var response SSHKeysResponse
	err := json.Unmarshal([]byte(jsonData), &response)

	if err != nil {
		t.Fatalf("Failed to unmarshal SSHKeysResponse: %v", err)
	}

	if response.StatusCode != 200 {
		t.Errorf("StatusCode = %d, want 200", response.StatusCode)
	}

	if len(response.Data) != 1 {
		t.Fatalf("Data length = %d, want 1", len(response.Data))
	}

	key := response.Data[0]
	if key.ID != 1 {
		t.Errorf("ID = %d, want 1", key.ID)
	}

	if key.Label != "My SSH Key" {
		t.Errorf("Label = %s, want 'My SSH Key'", key.Label)
	}

	if len(key.AssignedServers) != 1 {
		t.Fatalf("AssignedServers length = %d, want 1", len(key.AssignedServers))
	}

	server := key.AssignedServers[0]
	if server.ServerID != "server-123" {
		t.Errorf("ServerID = %s, want server-123", server.ServerID)
	}
}

func TestSSHKeyAddResponseUnmarshal(t *testing.T) {
	jsonData := `{
		"statusCode": 201,
		"message": "SSH key created successfully",
		"data": {
			"id": 42
		}
	}`

	var response SSHKeyAddResponse
	err := json.Unmarshal([]byte(jsonData), &response)

	if err != nil {
		t.Fatalf("Failed to unmarshal SSHKeyAddResponse: %v", err)
	}

	if response.StatusCode != 201 {
		t.Errorf("StatusCode = %d, want 201", response.StatusCode)
	}

	if response.Data.ID != 42 {
		t.Errorf("Data.ID = %d, want 42", response.Data.ID)
	}

	if response.Message != "SSH key created successfully" {
		t.Errorf("Message = %s, want 'SSH key created successfully'", response.Message)
	}
}

func TestSSHKeyDeleteResponseUnmarshal(t *testing.T) {
	tests := []struct {
		name        string
		jsonData    string
		wantSuccess bool
	}{
		{
			name: "successful deletion",
			jsonData: `{
				"statusCode": 200,
				"message": "SSH key deleted",
				"data": true
			}`,
			wantSuccess: true,
		},
		{
			name: "failed deletion",
			jsonData: `{
				"statusCode": 200,
				"message": "Failed to delete",
				"data": false
			}`,
			wantSuccess: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var response SSHKeyDeleteResponse
			err := json.Unmarshal([]byte(tt.jsonData), &response)

			if err != nil {
				t.Fatalf("Failed to unmarshal SSHKeyDeleteResponse: %v", err)
			}

			if response.Data != tt.wantSuccess {
				t.Errorf("Data = %v, want %v", response.Data, tt.wantSuccess)
			}
		})
	}
}

func TestSSHKeyAssignResponseUnmarshal(t *testing.T) {
	jsonData := `{
		"statusCode": 200,
		"message": "SSH keys assigned successfully",
		"data": [
			{
				"id": 1,
				"label": "Key 1",
				"key": "ssh-rsa AAAAB3... user1@host",
				"created_at": 1640000000,
				"updated_at": 1640000100
			},
			{
				"id": 2,
				"label": "Key 2",
				"key": "ssh-ed25519 AAAAC3... user2@host",
				"created_at": 1640000200,
				"updated_at": 1640000300
			}
		]
	}`

	var response SSHKeyAssignResponse
	err := json.Unmarshal([]byte(jsonData), &response)

	if err != nil {
		t.Fatalf("Failed to unmarshal SSHKeyAssignResponse: %v", err)
	}

	if response.StatusCode != 200 {
		t.Errorf("StatusCode = %d, want 200", response.StatusCode)
	}

	if len(response.Data) != 2 {
		t.Fatalf("Data length = %d, want 2", len(response.Data))
	}

	key1 := response.Data[0]
	if key1.ID != 1 {
		t.Errorf("Data[0].ID = %d, want 1", key1.ID)
	}

	if key1.Label != "Key 1" {
		t.Errorf("Data[0].Label = %s, want 'Key 1'", key1.Label)
	}

	key2 := response.Data[1]
	if key2.ID != 2 {
		t.Errorf("Data[1].ID = %d, want 2", key2.ID)
	}
}

func TestSSHKeyUnassignResponseUnmarshal(t *testing.T) {
	jsonData := `{
		"statusCode": 200,
		"message": "SSH key unassigned successfully",
		"data": true
	}`

	var response SSHKeyUnassignResponse
	err := json.Unmarshal([]byte(jsonData), &response)

	if err != nil {
		t.Fatalf("Failed to unmarshal SSHKeyUnassignResponse: %v", err)
	}

	if response.StatusCode != 200 {
		t.Errorf("StatusCode = %d, want 200", response.StatusCode)
	}

	if !response.Data {
		t.Error("Data = false, want true")
	}
}

func TestAssignedSSHKeysResponseUnmarshal(t *testing.T) {
	jsonData := `{
		"statusCode": 200,
		"message": "Success",
		"data": [
			{
				"id": 10,
				"label": "Production Key",
				"key": "ssh-rsa AAAAB3... prod@server",
				"created_at": 1640000000,
				"updated_at": 1640000100
			},
			{
				"id": 20,
				"label": "Development Key",
				"key": "ssh-ed25519 AAAAC3... dev@server",
				"created_at": 1640000200,
				"updated_at": 1640000300
			}
		]
	}`

	var response AssignedSSHKeysResponse
	err := json.Unmarshal([]byte(jsonData), &response)

	if err != nil {
		t.Fatalf("Failed to unmarshal AssignedSSHKeysResponse: %v", err)
	}

	if len(response.Data) != 2 {
		t.Fatalf("Data length = %d, want 2", len(response.Data))
	}

	key1 := response.Data[0]
	if key1.ID != 10 {
		t.Errorf("Data[0].ID = %d, want 10", key1.ID)
	}

	if key1.Label != "Production Key" {
		t.Errorf("Data[0].Label = %s, want 'Production Key'", key1.Label)
	}

	key2 := response.Data[1]
	if key2.ID != 20 {
		t.Errorf("Data[1].ID = %d, want 20", key2.ID)
	}

	if key2.Label != "Development Key" {
		t.Errorf("Data[1].Label = %s, want 'Development Key'", key2.Label)
	}
}

func TestSSHKeyResponseWithEmptyAssignedServers(t *testing.T) {
	jsonData := `{
		"statusCode": 200,
		"message": "Success",
		"data": [
			{
				"id": 5,
				"label": "Unassigned Key",
				"key": "ssh-rsa AAAAB3... user@host",
				"created_at": 1640000000,
				"updated_at": 1640000100,
				"assigned_servers": []
			}
		]
	}`

	var response SSHKeysResponse
	err := json.Unmarshal([]byte(jsonData), &response)

	if err != nil {
		t.Fatalf("Failed to unmarshal SSHKeysResponse: %v", err)
	}

	if len(response.Data) != 1 {
		t.Fatalf("Data length = %d, want 1", len(response.Data))
	}

	key := response.Data[0]
	if len(key.AssignedServers) != 0 {
		t.Errorf("AssignedServers length = %d, want 0", len(key.AssignedServers))
	}
}

func TestSSHKeyResponseWithMultipleAssignedServers(t *testing.T) {
	jsonData := `{
		"statusCode": 200,
		"message": "Success",
		"data": [
			{
				"id": 3,
				"label": "Shared Key",
				"key": "ssh-rsa AAAAB3... user@host",
				"created_at": 1640000000,
				"updated_at": 1640000100,
				"assigned_servers": [
					{
						"server_id": "server-1",
						"service_id": 1001,
						"domain": "example1.com",
						"hostname": "server1.example.com",
						"datacenter_name": "London"
					},
					{
						"server_id": "server-2",
						"service_id": 1002,
						"domain": "example2.com",
						"hostname": "server2.example.com",
						"datacenter_name": "New York"
					},
					{
						"server_id": "server-3",
						"service_id": 1003,
						"domain": "example3.com",
						"hostname": "server3.example.com",
						"datacenter_name": "Amsterdam"
					}
				]
			}
		]
	}`

	var response SSHKeysResponse
	err := json.Unmarshal([]byte(jsonData), &response)

	if err != nil {
		t.Fatalf("Failed to unmarshal SSHKeysResponse: %v", err)
	}

	key := response.Data[0]
	if len(key.AssignedServers) != 3 {
		t.Fatalf("AssignedServers length = %d, want 3", len(key.AssignedServers))
	}

	// Verify all three servers
	expectedDatacenters := []string{"London", "New York", "Amsterdam"}
	for i, server := range key.AssignedServers {
		if server.DatacenterName != expectedDatacenters[i] {
			t.Errorf("AssignedServers[%d].DatacenterName = %s, want %s",
				i, server.DatacenterName, expectedDatacenters[i])
		}
	}
}

func TestSSHKeyTimestamps(t *testing.T) {
	jsonData := `{
		"statusCode": 200,
		"message": "Success",
		"data": [
			{
				"id": 7,
				"label": "Time Test Key",
				"key": "ssh-rsa AAAAB3... user@host",
				"created_at": 1672531200,
				"updated_at": 1704067200,
				"assigned_servers": []
			}
		]
	}`

	var response SSHKeysResponse
	err := json.Unmarshal([]byte(jsonData), &response)

	if err != nil {
		t.Fatalf("Failed to unmarshal SSHKeysResponse: %v", err)
	}

	key := response.Data[0]
	if key.CreatedAt != 1672531200 {
		t.Errorf("CreatedAt = %d, want 1672531200", key.CreatedAt)
	}

	if key.UpdatedAt != 1704067200 {
		t.Errorf("UpdatedAt = %d, want 1704067200", key.UpdatedAt)
	}

	// Verify updated_at is after created_at
	if key.UpdatedAt < key.CreatedAt {
		t.Error("UpdatedAt should be >= CreatedAt")
	}
}
