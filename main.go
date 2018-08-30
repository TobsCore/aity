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
	session, err := mongo.NewSession("localhost:27017")
	if err != nil {
		panic("Cannot connect to db")
	}
	pers := mongo.NewStorageService(session, "aity")
	server := server{router: mux.NewRouter(), persistence: pers}
	server.routes()
	log.Printf("Starting server on port %d\n", port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), server.router))
}
