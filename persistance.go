package main

import (
	"errors"
	"github.com/globalsign/mgo"
	"github.com/tobscore/aity/model"
	"log"
	"strings"
)

// The persistence is a layer that stores the users' data and abstracts access over these data.
type Persistence struct {
	tracks     map[string]model.Track
	progresses map[string][]model.Progress
}

type MoviesDAO struct {
	Server   string
	Database string
}

var db *mgo.Database

const (
	TRACKS_COLL = "tracks"
)

func (m *MoviesDAO) Connect() {
	session, err := mgo.Dial(m.Server)
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB(m.Database)
}


func initDefaultPersistance() *Persistence {
	tracks := make(map[string]model.Track)
	progresses := make(map[string][]model.Progress)
	p := Persistence{tracks: tracks, progresses:progresses}

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

func (p *Persistence) getTrackByUsername(username string) (model.Track, error) {
	if res, ok := p.tracks[toLower(username)]; ok {
		return res, nil
	}
	return model.Track{}, errors.New("Cannot find track for user " + username)
}

func (p *Persistence) addTrack(username string, track model.Track) bool {
	// TODO: Probably should check the input
	p.tracks[username] = track
	return true
}

func (p *Persistence) getProgressByUsername(username string) []model.Progress {
	return p.progresses[toLower(username)]
}

func (p *Persistence) accProgresses(username string) []model.Progress {
	allProgresses := p.getProgressByUsername(username)
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
func (p *Persistence) addProgress(username string, progress model.Progress) model.Progress {
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
