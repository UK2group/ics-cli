package cmd

// ServerResponse represents the server list API response.
// This is returned when listing all servers.
type ServerResponse struct {
	StatusCode int      `json:"statusCode"` // HTTP status code returned by the API
	Message    string   `json:"message"`    // Human-readable message about the response
	Data       []Server `json:"data"`       // Array of server objects
}

// Server represents an individual server in the API response.
// Contains basic information about a server from the list endpoint.
type Server struct {
	ID             string `json:"id"`              // Unique identifier of the server
	Hostname       string `json:"hostname"`        // Server hostname
	MacAddress     string `json:"mac_address"`     // Primary MAC address of the server
	PublicIP       string `json:"public_ip"`       // Primary public IP address
	ServiceID      int    `json:"service_id"`      // Service ID associated with the server
	DatacenterName string `json:"datacenter_name"` // Name of the datacenter where the server is located
	FriendlyName   string `json:"friendly_name"`   // User-defined friendly name for the server
	ServerType     string `json:"server_type"`     // Type of server (e.g., "dedicated", "cloud")
}

// ServerDetailResponse represents the server detail API response.
// This is returned when fetching details for a specific server.
type ServerDetailResponse struct {
	StatusCode int          `json:"statusCode"` // HTTP status code returned by the API
	Message    string       `json:"message"`    // Human-readable message about the response
	Data       ServerDetail `json:"data"`       // Detailed server information
}

// ServerDetail represents the detailed information about a server.
// Contains comprehensive information about a specific server.
type ServerDetail struct {
	Hostname                string `json:"hostname"`                  // Server hostname
	MacAddress              string `json:"mac_address"`               // Primary MAC address of the server
	PublicIP                string `json:"public_ip"`                 // Primary public IP address
	ServiceID               int    `json:"service_id"`                // Service ID associated with the server
	FriendlyName            string `json:"friendly_name"`             // User-defined friendly name for the server
	ServerID                int    `json:"id"`                        // Unique identifier of the server
	OperatingSystemID       string `json:"operatingSystemId"`         // ID of the installed operating system
	OperatingSystemName     string `json:"operating_system"`          // Name of the installed operating system
	OperatingSystemUsername string `json:"operating_system_user"`     // Username for OS login
	OperatingSystemPassword string `json:"operating_system_password"` // Password for OS login
	DatacenterName          string `json:"datacenter"`                // Name of the datacenter where the server is located

	// IPAddresses contains all IP addresses assigned to the server
	IPAddresses []struct {
		IPAddress string `json:"ipAddress"` // The IP address
		IsPrimary bool   `json:"isPrimary"` // Whether this is the primary IP address
		Gateway   string `json:"gateway"`   // Gateway address for this IP
		Netmask   string `json:"netmask"`   // Network mask for this IP
		VlanID    string `json:"vlanId"`    // VLAN ID if the IP is on a VLAN
	} `json:"ip_addresses"`

	// NetworkPort contains network port information
	NetworkPort []struct {
		ID         int    `json:"id"`         // Port identifier
		PortNumber int    `json:"portNumber"` // Physical port number
		MacAddress string `json:"macAddress"` // MAC address for this port
		IPAddress  string `json:"ipAddress"`  // IP address assigned to this port
		Speed      string `json:"speed"`      // Current port speed
		MaxSpeed   string `json:"maxSpeed"`   // Maximum supported port speed
	} `json:"network_port"`

	// ProvisioningStatus contains information about the server's provisioning state
	ProvisioningStatus struct {
		IsProvisioning bool   `json:"isProvisioning"` // Whether the server is currently being provisioned
		StatusMessage  string `json:"statusMessage"`  // Human-readable status message
	} `json:"provisioning_status"`
}

// PowerStatusResponse represents the power status API response.
// This is returned when checking the power state of a server.
type PowerStatusResponse struct {
	StatusCode int    `json:"statusCode"` // HTTP status code returned by the API
	Message    string `json:"message"`    // Human-readable message about the response
	Data       struct {
		IsPoweredOn bool `json:"is_powered_on"` // Whether the server is powered on
	} `json:"data"`
}

// AssignedSSHKeysResponse represents the response from the server SSH keys API.
// This is returned when listing SSH keys assigned to a server.
type AssignedSSHKeysResponse struct {
	StatusCode int              `json:"statusCode"` // HTTP status code returned by the API
	Message    string           `json:"message"`    // Human-readable message about the response
	Data       []AssignedSSHKey `json:"data"`       // Array of SSH key objects
}

// AssignedSSHKey represents an individual SSH key in the API response.
// Contains information about an SSH key assigned to a server.
type AssignedSSHKey struct {
	ID        int    `json:"id"`         // Unique identifier of the SSH key
	Label     string `json:"label"`      // User-defined label for the SSH key
	Key       string `json:"key"`        // The actual SSH public key content
	CreatedAt int64  `json:"created_at"` // Unix timestamp when the key was created
	UpdatedAt int64  `json:"updated_at"` // Unix timestamp when the key was last updated
}

// RemoteAccessResponse represents the response from the remote access API.
// This is returned when requesting console access to a server.
type RemoteAccessResponse struct {
	StatusCode int    `json:"statusCode"` // HTTP status code returned by the API
	Message    string `json:"message"`    // Human-readable message about the response
	Data       struct {
		Redirect string `json:"redirect"` // URL to access the remote console
	} `json:"data"`
}

// CustomPXEResponse represents the response from the custom PXE API.
// This is returned when setting a custom PXE URL for a server.
type CustomPXEResponse struct {
	StatusCode int    `json:"statusCode"` // HTTP status code returned by the API
	Message    string `json:"message"`    // Human-readable message about the response
	Data       int    `json:"data"`       // Result data (typically a success indicator)
}

// GenericServerResponse represents the response from several management APIs.
// This is returned when performing power operations like power on/off/reboot.
type GenericServerResponse struct {
	StatusCode int    `json:"statusCode"` // HTTP status code returned by the API
	Message    string `json:"message"`    // Human-readable message about the response
	Data       struct {
		Success bool `json:"success"` // Whether the operation was successful
	} `json:"data"`
}

// OSListResponse represents the response from the available OS list API.
// This is returned when fetching available operating systems for a server.
type OSListResponse struct {
	StatusCode int    `json:"statusCode"` // HTTP status code returned by the API
	Message    string `json:"message"`    // Human-readable message about the response
	Data       struct {
		OSList []OS `json:"osList"` // List of available operating systems
	} `json:"data"`
}

// OS represents an individual operating system in the API response.
// Contains details about an operating system available for installation.
type OS struct {
	ID       string   `json:"id"`       // Unique identifier of the operating system
	Name     string   `json:"name"`     // Display name of the operating system
	Version  string   `json:"version"`  // Version of the operating system
	Licenses []string `json:"licenses"` // Licenses required for this operating system
}

// InventoryResponse represents the response from the inventory API
type InventoryResponse struct {
	StatusCode int                `json:"statusCode"`
	Message    string             `json:"message"`
	Data       []InventoryDetails `json:"data"`
}

// InventoryDetails represents a server from inventory
type InventoryDetails struct {
	SkuID            int        `json:"sku_id"`
	Quantity         int        `json:"quantity"`
	AutoProvisionQty int        `json:"auto_provision_quantity"`
	DatacenterID     int        `json:"datacenter_id"`
	RegionID         int        `json:"region_id"`
	LocationCode     string     `json:"location_code"`
	CPUBrand         string     `json:"cpu_brand"`
	CPUModel         string     `json:"cpu_model"`
	CPUClockSpeedGhz float64    `json:"cpu_clock_speed_ghz"`
	CPUCores         int        `json:"cpu_cores"`
	CPUCount         int        `json:"cpu_count"`
	TotalSSDSizeGB   int        `json:"total_ssd_size_gb"`
	TotalHDDSizeGB   int        `json:"total_hdd_size_gb"`
	TotalNVMESizeGB  int        `json:"total_nvme_size_gb"`
	RAIDEnabled      bool       `json:"raid_enabled"`
	TotalRAMGB       int        `json:"total_ram_gb"`
	NICSpeedMbps     int        `json:"nic_speed_mbps"`
	QTProductID      int        `json:"qt_product_id"`
	Status           string     `json:"status"`
	Metadata         []Metadata `json:"metadata"`
	CurrencyCode     string     `json:"currency_code"`
	SkuProductName   string     `json:"sku_product_name"`
	Price            string     `json:"price"`
}

// Metadata represents additional server information
type Metadata struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Value       string `json:"value"`
}

// GroupedInventory represents inventory grouped by location and SKU
type GroupedInventory struct {
	LocationCode   string
	SkuProductName string
	Price          string
	TotalQuantity  int
}

// AddonsResponse represents the response from the add-ons API
type AddonsResponse struct {
	StatusCode int        `json:"statusCode"`
	Message    string     `json:"message"`
	Data       AddonTypes `json:"data"` // Changed from array to object
}

// AddonTypes contains the different types of add-ons
type AddonTypes struct {
	OperatingSystems OperatingSystemsSection `json:"operating_systems"`
	Licenses         LicenseSection          `json:"licenses"`
	SupportLevels    SupportSection          `json:"support_levels"`
}

// OperatingSystemsSection contains operating system options
type OperatingSystemsSection struct {
	Name     string      `json:"name"`
	Required string      `json:"required"`
	Products []OSProduct `json:"products"`
}

// OSProduct represents an operating system product
type OSProduct struct {
	Name         string      `json:"name"`
	OSType       string      `json:"os_type"`
	ProductCode  string      `json:"product_code"`
	Price        float64     `json:"price"`
	PricePerCore interface{} `json:"price_per_core"`
}

// LicenseSection contains license options
type LicenseSection struct {
	Name     string           `json:"name"`
	Products []LicenseProduct `json:"products"`
}

// LicenseProduct represents a license product
type LicenseProduct struct {
	Name        string  `json:"name"`
	ProductCode string  `json:"product_code"`
	Price       float64 `json:"price"`
}

// SupportSection contains support options
type SupportSection struct {
	Name     string           `json:"name"`
	Products []SupportProduct `json:"products"`
}

// SupportProduct represents a support product
type SupportProduct struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	ProductCode string  `json:"product_code"`
	Price       float64 `json:"price"`
}

// OrderRequest represents the request body for placing an order
type OrderRequest struct {
	SKUProductName             string `json:"sku_product_name"`
	Quantity                   int    `json:"quantity"`
	LocationCode               string `json:"location_code"`
	OperatingSystemProductCode string `json:"operating_system_product_code"`
	LicenseProductCode         string `json:"license_product_code,omitempty"`
	AdditionalBandwidthTB      int    `json:"additional_bandwidth_tb,omitempty"`
	SupportLevelProductCode    string `json:"support_level_product_code,omitempty"`
	SSHKeyIDs                  []int  `json:"ssh_key_ids,omitempty"`
}

// OrderResponse represents the response from the order API
// This is returned when successfully placing a server order.
type OrderResponse struct {
	StatusCode int    `json:"statusCode"` // HTTP status code returned by the API
	Message    string `json:"message"`    // Human-readable message about the response
	Data       struct {
		OrderServiceIDs []int `json:"order_service_ids"` // Array of service IDs created by the order
	} `json:"data"`
}
