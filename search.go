package main

import (
	"encoding/json"
	"os"
	"sort"
)


const BLOCK_SIZE = 10

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

func groupByLocations(listings []Listing) map[string][]Listing {
	byLocation := map[string][]Listing{}
	for _, l := range listings {
		byLocation[l.LocationID] = append(byLocation[l.LocationID], l)
	}
	return byLocation
}

func canPlace(row []bool, blocks int) int {
    maxStart := len(row) - blocks
    for start := 0; start <= maxStart; start++ {
        ok := true
        for offset := 0; offset < blocks; offset++ {
            if row[start+offset] {
                ok = false
                break
            }
        }
        if ok {
            return start
        }
    }
    return -1
}

func place(row []bool, start, blocks int) {
    for k := 0; k < blocks; k++ {
        row[start+k] = true
    }
}

func packOrientation(length, width int, vehicles []int) ([]int, int) {
    rows, cols := width/BLOCK_SIZE, length/BLOCK_SIZE
    space := make([][]bool, rows)
    for i := range space {
        space[i] = make([]bool, cols)
    }

    var remaining []int
    packed := 0

    for _, vehicle := range vehicles {
        blocks := vehicle / BLOCK_SIZE
        placed := false

        // scan each row looking for a fit
        for r := 0; r < rows && !placed; r++ {
            if start := canPlace(space[r], blocks); start >= 0 {
                place(space[r], start, blocks)
                packed++
                placed = true
            }
        }

        if !placed {
            remaining = append(remaining, vehicle)
        }
    }

    return remaining, packed
}

func packVehicles(listing Listing, vehicles []int) ([]int, int) {
    orientations := [][2]int{
        {listing.Length, listing.Width},
        {listing.Width, listing.Length},
    }

    bestRemaning := vehicles
    bestCount := 0

    for _, o := range orientations {
        rem, cnt := packOrientation(o[0], o[1], vehicles)
        if cnt > bestCount {
            bestCount = cnt
            bestRemaning = rem
        }
    }
    return bestRemaning, bestCount
}

func findValidCombinations(vehicles []int, listings []Listing) []Result {
	locations := groupByLocations(listings)
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
