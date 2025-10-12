package cmd

// VPSResponse represents the VPS list API response.
// This is returned when listing all VPS servers.
type VPSResponse struct {
	StatusCode int        `json:"statusCode"` // HTTP status code returned by the API
	Message    string     `json:"message"`    // Human-readable message about the response
	Data       []VPSBasic `json:"data"`       // Array of VPS objects
}

// VPSBasic represents an individual VPS in the API response.
// Contains basic information about a VPS from the list endpoint.
type VPSBasic struct {
	ID            int    `json:"id"`         // Unique identifier of the VPS
	ServiceID     int    `json:"service_id"` // Service ID associated with the VPS
	Label         string `json:"label"`      // User-defined label for the VPS
	Hostname      string `json:"hostname"`   // VPS hostname
	IPAddress     string `json:"ip_address"` // Primary IP address
	Backups       int    `json:"backups"`    // Whether backups are enabled
	HourlyBilling int    `json:"hourly_billing"` // Whether hourly billing is enabled (0 or 1)
	PlanID        int    `json:"plan_id"`    // Plan ID
	Location      struct {
		ID   int    `json:"id"`   // Location ID
		Name string `json:"name"` // Location name
	} `json:"location"` // Datacenter location information
}

// VPSDetailResponse represents the VPS detail API response.
// This is returned when fetching details for a specific VPS.
type VPSDetailResponse struct {
	StatusCode int       `json:"statusCode"` // HTTP status code returned by the API
	Message    string    `json:"message"`    // Human-readable message about the response
	Data       VPSDetail `json:"data"`       // Detailed VPS information
}

// VPSDetail represents the detailed information about a VPS.
// Contains comprehensive information about a specific VPS.
type VPSDetail struct {
	ID                      int    `json:"id"`                         // Unique identifier of the VPS
	Hostname                string `json:"hostname"`                   // VPS hostname
	Memory                  int    `json:"memory"`                     // RAM in MB
	CPUs                    int    `json:"cpus"`                       // Number of CPU cores
	CreatedAt               string `json:"created_at"`                 // Creation timestamp
	Built                   bool   `json:"built"`                      // Whether VPS is built
	Booted                  bool   `json:"booted"`                     // Whether VPS is booted
	Label                   string `json:"label"`                      // User-defined label
	OperatingSystem         string `json:"operating_system"`           // Operating system type
	OperatingSystemDistro   string `json:"operating_system_distro"`    // OS distribution
	Note                    string `json:"note"`                       // User notes
	Suspended               bool   `json:"suspended"`                  // Whether VPS is suspended
	Domain                  string `json:"domain"`                     // Domain name
	TotalDiskSize           int    `json:"total_disk_size"`            // Total disk size in GB
	InitialRootPassword     string `json:"initial_root_password"`      // Initial root password
	HourlyBilling           int    `json:"hourly_billing"`             // Whether hourly billing is enabled (0 or 1)
	Backups                 int    `json:"backups"`                    // Whether backups are enabled (0 or 1)

	// IPAddresses contains all IP addresses assigned to the VPS
	IPAddresses []struct {
		IPAddress struct {
			Address        string `json:"address"`         // The IP address
			Broadcast      string `json:"broadcast"`       // Broadcast address
			NetworkAddress string `json:"network_address"` // Network address
			Gateway        string `json:"gateway"`         // Gateway address
			CreatedAt      string `json:"created_at"`      // Creation timestamp
			RdnsEntry      string `json:"rdnsEntry"`       // Reverse DNS entry
		} `json:"ip_address"`
	} `json:"ip_addresses"`
}

// VPSPowerStatusResponse represents the power status API response.
// This is returned when checking the power state of a VPS.
type VPSPowerStatusResponse struct {
	StatusCode int    `json:"statusCode"` // HTTP status code returned by the API
	Message    string `json:"message"`    // Human-readable message about the response
	Data       struct {
		IsPoweredOn bool `json:"is_powered_on"` // Whether the VPS is powered on
	} `json:"data"`
}

// VPSGenericResponse represents the response from several management APIs.
// This is returned when performing power operations like power on/off/reboot.
type VPSGenericResponse struct {
	StatusCode int    `json:"statusCode"` // HTTP status code returned by the API
	Message    string `json:"message"`    // Human-readable message about the response
	Data       struct {
		Success bool `json:"success"` // Whether the operation was successful
	} `json:"data"`
}

// VPSEmptyDataResponse represents responses that return an empty data array.
// This is used for operations like reboot that return "data": []
type VPSEmptyDataResponse struct {
	StatusCode int           `json:"status_code"` // HTTP status code returned by the API
	Message    string        `json:"message"`     // Human-readable message about the response
	Data       []interface{} `json:"data"`        // Empty array
}

// VPSPlanResponse represents the response from the VPS plans API
type VPSPlanResponse struct {
	StatusCode int       `json:"statusCode"`
	Message    string    `json:"message"`
	Data       []VPSPlan `json:"data"`
}

// VPSPlan represents a VPS plan/package
type VPSPlan struct {
	ID          int    `json:"id"`           // Plan ID
	Label       string `json:"label"`        // Plan description/label
	Price       string `json:"price"`        // Monthly price
	ProductName string `json:"product_name"` // Plan product name/SKU
}

// VPSOSResponse represents the response from the VPS OS list API
type VPSOSResponse struct {
	StatusCode int     `json:"statusCode"`
	Message    string  `json:"message"`
	Data       []VPSOS `json:"data"`
}

// VPSOS represents an operating system available for VPS
type VPSOS struct {
	Name                 string  `json:"name"`                   // Display name of the operating system
	ComponentCode        string  `json:"component_code"`         // OS component code
	Price                float64 `json:"price"`                  // Monthly price
	PriceHourly          float64 `json:"price_hourly"`           // Hourly price
	Currency             string  `json:"currency"`               // Currency code
	MinDiskSize          string  `json:"min_disk_size"`          // Minimum disk size in GB
	OperatingSystemArch  string  `json:"operating_system_arch"`  // Architecture (x64, etc)
	MinMemorySize        string  `json:"min_memory_size"`        // Minimum memory in MB
}

// VPSLocationResponse represents the response from the VPS locations API
type VPSLocationResponse struct {
	StatusCode int           `json:"statusCode"`
	Message    string        `json:"message"`
	Data       []VPSLocation `json:"data"`
}

// VPSLocation represents a VPS datacenter location
type VPSLocation struct {
	ID   int    `json:"id"`   // Location ID
	Name string `json:"name"` // Location name
}

// VPSDeployRequest represents the request body for deploying a VPS
type VPSDeployRequest struct {
	ProductName     string `json:"product_name"`     // VPS plan/SKU (e.g., "V1", "V2")
	Label           string `json:"label"`            // Friendly name/label for the VPS
	Hostname        string `json:"hostname"`         // Hostname for the VPS
	Backups         bool   `json:"backups"`          // Enable backups
	BillHourly      bool   `json:"bill_hourly"`      // Enable hourly billing
	OSComponentCode string `json:"os_component_code"` // Operating system component code
}

// VPSDeployResponse represents the response from the VPS deploy API
type VPSDeployResponse struct {
	StatusCode int           `json:"status_code"` // HTTP status code returned by the API
	Message    string        `json:"message"`     // Human-readable message about the response
	Data       []interface{} `json:"data"`        // Empty array on success
}
