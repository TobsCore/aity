package main

type Progress struct {
	ID       string   `json:"id,omitempty"`
	TrackID  string   `json:"trackid"`
	Distance Distance `json:"distance"`
	Date     string   `json:"date"`
}

type Distance struct {
	Value int    `json:"distance"`
	Unit  string `json:"unit"`
}
