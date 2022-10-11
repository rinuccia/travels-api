package model

// Location represent location data model
type Location struct {
	LocationId uint32 `json:"location_id" validate:"required"`
	Place      string `json:"place" validate:"required"`
	Country    string `json:"country" validate:"required,min=2,max=50"`
}

type Locs []Location

// Locations represents a list of all locations
type Locations struct {
	List []Location `json:"list"`
}

// AvgRating represent average location rating data model
type AvgRating struct {
	Avg float32 `json:"avg"`
}
