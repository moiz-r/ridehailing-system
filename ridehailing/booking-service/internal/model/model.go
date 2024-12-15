package model

import "time"

type Booking struct {
	ID     int       `json:"id"`      // Booking ID
	UserID int       `json:"user_id"` // User ID
	RideID int       `json:"ride_id"` // Ride ID
	Time   time.Time `json:"time"`    // Timestamp of the booking
}

type Ride struct {
	Source      string  `json:"source"`      // Source location
	Destination string  `json:"destination"` // Destination location
	Distance    float64 `json:"distance"`    // Distance in kilometers with decimal precision
	Cost        float64 `json:"cost"`        // Cost of the ride with decimal precision
}

type BookingDTO struct {
	BookingID   int       `json:"booking_id"`  // Booking ID
	UserID      int       `json:"user_id"`     // User ID
	RideID      int       `json:"ride_id"`     // Ride ID
	Source      string    `json:"source"`      // Source location
	Destination string    `json:"destination"` // Destination location
	Distance    float64   `json:"distance"`    // Distance in kilometers with decimal precision
	Cost        float64   `json:"cost"`        // Cost of the ride with decimal precision
	Time        time.Time `json:"time"`        // Timestamp of the booking
}
