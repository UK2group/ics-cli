package cmd

// SSHKeysResponse represents the response from the SSH keys API.
// This is returned when listing all SSH keys for a customer.
type SSHKeysResponse struct {
	StatusCode int      `json:"statusCode"` // HTTP status code returned by the API
	Message    string   `json:"message"`    // Human-readable message about the response
	Data       []SSHKey `json:"data"`       // Array of SSH key objects
}

// SSHKey represents an individual SSH key in the API response.
// Contains detailed information about a single SSH key including assigned servers.
type SSHKey struct {
	ID              int         `json:"id"`               // Unique identifier of the SSH key
	Label           string      `json:"label"`            // User-defined label for the SSH key
	Key             string      `json:"key"`              // The actual SSH public key content
	CreatedAt       int64       `json:"created_at"`       // Unix timestamp when the key was created
	UpdatedAt       int64       `json:"updated_at"`       // Unix timestamp when the key was last updated
	AssignedServers []SSHServer `json:"assigned_servers"` // List of servers this key is assigned to
}

// SSHServer represents a server that has the SSH key assigned.
// Contains basic information about a server associated with an SSH key.
type SSHServer struct {
	ServerID       string `json:"server_id"`       // Unique identifier of the server
	ServiceID      int    `json:"service_id"`      // Service ID associated with the server
	Domain         string `json:"domain"`          // Domain name associated with the server
	Hostname       string `json:"hostname"`        // Server hostname
	DatacenterName string `json:"datacenter_name"` // Name of the datacenter where the server is located
}

// SSHKeyAddResponse represents the response from adding a new SSH key.
// This is returned when creating a new SSH key through the API.
type SSHKeyAddResponse struct {
	StatusCode int    `json:"statusCode"` // HTTP status code returned by the API
	Message    string `json:"message"`    // Human-readable message about the response
	Data       struct {
		ID int `json:"id"` // Unique identifier of the newly created SSH key
	} `json:"data"`
}

// SSHKeyDeleteResponse represents the response from deleting an SSH key.
// This is returned when removing an SSH key through the API.
type SSHKeyDeleteResponse struct {
	StatusCode int    `json:"statusCode"` // HTTP status code returned by the API
	Message    string `json:"message"`    // Human-readable message about the response
	Data       bool   `json:"data"`       // Success indicator (true if deletion was successful)
}

// SSHKeyAssignResponse represents the response from assigning an SSH key to a server.
// This is returned when adding an SSH key to a specific server.
type SSHKeyAssignResponse struct {
	StatusCode int    `json:"statusCode"` // HTTP status code returned by the API
	Message    string `json:"message"`    // Human-readable message about the response
	Data       bool   `json:"data"`       // Success indicator (true if assignment was successful)
}

// SSHKeyUnassignResponse represents the response from removing an SSH key from a server.
// This is returned when removing an SSH key from a specific server.
type SSHKeyUnassignResponse struct {
	StatusCode int    `json:"statusCode"` // HTTP status code returned by the API
	Message    string `json:"message"`    // Human-readable message about the response
	Data       bool   `json:"data"`       // Success indicator (true if unassignment was successful)
}
