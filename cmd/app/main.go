package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"simple-crud-app/internal/datastore"
	"simple-crud-app/internal/delivery"
)

func main() {
	// Create a new router
	r := mux.NewRouter()

	db, err := datastore.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = datastore.CreateTables(db)
	if err != nil {
		log.Fatal(err)
	}

	delivery.SetupRoutes(r, db)

	log.Fatal(http.ListenAndServe(":8080", r))
}
