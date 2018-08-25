package main

import (
	"errors"
	"strings"
)

// The persistence is a layer that stores the users' data and abstracts access over these data.
type Persistence struct {
	tracks map[string]Track
}

func initDefaultPersistance() *Persistence {
	tracks := make(map[string]Track)
	p := Persistence{tracks: tracks}

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
	usernameLower := strings.ToLower(username)
	if res, ok := p.tracks[usernameLower]; ok {
		return res, nil
	}
	return Track{}, errors.New("Cannot find track for user " + username)
}

func (p *Persistence) addTrack(username string, track Track) bool {
	// TODO: Probably should check the input
	p.tracks[username] = track
	return true
}
func (p *Persistence) addProgress(username string, progress Progress) (*Progress, error) {
	// Only add progress if user has a track
	resProgress := Progress{}
	if _, err := p.getTrackByUsername(username); err != nil {
		return &resProgress, nil
	}

	return &Progress{}, nil
}
