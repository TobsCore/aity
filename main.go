package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"strconv"
)

const port = 63971

var progress []Progress
var progressID = 3

type server struct {
	router *mux.Router
}

func main() {
	progress = append(progress, Progress{ID: "1", TrackID: "1", Distance: Distance{Value: 1500, Unit: "meter"}, Date: "07/08/2018"})
	progress = append(progress, Progress{ID: "2", TrackID: "1", Distance: Distance{Value: 3500, Unit: "meter"}, Date: "07/08/2018"})

	server := server{router: mux.NewRouter()}
	server.routes()
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), server.router))
}
