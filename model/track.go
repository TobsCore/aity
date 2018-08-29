package model

type Track struct {
	Start    Location `json:"start"`
	End      Location `json:"end"`
	Distance Distance `json:"totalDistance"`
}

type Location struct {
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
