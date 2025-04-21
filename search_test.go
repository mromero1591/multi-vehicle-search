package main

import (
	"os"
	"reflect"
	"sort"
	"testing"
)

func TestFlattenVehicles(t *testing.T) {
	tests := []struct {
		name     string
		requests []VehicleRequest
		want     []int
	}{
		{
			name: "Single vehicle",
			requests: []VehicleRequest{
				{Length: 20, Quantity: 1},
			},
			want: []int{20},
		},
		{
			name: "Multiple vehicles of same size",
			requests: []VehicleRequest{
				{Length: 10, Quantity: 3},
			},
			want: []int{10, 10, 10},
		},
		{
			name: "Multiple vehicles of different sizes",
			requests: []VehicleRequest{
				{Length: 10, Quantity: 1},
				{Length: 20, Quantity: 2},
				{Length: 30, Quantity: 1},
			},
			want: []int{30, 20, 20, 10},
		},
		{
			name:     "Empty request",
			requests: []VehicleRequest{},
			want:     []int{},
		},
		{
			name: "Zero quantity",
			requests: []VehicleRequest{
				{Length: 10, Quantity: 0},
			},
			want: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := flattenVehicles(tt.requests)
			
			// Since we know flattenVehicles should sort in descending order
			sort.Slice(tt.want, func(i, j int) bool {
				return tt.want[i] > tt.want[j]
			})
			
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("flattenVehicles() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFits(t *testing.T) {
	tests := []struct {
		name    string
		listing Listing
		length  int
		want    bool
	}{
		{
			name:    "Fits in length",
			listing: Listing{Length: 20, Width: 10},
			length:  15,
			want:    true,
		},
		{
			name:    "Fits in width",
			listing: Listing{Length: 10, Width: 20},
			length:  15,
			want:    true,
		},
		{
			name:    "Fits exactly in length",
			listing: Listing{Length: 20, Width: 10},
			length:  20,
			want:    true,
		},
		{
			name:    "Too long for both dimensions",
			listing: Listing{Length: 15, Width: 15},
			length:  20,
			want:    false,
		},
		{
			name:    "Width too narrow",
			listing: Listing{Length: 20, Width: 5},
			length:  15,
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fits(tt.listing, tt.length); got != tt.want {
				t.Errorf("fits() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFindValidCombinations(t *testing.T) {
	testListings := []Listing{
		{ID: "l1", LocationID: "loc1", Length: 20, Width: 10, PriceInCents: 100},
		{ID: "l2", LocationID: "loc1", Length: 30, Width: 15, PriceInCents: 200},
		{ID: "l3", LocationID: "loc1", Length: 10, Width: 10, PriceInCents: 50},
		{ID: "l4", LocationID: "loc2", Length: 40, Width: 20, PriceInCents: 300},
		{ID: "l5", LocationID: "loc2", Length: 10, Width: 10, PriceInCents: 75},
		{ID: "l6", LocationID: "loc3", Length: 10, Width: 10, PriceInCents: 30},
		{ID: "l7", LocationID: "loc3", Length: 5, Width: 15, PriceInCents: 20},
	}

	tests := []struct {
		name     string
		vehicles []int
		want     []Result
	}{
		{
			name:     "Single small vehicle",
			vehicles: []int{10},
			want: []Result{
				{LocationID: "loc3", ListingIDs: []string{"l6"}, TotalPriceInCents: 30},
				{LocationID: "loc1", ListingIDs: []string{"l3"}, TotalPriceInCents: 50},
				{LocationID: "loc2", ListingIDs: []string{"l5"}, TotalPriceInCents: 75},
			},
		},
		{
			name:     "Multiple vehicles with exact fit",
			vehicles: []int{30, 10},
			want: []Result{
				{LocationID: "loc1", ListingIDs: []string{"l2", "l3"}, TotalPriceInCents: 250},
				{LocationID: "loc2", ListingIDs: []string{"l4", "l5"}, TotalPriceInCents: 375},
			},
		},
		{
			name:     "Vehicle won't fit",
			vehicles: []int{50},
			want:     []Result{},
		},
		{
			name:     "Not enough listings at location",
			vehicles: []int{20, 20, 20},
			want:     []Result{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := findValidCombinations(tt.vehicles, testListings)

			// Sort results by total price for stable comparison
			sort.Slice(got, func(i, j int) bool {
				return got[i].TotalPriceInCents < got[j].TotalPriceInCents
			})

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findValidCombinations() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoadListings(t *testing.T) {
	// Create a temporary test JSON file
	testFilePath := "test_listings.json"
	testData := `[
		{"id": "test1", "location_id": "loc1", "length": 20, "width": 10, "price_in_cents": 100},
		{"id": "test2", "location_id": "loc2", "length": 30, "width": 15, "price_in_cents": 200}
	]`

	// Write the test data to a file
	if err := os.WriteFile(testFilePath, []byte(testData), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer os.Remove(testFilePath) // Clean up after the test

	// Run the test
	listings, err := loadListings(testFilePath)
	if err != nil {
		t.Errorf("loadListings() error = %v", err)
		return
	}

	// Check the result
	expected := []Listing{
		{ID: "test1", LocationID: "loc1", Length: 20, Width: 10, PriceInCents: 100},
		{ID: "test2", LocationID: "loc2", Length: 30, Width: 15, PriceInCents: 200},
	}

	if !reflect.DeepEqual(listings, expected) {
		t.Errorf("loadListings() = %v, want %v", listings, expected)
	}
}

