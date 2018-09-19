package model

type Track struct {
	Id       string   `json:"track-id"`
	Active   bool     `json:"active"`
	Start    Location `json:"start"`
	End      Location `json:"end"`
	Distance Distance `json:"totalDistance"`
}

type Location struct {
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
