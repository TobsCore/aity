package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/tobscore/aity/model"
	"github.com/tobscore/aity/unit"
	"net/http"
)

func (s *server) routes() {
	s.router.HandleFunc("/authenticate", s.Authenticate).Methods("POST")

	s.router.HandleFunc("/{username}/track", s.GetCurrentTrackInfo).Methods("GET")
	s.router.HandleFunc("/{username}/track", s.CreateNewTrack).Methods("POST")
	s.router.HandleFunc("/{username}/track/progress", s.GetCurrentTrackProgress).Methods("GET")
	s.router.HandleFunc("/{username}/track/progress", s.AddToCurrentTrackProgress).Methods("POST")
}

func (s *server) Authenticate(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("token")
	if len(token) == 0 {
		http.Error(w, "Auth token missing", http.StatusBadRequest)
		return
	}

	u, err, statusCode := lookupGoogleUser(token)
	if err != nil {
		http.Error(w, err.Error(), statusCode)
		return
	}
	var exists bool
	var user model.User
	exists = s.persistence.UserExists(u.Email)
	if exists {
		user, _ = s.persistence.GetUserByEmail(u.Email)
	} else {
		user, _ = s.persistence.CreateUser(u.toUser())
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(user)
}

// GetCurrentTrackInfo returns information about the current track of a user.
func (s *server) GetCurrentTrackInfo(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	track, err := s.persistence.GetTrackByUsername(username)

	if err != nil {
		http.Error(w, "No track found for user "+username, http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(track)
}

func (s *server) CreateNewTrack(w http.ResponseWriter, r *http.Request) {
	var track model.Track
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
	track.Distance = model.Distance(distance)
	_, err = s.persistence.CreateTrack(username, &track)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the saved object to the sender
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(track)
}

func (s *server) GetCurrentTrackProgress(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	progressList, err := s.persistence.GetProgressByUsername(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(progressList) == 0 {
		http.Error(w, "No entries found for user "+username, http.StatusNotFound)
		return
	}
	accumulatedProgressList := model.AccProgresses(progressList)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(accumulatedProgressList)
}

func (s *server) AddToCurrentTrackProgress(w http.ResponseWriter, r *http.Request) {
	var progress model.Progress
	_ = json.NewDecoder(r.Body).Decode(&progress)
	username := mux.Vars(r)["username"]

	err := s.persistence.AddProgress(username, &progress)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(progress)
}
