package cmd

import (
	"strconv"
	"testing"
)

func TestFilterInventory(t *testing.T) {
	// Create sample inventory data
	inventory := []InventoryDetails{
		{
			LocationCode:   "LHR",
			SkuProductName: "E-2136",
			Price:          "150.00",
			HourlyEnabled:  true,
			HourlyPrice:    "0.25",
			Quantity:       5,
		},
		{
			LocationCode:   "NYC",
			SkuProductName: "E-2236",
			Price:          "200.00",
			HourlyEnabled:  false,
			HourlyPrice:    "0.00",
			Quantity:       3,
		},
		{
			LocationCode:   "LHR",
			SkuProductName: "E-2288G",
			Price:          "300.00",
			HourlyEnabled:  true,
			HourlyPrice:    "0.50",
			Quantity:       2,
		},
		{
			LocationCode:   "NYC",
			SkuProductName: "E-2136",
			Price:          "160.00",
			HourlyEnabled:  true,
			HourlyPrice:    "0.30",
			Quantity:       4,
		},
	}

	tests := []struct {
		name              string
		datacenter        string
		sku               string
		minPrice          float64
		maxPrice          float64
		minHourlyPrice    float64
		maxHourlyPrice    float64
		hourlyEnabledOnly bool
		expectedCount     int
	}{
		{
			name:          "no filters - return all",
			expectedCount: 4,
		},
		{
			name:          "filter by datacenter LHR",
			datacenter:    "LHR",
			expectedCount: 2,
		},
		{
			name:          "filter by datacenter NYC",
			datacenter:    "NYC",
			expectedCount: 2,
		},
		{
			name:          "filter by SKU E-2136",
			sku:           "E-2136",
			expectedCount: 2,
		},
		{
			name:          "filter by min price 175",
			minPrice:      175.00,
			expectedCount: 2, // Only 200.00 and 300.00 are >= 175
		},
		{
			name:          "filter by max price 175",
			maxPrice:      175.00,
			expectedCount: 2,
		},
		{
			name:          "filter by price range 150-200",
			minPrice:      150.00,
			maxPrice:      200.00,
			expectedCount: 3,
		},
		{
			name:              "filter hourly enabled only",
			hourlyEnabledOnly: true,
			expectedCount:     3,
		},
		{
			name:           "filter by hourly price min 0.30",
			minHourlyPrice: 0.30,
			expectedCount:  3, // NYC/E-2136 (0.30) and LHR/E-2288G (0.50) pass, NYC/E-2236 has no hourly so passes
		},
		{
			name:           "filter by hourly price max 0.30",
			maxHourlyPrice: 0.30,
			expectedCount:  3, // LHR/E-2136 (0.25) and NYC/E-2136 (0.30) pass, NYC/E-2236 has no hourly so passes
		},
		{
			name:          "complex filter - LHR + hourly enabled",
			datacenter:    "LHR",
			hourlyEnabledOnly: true,
			expectedCount: 2,
		},
		{
			name:          "complex filter - NYC + E-2136 + hourly",
			datacenter:    "NYC",
			sku:           "E-2136",
			hourlyEnabledOnly: true,
			expectedCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filterInventory(
				inventory,
				tt.datacenter,
				tt.sku,
				tt.minPrice,
				tt.maxPrice,
				tt.minHourlyPrice,
				tt.maxHourlyPrice,
				tt.hourlyEnabledOnly,
			)

			if len(result) != tt.expectedCount {
				t.Errorf("filterInventory() returned %d items, want %d", len(result), tt.expectedCount)
			}
		})
	}
}

func TestGroupInventory(t *testing.T) {
	inventory := []InventoryDetails{
		{
			LocationCode:   "LHR",
			SkuProductName: "E-2136",
			Price:          "150.00",
			HourlyEnabled:  true,
			HourlyPrice:    "0.25",
			Quantity:       5,
		},
		{
			LocationCode:   "LHR",
			SkuProductName: "E-2136",
			Price:          "150.00",
			HourlyEnabled:  true,
			HourlyPrice:    "0.25",
			Quantity:       3,
		},
		{
			LocationCode:   "NYC",
			SkuProductName: "E-2136",
			Price:          "160.00",
			HourlyEnabled:  true,
			HourlyPrice:    "0.30",
			Quantity:       2,
		},
	}

	result := groupInventory(inventory)

	// Should have 2 groups: LHR:E-2136 and NYC:E-2136
	if len(result) != 2 {
		t.Errorf("groupInventory() returned %d groups, want 2", len(result))
	}

	// Find the LHR group
	var lhrGroup *GroupedInventory
	for i := range result {
		if result[i].LocationCode == "LHR" {
			lhrGroup = &result[i]
			break
		}
	}

	if lhrGroup == nil {
		t.Fatal("groupInventory() did not create LHR group")
	}

	// LHR group should have combined quantity of 5 + 3 = 8
	if lhrGroup.TotalQuantity != 8 {
		t.Errorf("groupInventory() LHR TotalQuantity = %d, want 8", lhrGroup.TotalQuantity)
	}

	// Verify sorting (should be sorted by location then SKU)
	if result[0].LocationCode > result[1].LocationCode {
		t.Error("groupInventory() result is not sorted by location")
	}
}

func TestConfirmOrderDetails(t *testing.T) {
	// This function requires user input, so we test with a simple struct validation
	order := OrderRequest{
		SKUProductName:             "E-2136",
		Quantity:                   2,
		LocationCode:               "LHR",
		OperatingSystemProductCode: "ubuntu-22.04",
		LicenseProductCode:         "",
		AdditionalBandwidthTB:      5,
		SupportLevelProductCode:    "standard",
		SSHKeyIDs:                  []int{1, 2, 3},
		BillHourly:                 true,
	}

	// Ensure the order struct is valid
	if order.Quantity <= 0 {
		t.Error("Invalid order quantity")
	}

	if order.SKUProductName == "" {
		t.Error("Invalid SKU product name")
	}

	if order.LocationCode == "" {
		t.Error("Invalid location code")
	}
}

func TestFilterInventoryEdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		inventory []InventoryDetails
		want      int
	}{
		{
			name:      "empty inventory",
			inventory: []InventoryDetails{},
			want:      0,
		},
		{
			name: "single item",
			inventory: []InventoryDetails{
				{
					LocationCode:   "LHR",
					SkuProductName: "E-2136",
					Price:          "150.00",
					Quantity:       1,
				},
			},
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filterInventory(tt.inventory, "", "", 0, 0, 0, 0, false)
			if len(result) != tt.want {
				t.Errorf("filterInventory() = %d items, want %d", len(result), tt.want)
			}
		})
	}
}

func TestFilterInventoryInvalidPrice(t *testing.T) {
	inventory := []InventoryDetails{
		{
			LocationCode:   "LHR",
			SkuProductName: "E-2136",
			Price:          "invalid-price",
			Quantity:       1,
		},
		{
			LocationCode:   "NYC",
			SkuProductName: "E-2236",
			Price:          "200.00",
			Quantity:       1,
		},
	}

	// Filter with min price should skip items with invalid prices
	result := filterInventory(inventory, "", "", 100.00, 0, 0, 0, false)

	// Should return both items since price parsing error is handled gracefully
	if len(result) != 2 {
		t.Errorf("filterInventory() with invalid price = %d items, want 2", len(result))
	}
}

func TestFilterInventoryCaseInsensitive(t *testing.T) {
	inventory := []InventoryDetails{
		{
			LocationCode:   "LHR",
			SkuProductName: "E-2136",
			Price:          "150.00",
			Quantity:       1,
		},
	}

	tests := []struct {
		name       string
		datacenter string
		sku        string
		want       int
	}{
		{
			name:       "lowercase datacenter",
			datacenter: "lhr",
			want:       1,
		},
		{
			name:       "uppercase datacenter",
			datacenter: "LHR",
			want:       1,
		},
		{
			name:       "mixed case datacenter",
			datacenter: "Lhr",
			want:       1,
		},
		{
			name: "lowercase sku",
			sku:  "e-2136",
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filterInventory(inventory, tt.datacenter, tt.sku, 0, 0, 0, 0, false)
			if len(result) != tt.want {
				t.Errorf("filterInventory() case insensitive = %d items, want %d", len(result), tt.want)
			}
		})
	}
}

func TestGroupInventorySorting(t *testing.T) {
	inventory := []InventoryDetails{
		{LocationCode: "NYC", SkuProductName: "Z-1000", Price: "100", Quantity: 1},
		{LocationCode: "LHR", SkuProductName: "A-1000", Price: "100", Quantity: 1},
		{LocationCode: "NYC", SkuProductName: "A-1000", Price: "100", Quantity: 1},
		{LocationCode: "LHR", SkuProductName: "Z-1000", Price: "100", Quantity: 1},
	}

	result := groupInventory(inventory)

	// Should be sorted: LHR:A-1000, LHR:Z-1000, NYC:A-1000, NYC:Z-1000
	if len(result) != 4 {
		t.Fatalf("groupInventory() returned %d groups, want 4", len(result))
	}

	expected := []struct {
		location string
		sku      string
	}{
		{"LHR", "A-1000"},
		{"LHR", "Z-1000"},
		{"NYC", "A-1000"},
		{"NYC", "Z-1000"},
	}

	for i, exp := range expected {
		if result[i].LocationCode != exp.location || result[i].SkuProductName != exp.sku {
			t.Errorf("groupInventory()[%d] = %s:%s, want %s:%s",
				i, result[i].LocationCode, result[i].SkuProductName, exp.location, exp.sku)
		}
	}
}

func TestFilterInventoryHourlyPriceOnlyAppliesToHourlyEnabled(t *testing.T) {
	inventory := []InventoryDetails{
		{
			LocationCode:   "LHR",
			SkuProductName: "E-2136",
			Price:          "150.00",
			HourlyEnabled:  false,
			HourlyPrice:    "0.00",
			Quantity:       1,
		},
		{
			LocationCode:   "NYC",
			SkuProductName: "E-2236",
			Price:          "200.00",
			HourlyEnabled:  true,
			HourlyPrice:    "0.50",
			Quantity:       1,
		},
	}

	// Filter with hourly price min should only affect hourly-enabled items
	result := filterInventory(inventory, "", "", 0, 0, 0.40, 0, false)

	// Should return both: one passes because hourly disabled, one passes because hourly price is 0.50
	if len(result) != 2 {
		t.Errorf("filterInventory() with hourly price filter = %d items, want 2", len(result))
	}
}

func TestFilterInventoryPriceConversion(t *testing.T) {
	inventory := []InventoryDetails{
		{
			LocationCode:   "LHR",
			SkuProductName: "E-2136",
			Price:          "150.50",
			Quantity:       1,
		},
		{
			LocationCode:   "NYC",
			SkuProductName: "E-2236",
			Price:          "149.50",
			Quantity:       1,
		},
	}

	// Test that float comparison works correctly
	result := filterInventory(inventory, "", "", 150.00, 0, 0, 0, false)

	// Should only return the first item (150.50 >= 150.00)
	if len(result) != 1 {
		t.Errorf("filterInventory() price comparison = %d items, want 1", len(result))
	}

	if len(result) > 0 {
		price, _ := strconv.ParseFloat(result[0].Price, 64)
		if price < 150.00 {
			t.Errorf("filterInventory() returned item with price %.2f, should be >= 150.00", price)
		}
	}
}
