package main

import (
	"log"
	"net/http"
	"simple-crud-app/internal/datastore"
	"simple-crud-app/internal/delivery"
	"github.com/gorilla/mux"
)

func main() {
    
	r := mux.NewRouter()

	db, err := datastore.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = datastore.CreateTables(db)
	if err != nil {
		log.Fatal("Error creating table.",err)
	}

    if err := datastore.CreateInitialUserIfNeeded(db); err != nil {
		log.Printf("Error setting up initial user: %v", err)
	}

	delivery.SetupRoutes(r, db)

	log.Printf("App running on http://localhost:%s", "8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
