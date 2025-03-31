package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

// getSSHKeyFromLabel finds the SSH Key ID from the corresponding Label
func getSSHKeyFromLabel(sshKeyLabel string) (*SSHKey, error) {
	// Use makeAPIRequest to get the ssh key listt
	var sshKeyList SSHKeysResponse

	err := makeAPIRequest(
		"GET",
		30,
		"https://api.ingenuitycloudservices.com/rest-api/ssh-keys",
		nil,
		&sshKeyList,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get ssh key list: %w", err)
	}

	// Find the server with the matching service ID
	for i, sshKey := range sshKeyList.Data {
		if sshKey.Label == sshKeyLabel {
			return &sshKeyList.Data[i], nil
		}
	}

	return nil, fmt.Errorf("SSH Key with name %s not found", sshKeyLabel)
}

// cleanSSHKey removes extra whitespace, newlines, and comments from an SSH key
func cleanSSHKey(key string) string {
	// Remove leading/trailing whitespace
	key = strings.TrimSpace(key)

	// Split into lines and keep only the key line
	lines := strings.Split(key, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "ssh-rsa") ||
			strings.HasPrefix(line, "ssh-ed25519") ||
			strings.HasPrefix(line, "ssh-dss") ||
			strings.HasPrefix(line, "ecdsa-sha2-") {
			return line
		}
	}

	// If we didn't find a key line, return the original (trimmed)
	return key
}

// addSSHKey sends a request to add a new SSH key
func addSSHKey(requestData map[string]string) (bool, error) {

	// Create request body
	requestBody, err := json.Marshal(requestData)
	if err != nil {
		return false, fmt.Errorf("error creating request body: %w", err)
	}

	// Make API request
	var response SSHKeyAddResponse
	err = makeAPIRequest(
		"POST",
		60,
		"https://api.ingenuitycloudservices.com/rest-api/ssh-keys",
		bytes.NewBuffer(requestBody),
		&response,
	)

	if err != nil {
		return false, err
	}

	return true, nil
}

// deleteSSHKey sends a request to delete an SSH key
func deleteSSHKey(keyID int) (bool, error) {
	var response GenericServerResponse

	err := makeAPIRequest(
		"DELETE",
		60,
		fmt.Sprintf("https://api.ingenuitycloudservices.com/rest-api/ssh-keys/%d", keyID),
		nil,
		&response,
	)

	if err != nil {
		return false, fmt.Errorf("%w", err)
	}

	return true, nil
}

// getSSHKeys retrieves all SSH keys from the API
func getSSHKeys() ([]SSHKey, error) {
	var response SSHKeysResponse

	err := makeAPIRequest(
		"GET",
		60,
		"https://api.ingenuitycloudservices.com/rest-api/ssh-keys",
		nil,
		&response,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get SSH keys: %w", err)
	}

	return response.Data, nil
}

// setSSHKeyLabel updates the label of an SSH key
func setSSHKeyLabel(keyID int, newName string) (bool, error) {
	var response GenericServerResponse

	// Pass url as pxe_script_url in the request body
	requestData := map[string]string{
		"label": newName,
	}

	requestBody, err := json.Marshal(requestData)
	if err != nil {
		return false, fmt.Errorf("error creating request body: %w", err)
	}

	err = makeAPIRequest(
		"PUT",
		60,
		fmt.Sprintf("https://api.ingenuitycloudservices.com/rest-api/ssh-keys/%d", keyID),
		bytes.NewBuffer(requestBody),
		&response,
	)

	if err != nil {
		return false, fmt.Errorf("%w", err)
	}

	return true, nil
}

// assignSSHKey assigns an SSH key to a Baremetal Server
func assignSSHKeys(serverID string, keyIDs []int) (bool, error) {
	var response SSHKeyAssignResponse

	// Create request data with an array of SSH key IDs
	requestData := map[string][]int{
		"ssh_key_ids": keyIDs, // Pass the entire array of keyIDs
	}

	// Create request body
	requestBody, err := json.Marshal(requestData)
	if err != nil {
		return false, fmt.Errorf("error creating request body: %w", err)
	}

	// Make API request
	err = makeAPIRequest(
		"PATCH",
		60,
		fmt.Sprintf("https://api.ingenuitycloudservices.com/rest-api/servers/%s/ssh-keys/assign", serverID),
		bytes.NewBuffer(requestBody),
		&response,
	)

	if err != nil {
		return false, fmt.Errorf("failed to assign SSH keys: %w", err)
	}

	return response.Data, nil
}

// unassignSSHKey sends a request to unassign an SSH key from a server
func unassignSSHKey(serverID string, keyID int) (bool, error) {
	var response SSHKeysResponse

	// Create request data with an array of SSH key IDs
	requestData := map[string][]int{
		"ssh_key_ids": {keyID}, // Array with single keyID
	}

	// Create request body
	requestBody, err := json.Marshal(requestData)
	if err != nil {
		return false, fmt.Errorf("error creating request body: %w", err)
	}

	// Make API request
	err = makeAPIRequest(
		"PATCH",
		60,
		fmt.Sprintf("https://api.ingenuitycloudservices.com/rest-api/servers/%s/ssh-keys/un-assign", serverID),
		bytes.NewBuffer(requestBody),
		&response,
	)

	if err != nil {
		return false, fmt.Errorf("failed to unassign SSH key: %w", err)
	}

	return true, nil
}
