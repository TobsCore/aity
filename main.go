package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"strconv"
)

const port = 63971

type server struct {
	router      *mux.Router
	persistence *Persistence
}

func main() {
	var p = initDefaultPersistance()
	server := server{router: mux.NewRouter(), persistence: p}
	server.routes()
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), server.router))
}
