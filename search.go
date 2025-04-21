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

	listings := []Listing{}
	if err := json.Unmarshal(data, &listings); err != nil {
		return nil, err
	}

	return listings, nil
}

func flattenVehicles(request []VehicleRequest) []int {
	vehicles := []int{}
	for _, r := range request {
		for i := 0; i < r.Quantity; i++ {
			vehicles = append(vehicles, r.Length)
		}
	}

	sort.Slice(vehicles, func(i, j int) bool {
		return vehicles[i] > vehicles[j]
	})

	return vehicles
}

func mapLocations(listings []Listing) map[string][]Listing {
	byLocation := map[string][]Listing{}
	for _, l := range listings {
		byLocation[l.LocationID] = append(byLocation[l.LocationID], l)
	}
	return byLocation
}

func packVehicles(listing Listing, vehicles []int) ([]int, int) {
	l1, w1 := listing.Length, listing.Width
	l2, w2 := listing.Width, listing.Length

	orientations := []struct {
		length int
		width  int
	}{
		{l1, w1},
		{l2, w2},
	}

	maxUsed := 0
	bestRemaining := vehicles

	for _, o := range orientations {
		rows := o.width / 10
		cols := o.length / 10
		space := make([][]bool, rows)
		for i := range space {
			space[i] = make([]bool, cols)
		}

		tmpRemaining := []int{}
		count := 0

		for _, v := range vehicles {
			lBlocks := v / 10
			fit := false

			for row := 0; row < rows; row++ {
				for col := 0; col <= cols-lBlocks; col++ {
					canFit := true
					for k := 0; k < lBlocks; k++ {
						if space[row][col+k] {
							canFit = false
							break
						}
					}
					if canFit {
						for k := 0; k < lBlocks; k++ {
							space[row][col+k] = true
						}
						fit = true
						break
					}
				}
				if fit {
					break
				}
			}

			if fit {
				count++ 
			} else {
				tmpRemaining = append(tmpRemaining, v) 
			}
		}

		if count > maxUsed {
			maxUsed = count
			bestRemaining = tmpRemaining
		}
	}

	return bestRemaining, maxUsed
}

func findValidCombinations(vehicles []int, listings []Listing) []Result {
	locations := mapLocations(listings)
	results := []Result{}

	for locationID, locationListings := range locations {
		sort.Slice(locationListings, func(i, j int) bool {
			return locationListings[i].PriceInCents < locationListings[j].PriceInCents
		})

		remaining := append([]int(nil), vehicles...) 
		totalPrice := 0
		listingIDs := []string{}

		for _, listing := range locationListings {
			if len(remaining) == 0 {
				break // all vehicles packed
			}

			updatedRemaining, packed := packVehicles(listing, remaining)
			if packed > 0 {
				listingIDs = append(listingIDs, listing.ID)
				totalPrice += listing.PriceInCents
				remaining = updatedRemaining
			}
		}

		if len(remaining) == 0 {
			results = append(results, Result{
				LocationID:         locationID,
				ListingIDs:         listingIDs,
				TotalPriceInCents:  totalPrice,
			})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].TotalPriceInCents < results[j].TotalPriceInCents
	})

	return results
}
