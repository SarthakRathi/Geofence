package models

import "time"

type Coordinate struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type Zone struct {
	ID       string       `json:"id"`
	Name     string       `json:"name"`
	Boundary []Coordinate `json:"boundary"`
}

type LocationEvent struct {
	VehicleID string  `json:"vehicle_id"`
	Lat       float64 `json:"lat"`
	Lon       float64 `json:"lon"`
	Timestamp int64   `json:"timestamp"`
}

type VehicleStatus struct {
	VehicleID     string    `json:"vehicle_id"`
	CurrentZoneID string    `json:"current_zone_id"`
	LastSeen      time.Time `json:"last_seen"`
}
