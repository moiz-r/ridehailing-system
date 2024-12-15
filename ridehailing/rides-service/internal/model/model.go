package model

type Ride struct {
	ID          int     `json:"id"`          // Ride ID
	Source      string  `json:"source"`      // Source location
	Destination string  `json:"destination"` // Destination location
	Distance    float64 `json:"distance"`    // Distance in kilometers with decimal precision
	Cost        float64 `json:"cost"`        // Cost of the ride with decimal precision
}
