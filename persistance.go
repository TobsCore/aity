package main

import (
	"github.com/tobscore/aity/model"
	"strings"
)

// The persistence is a layer that stores the users' data and abstracts access over these data.
type Persistence struct {
	tracks     map[string]model.Track
	progresses map[string][]model.Progress
}

func initDefaultPersistance() *Persistence {
	tracks := make(map[string]model.Track)
	progresses := make(map[string][]model.Progress)
	p := Persistence{tracks: tracks, progresses: progresses}

	// Add example track as to test the core.
	p.tracks["tobscore"] = model.Track{
		Start: model.Location{
			Name:      "Karlsruhe",
			Latitude:  8.4044370,
			Longitude: 49.013506,
		},
		End: model.Location{
			Name:      "Berlin",
			Latitude:  13.377637,
			Longitude: 52.516275,
		},
		Distance: 624000,
	}
	return &p
}

func (p *Persistence) GetProgressByUsername(username string) []model.Progress {
	return p.progresses[toLower(username)]
}

func (p *Persistence) AccProgresses(username string) []model.Progress {
	allProgresses := p.GetProgressByUsername(username)
	tempProgressInfo := make(map[string]model.Distance, len(allProgresses))
	var resProgresses []model.Progress
	for _, prog := range allProgresses {
		tmpDistane := tempProgressInfo[prog.Date] + prog.Distance
		tempProgressInfo[prog.Date] = tmpDistane
	}

	// Fill the result array from the info in the map
	for date, dist := range tempProgressInfo {
		resProgresses = append(resProgresses, model.Progress{
			Date:     date,
			Distance: dist,
		})
	}
	return resProgresses
}

// AddProgress adds the given progress to the list of progresses for the given user. It then returns the user's daily progress. If there was progress stored for the day previously, the given amount will be added upon the old value.
func (p *Persistence) AddProgress(username string, progress model.Progress) model.Progress {
	progressMap := p.progresses
	progressList := progressMap[toLower(username)]
	progressMap[toLower(username)] = append(progressList, progress)
	return progress
}

func (p *Persistence) UserKnown(username string) bool {
	_, exists := p.tracks[toLower(username)]
	return exists
}

func toLower(username string) string {
	return strings.ToLower(username)
}
