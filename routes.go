package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/tobscore/aity/model"
	"github.com/tobscore/aity/unit"
	"log"
	"net/http"
)

func (s *server) routes() {
	s.router.HandleFunc("/authenticate", s.Authenticate).Methods("POST")

	s.router.HandleFunc("/{user}/track", s.GetCurrentTrackInfo).Methods("GET")
	s.router.HandleFunc("/{user}/track", s.CreateNewTrack).Methods("PUT")
	s.router.HandleFunc("/{user}/track", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "", http.StatusNotImplemented)
		return
	}).Methods("POST")
	s.router.HandleFunc("/{user}/track/progress", s.GetCurrentTrackProgress).Methods("GET")
	s.router.HandleFunc("/{user}/track/progress", s.AddToCurrentTrackProgress).Methods("POST")
}

func (s *server) Authenticate(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("token")
	if len(token) == 0 {
		http.Error(w, "Auth token missing", http.StatusBadRequest)
		return
	}

	u, err, statusCode := model.LookupGoogleUser(token)
	if err != nil {
		http.Error(w, err.Error(), statusCode)
		return
	}

	// Check if the user exists already. If it does the user is returned. If it doesn't the user document is created in the database and the username is set to the user's name (which coming from google is the actual name).
	var exists bool
	var user model.User
	exists = s.persistence.UserExists(u.Email)
	if exists {
		user, _ = s.persistence.GetUserByEmail(u.Email)
	} else {
		user = *u.ToUser()
		err = s.persistence.CreateUser(&user)
		if err != nil {
			log.Println(err)
		}
	}

	// Create a JWT Token for the user with the application's secret
	authToken, err := TokenForUser(user.Email)
	if err != nil {
		log.Printf("Error generating token for user %+v\n", user)
		log.Println(err.Error())
		http.Error(w, "Error generating JWT token", http.StatusInternalServerError)
		return
	}

	// Create an auth response that contains information, whether the user existed before and the user information.
	res := &model.AuthResponse{
		AlreadyRegistered:  exists,
		UserInfo:           user,
		AuthToken:          authToken,
		ExpirationDuration: tokenExpiration.Seconds(),
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(res)
}

// GetCurrentTrackInfo returns information about the current track of a user.
func (s *server) GetCurrentTrackInfo(w http.ResponseWriter, r *http.Request) {
	if valid, s, c := ValidateRequest(r); !valid {
		http.Error(w, s, c)
		return
	}

	user := mux.Vars(r)["user"]
	track, err := s.persistence.GetTrackByUsername(user)

	if err != nil {
		http.Error(w, "No track found for user "+user, http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(track)
}

func (s *server) CreateNewTrack(w http.ResponseWriter, r *http.Request) {
	if valid, s, c := ValidateRequest(r); !valid {
		http.Error(w, s, c)
		return
	}

	var track model.Track
	_ = json.NewDecoder(r.Body).Decode(&track)
	user := mux.Vars(r)["user"]

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
	_, err = s.persistence.CreateTrack(user, &track)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the saved object to the sender
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(track)
}

func (s *server) GetCurrentTrackProgress(w http.ResponseWriter, r *http.Request) {
	if valid, s, c := ValidateRequest(r); !valid {
		http.Error(w, s, c)
		return
	}

	user := mux.Vars(r)["user"]
	progressModelList, err := s.persistence.GetProgressByUsername(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(progressModelList) == 0 {
		http.Error(w, "No entries found for user "+user, http.StatusNotFound)
		return
	}
	// Convert the list of progress model objects (DB) to usable progress objects.
	progressList := make([]model.Progress, len(progressModelList))
	for i, progressModel := range progressModelList {
		progressList[i] = *progressModel.ToProgress()
	}

	// Generate a list of accumulated progresses, where one date has one progress distance.
	accumulatedProgressList := model.AccProgresses(progressList)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(accumulatedProgressList)
}

func (s *server) AddToCurrentTrackProgress(w http.ResponseWriter, r *http.Request) {
	if valid, s, c := ValidateRequest(r); !valid {
		http.Error(w, s, c)
		return
	}

	var progress model.Progress
	_ = json.NewDecoder(r.Body).Decode(&progress)
	user := mux.Vars(r)["user"]

	err := s.persistence.AddProgress(user, &progress)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(progress)
}
