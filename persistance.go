package main

import (
	"errors"
	"strings"
)

// The persistence is a layer that stores the users' data and abstracts access over these data.
type Persistence struct {
	tracks     map[string]Track
	progresses map[string][]Progress
}

func initDefaultPersistance() *Persistence {
	tracks := make(map[string]Track)
	progresses := make(map[string][]Progress)
	p := Persistence{tracks: tracks, progresses:progresses}

	// Add example track as to test the core.
	p.tracks["tobscore"] = Track{
		Start: Location{
			Name:      "Karlsruhe",
			Latitude:  8.4044370,
			Longitude: 49.013506,
		},
		End: Location{
			Name:      "Berlin",
			Latitude:  13.377637,
			Longitude: 52.516275,
		},
		Distance: 624000,
	}
	return &p
}

func (p *Persistence) getTrackByUsername(username string) (Track, error) {
	if res, ok := p.tracks[toLower(username)]; ok {
		return res, nil
	}
	return Track{}, errors.New("Cannot find track for user " + username)
}

func (p *Persistence) addTrack(username string, track Track) bool {
	// TODO: Probably should check the input
	p.tracks[username] = track
	return true
}

func (p *Persistence) getProgressByUsername(username string) []Progress {
	return p.progresses[toLower(username)]
}

func (p *Persistence) accProgresses(username string) []Progress {
	allProgresses := p.getProgressByUsername(username)
	tempProgressInfo := make(map[string]Distance)
	var resProgresses []Progress
	for _, prog := range allProgresses {
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

// AddProgress adds the given progress to the list of progresses for the given user. It then returns the user's daily progress. If there was progress stored for the day previously, the given amount will be added upon the old value.
func (p *Persistence) addProgress(username string, progress Progress) Progress {
	progressMap := p.progresses
	progressList := progressMap[toLower(username)]
	progressMap[toLower(username)] = append(progressList, progress)
	return progress
}

func (p *Persistence) userKnown(username string) bool {
	_, exists := p.tracks[toLower(username)]
	return exists
}

func toLower(username string) string {
	return strings.ToLower(username)
}
