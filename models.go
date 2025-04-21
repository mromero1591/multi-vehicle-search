package main

type VehicleRequest struct {
	Length   int `json:"length"`
	Quantity int `json:"quantity"`
}

type Listing struct {
	ID             string `json:"id"`
	LocationID     string `json:"location_id"`
	Length         int    `json:"length"`
	Width          int    `json:"width"`
	PriceInCents   int    `json:"price_in_cents"`
}

type Result struct {
	LocationID        string   `json:"location_id"`
	ListingIDs        []string `json:"listing_ids"`
	TotalPriceInCents int      `json:"total_price_in_cents"`
}
