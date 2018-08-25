package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/tobscore/walk/unit"
	"net/http"
)

func (s *server) routes() {
	s.router.HandleFunc("/{username}/track", s.GetCurrentTrackInfo).Methods("GET")
	s.router.HandleFunc("/{username}/track", s.CreateNewTrack).Methods("POST")
	s.router.HandleFunc("/{username}/track/progress", s.GetCurrentTrackProgress).Methods("GET")
	s.router.HandleFunc("/{username}/track/progress", s.AddToCurrentTrackProgress).Methods("POST")
}

// GetCurrentTrackInfo returns information about the current track of a user.
func (s *server) GetCurrentTrackInfo(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	track, err := s.persistence.getTrackByUsername(username)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	} else {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(track)
	}
}

// GetCurrentTrackInfo returns information about the current track of a user.
func (s *server) CreateNewTrack(w http.ResponseWriter, r *http.Request) {
	var track Track
	_ = json.NewDecoder(r.Body).Decode(&track)
	username := mux.Vars(r)["username"]

	// Check if names have been set for start and end
	if track.Start.Name == "" {
		http.Error(w, "Start name cannot be empty", http.StatusBadRequest)
		return
	} else if track.End.Name == "" {
		http.Error(w, "End name cannot be empty", http.StatusBadRequest)
		return
	}
	// Calculate the difference between the two
	from, err := unit.NewCoord(track.Start.Latitude, track.Start.Longitude)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	to, err := unit.NewCoord(track.End.Latitude, track.End.Longitude)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	distance := unit.Distance(from, to)

	// Set the newly calculated distance to the track and save it in the persistence layer
	track.Distance = Distance(distance)
	s.persistence.addTrack(username, track)

	// Return the saved object to the sender
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(track)
}

func (s *server) GetCurrentTrackProgress(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (s *server) AddToCurrentTrackProgress(w http.ResponseWriter, r *http.Request) {
	var progress Progress
	_ = json.NewDecoder(r.Body).Decode(&progress)
	username := mux.Vars(r)["username"]

	dailyProgress, err := s.persistence.addProgress(username, progress)
}
