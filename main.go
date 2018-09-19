package main

import (
	"github.com/tobscore/aity/mongo"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"strconv"
)

const port = 63971

type server struct {
	router      *mux.Router
	persistence *mongo.StorageService
}

func main() {
	mc, err := ReadConf()
	if err != nil {
		log.Fatal(err)
	}
	session, err := mongo.NewSession(mc)
	if err != nil {
		log.Fatal("Cannot connect to db")
	}
	p := mongo.NewStorageService(session, "aity")
	server := server{router: mux.NewRouter(), persistence: p}
	server.routes()

	// Start the server
	log.Printf("Starting server on port %d\n", port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), server.router))
}
