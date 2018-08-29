package model

type Progress struct {
	Date     string   `json:"date"`
	Distance Distance `json:"distance"`
}

type Distance int64
