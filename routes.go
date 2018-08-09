package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func (s *server) routes() {
	s.router.HandleFunc("/{username}/current", s.GetCurrentTrackInfo).Methods("GET")
	s.router.HandleFunc("/{username}/current/last", s.GetCurrentTrack).Methods("GET")
	s.router.HandleFunc("/{username}/current/track", s.AddDistanceToCurrent).Methods("POST")
}

func (s *server) GetCurrentTrack(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(progress)
}

func (s *server) AddDistanceToCurrent(w http.ResponseWriter, r *http.Request) {
	var inputProgress Progress
	_ = json.NewDecoder(r.Body).Decode(&inputProgress)
	inputProgress.ID = strconv.Itoa(progressID)
	progressID++
	progress = append(progress, inputProgress)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(progress)
}

// GetCurrentTrackInfo returns information about the current track of a user.
func (s *server) GetCurrentTrackInfo(w http.ResponseWriter, r *http.Request) {

}
