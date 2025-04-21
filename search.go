package main

import (
	"encoding/json"
	"os"
	"sort"
)

func loadListings(path string) ([]Listing, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var listings []Listing
	if err := json.Unmarshal(data, &listings); err != nil {
		return nil, err
	}
	return listings, nil
}

func flattenVehicles(request []VehicleRequest) []int {
	var vehicles []int
	for _, r := range request {
		for i := 0; i < r.Quantity; i++ {
			vehicles = append(vehicles, r.Length)
		}
	}
	return vehicles
}

func fits(listing Listing, length int) bool {
	return (length <= listing.Length && 10 <= listing.Width) ||
		(10 <= listing.Length && length <= listing.Width)
}

func findValidCombinations(vehicles []int, listings []Listing) []Result {
	byLocation := map[string][]Listing{}
	for _, l := range listings {
		byLocation[l.LocationID] = append(byLocation[l.LocationID], l)
	}

	var results []Result
	for locationID, locationListings := range byLocation {
		sort.Slice(locationListings, func(i, j int) bool {
			return locationListings[i].PriceInCents < locationListings[j].PriceInCents
		})

		used := make([]bool, len(locationListings))
		listingIDs := []string{}
		totalPrice := 0

		// Greedily assign vehicles to cheapest listings that can fit them
		for _, vehicle := range vehicles {
			assigned := false
			for i, listing := range locationListings {
				if !used[i] && fits(listing, vehicle) {
					used[i] = true
					listingIDs = append(listingIDs, listing.ID)
					totalPrice += listing.PriceInCents
					assigned = true
					break
				}
			}
			if !assigned {
				listingIDs = nil
				totalPrice = 0
				break
			}
		}

		if len(listingIDs) == len(vehicles) {
			results = append(results, Result{
				LocationID:        locationID,
				ListingIDs:        listingIDs,
				TotalPriceInCents: totalPrice,
			})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].TotalPriceInCents < results[j].TotalPriceInCents
	})
	return results
}
