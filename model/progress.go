package model

import (
	"errors"
	"fmt"
	"time"
)

type Progress struct {
	Date     string   `json:"date"`
	Distance Distance `json:"distance"`
}

const TimeFormat = "2006-01-02"

func (p *Progress) GetDate() (time.Time, error) {
	date, err := time.Parse(TimeFormat, p.Date)
	if err != nil {
		return time.Time{}, errors.New(fmt.Sprintf("Cannot parse date %s", p.Date))
	}
	return date, nil
}

type Distance int64

func AccProgresses(progs []Progress) []Progress {
	tempProgressInfo := make(map[string]Distance, len(progs))
	var resProgresses []Progress
	for _, prog := range progs {
		tmpDistane := tempProgressInfo[prog.Date] + prog.Distance
		tempProgressInfo[prog.Date] = tmpDistane
	}

	// Fill the result array from the info in the map
	for date, dist := range tempProgressInfo {
		resProgresses = append(resProgresses, Progress{
			Date:     date,
			Distance: dist,
		})
	}
	return resProgresses
}
