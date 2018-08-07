package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"encoding/json"
	"strconv"
)

type Progress struct {
	ID               string         `json:"id,omitempty"`
	TrackID          string         `json:"trackid"`
	ProgressInMeters ProgressFormat `json:"progress"`
	Date             string         `json:"date"`
}

type ProgressFormat struct {
	Distance int `json:"distance"`
	Unit     string `json:"unit"`
}

const port = 63971
var progress []Progress
var progressID = 3

func main() {
	progress = append(progress, Progress{ID: "1", TrackID: "1", ProgressInMeters: ProgressFormat{Distance: 1500, Unit: "meter"}, Date: "07/08/2018"})
	progress = append(progress, Progress{ID: "2", TrackID: "1", ProgressInMeters: ProgressFormat{Distance: 3500, Unit: "meter"}, Date: "07/08/2018"})

	router := mux.NewRouter()
	router.HandleFunc("/{userid}/current/last", GetCurrentTrack).Methods("GET")
	router.HandleFunc("/{userid}/current/track", AddDistanceToCurrent).Methods("POST")
	log.Fatal(http.ListenAndServe(":" + strconv.Itoa(port), router))
}

func GetCurrentTrack(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(progress)
}

func AddDistanceToCurrent(w http.ResponseWriter, r *http.Request) {
	var inputProgress Progress
	_ = json.NewDecoder(r.Body).Decode(&inputProgress)
	inputProgress.ID = strconv.Itoa(progressID)
	progressID++
	progress = append(progress, inputProgress)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(progress)
}
