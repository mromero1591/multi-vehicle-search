package main

import (
	"encoding/json"
	"log"
	"net/http"
)

var listings []Listing

func handler(w http.ResponseWriter, r *http.Request) {
	var request []VehicleRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	vehicles := flattenVehicles(request)
	results := findValidCombinations(vehicles, listings)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func main() {
	var err error
	listings, err = loadListings("./data/listings.json")
	if err != nil {
		log.Fatalf("failed to load listings: %v", err)
	}

	http.HandleFunc("/", handler)
	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
