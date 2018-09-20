package main

import (
	"encoding/json"
	"errors"
	"github.com/globalsign/mgo/bson"
	"github.com/gorilla/mux"
	"github.com/tobscore/aity/model"
	"github.com/tobscore/aity/unit"
	"log"
	"net/http"
)

func (s *server) routes() {
	s.router.HandleFunc("/authenticate", s.Authenticate).Methods("POST")

	s.router.HandleFunc("/{user}/tracks", s.CreateTrack).Methods("POST")
	s.router.HandleFunc("/{user}/tracks", s.GetAllTracks).Methods("GET")

	s.router.HandleFunc("/{user}/track/current", s.GetCurrTrack).Methods("GET")
	s.router.HandleFunc("/{user}/track/{trackid:[a-z0-9]+}", s.GetTrack).Methods("GET")

	s.router.HandleFunc("/{user}/track/current/progress", s.CreateProgressForCurr).Methods("POST")
	s.router.HandleFunc("/{user}/track/current/progress", s.GetProgressForCurr).Methods("GET")

	s.router.HandleFunc("/{user}/track/{trackid:[a-z0-9]+}/progress", s.CreateProgress).Methods("POST")
	s.router.HandleFunc("/{user}/track/{trackid:[a-z0-9]+}/progress", s.GetProgress).Methods("GET")
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

// GetTrack returns information about the given track of a user
func (s *server) GetTrack(w http.ResponseWriter, r *http.Request) {
	if valid, s, c := ValidateRequest(r); !valid {
		http.Error(w, s, c)
		return
	}

	user := mux.Vars(r)["user"]
	trackId := mux.Vars(r)["trackid"]

	// Check if the given id is valid
	if !bson.IsObjectIdHex(trackId) {
		http.Error(w, "Invalid track id: "+trackId, http.StatusNotFound)
		return
	}

	track, err := s.persistence.GetTrackById(trackId, user)
	if err != nil {
		http.Error(w, "No track with id "+trackId+" found for user "+user, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(track)
}

// GetAllTracks returns information about the current track of a user.
func (s *server) GetAllTracks(w http.ResponseWriter, r *http.Request) {
	if valid, s, c := ValidateRequest(r); !valid {
		http.Error(w, s, c)
		return
	}

	user := mux.Vars(r)["user"]
	tracks, err := s.persistence.GetAllTracksByUsername(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(tracks)
}

// GetCurrTrack returns information about the current track of a user.
func (s *server) GetCurrTrack(w http.ResponseWriter, r *http.Request) {
	if valid, s, c := ValidateRequest(r); !valid {
		http.Error(w, s, c)
		return
	}

	user := mux.Vars(r)["user"]
	track, err := s.persistence.GetCurrentTrackByUsername(user)

	if err != nil {
		http.Error(w, "No track found for user "+user, http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(track)
}

func (s *server) CreateTrack(w http.ResponseWriter, r *http.Request) {
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
	trackid, err := s.persistence.CreateTrack(user, &track)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the id of the returned track to the id of the created track (db)
	track.Id = trackid

	// Update the current track so it is the newly created track
	err = s.persistence.UpdateCurrentTrack(user, track)
	// And set the returned track object to active
	track.Active = true
	if err != nil {
		log.Println("Cannot update current track for user " + user)
		http.Error(w, "Cannot update current track", http.StatusInternalServerError)
		return
	}

	// Return the saved object to the sender
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(track)
}

func (s *server) GetProgressForCurr(w http.ResponseWriter, r *http.Request) {
	// Validate user request. Is the token provided and is it still valid?
	if valid, s, c := ValidateRequest(r); !valid {
		http.Error(w, s, c)
		return
	}

	user := mux.Vars(r)["user"]
	track, err := s.persistence.GetCurrentTrackByUsername(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	accProgressList, err := s.getProgress(r, user, track.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(accProgressList)
}

func (s *server) GetProgress(w http.ResponseWriter, r *http.Request) {
	// Validate user request. Is the token provided and is it still valid?
	if valid, s, c := ValidateRequest(r); !valid {
		http.Error(w, s, c)
		return
	}

	user := mux.Vars(r)["user"]
	trackId := mux.Vars(r)["trackid"]

	// Check if the given id is valid
	if !bson.IsObjectIdHex(trackId) {
		http.Error(w, "Invalid track id: "+trackId, http.StatusNotFound)
		return
	}

	accProgressList, err := s.getProgress(r, user, trackId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(accProgressList)
}

func (s *server) getProgress(r *http.Request, user, trackId string) ([]model.Progress, error) {
	progressModelList, err := s.persistence.GetProgressByUsername(trackId, user)

	progressList := make([]model.Progress, len(progressModelList))
	if err != nil {
		return progressList, err
	}
	if len(progressModelList) == 0 {
		return progressList, errors.New("No progress entries found for user " + user + " and track " + trackId)
	}
	// Convert the list of progress model objects (DB) to usable progress objects.
	for i, progressModel := range progressModelList {
		progressList[i] = *progressModel.ToProgress()
	}

	// Generate a list of accumulated progresses, where one date has one progress distance.
	accProgressList := model.AccProgresses(progressList)
	return accProgressList, nil
}

func (s *server) CreateProgressForCurr(w http.ResponseWriter, r *http.Request) {
	if valid, s, c := ValidateRequest(r); !valid {
		http.Error(w, s, c)
		return
	}

	user := mux.Vars(r)["user"]
	track, err := s.persistence.GetCurrentTrackByUsername(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	progress, err := s.createProgress(r, user, track.Id)
	if err != nil {
		http.Error(w, "Invalid track id: "+track.Id, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(*progress)
}

func (s *server) CreateProgress(w http.ResponseWriter, r *http.Request) {
	if valid, s, c := ValidateRequest(r); !valid {
		http.Error(w, s, c)
		return
	}

	user := mux.Vars(r)["user"]
	trackId := mux.Vars(r)["trackid"]

	progress, err := s.createProgress(r, user, trackId)
	if err != nil {
		http.Error(w, "Invalid track id: "+trackId, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(*progress)
}

func (s *server) createProgress(r *http.Request, user, trackId string) (*model.Progress, error) {
	var progress model.Progress
	_ = json.NewDecoder(r.Body).Decode(&progress)

	// Check if the given id is valid
	if !bson.IsObjectIdHex(trackId) {
		return &progress, errors.New("invalid track id")
	}

	err := s.persistence.AddProgressToTrack(trackId, user, &progress)
	if err != nil {
		return &progress, err
	}
	return &progress, nil
}
