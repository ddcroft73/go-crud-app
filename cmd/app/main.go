package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"simple-crud-app/internal/datastore"
	"simple-crud-app/internal/delivery"
	"github.com/gorilla/handlers"
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
		log.Fatal("Error creating table.", err)
	}

	if err := datastore.CreateInitialUserIfNeeded(db); err != nil {
		log.Printf("Error setting up initial user: %v", err)
	}

	delivery.SetupRoutes(r, db)
    

	// Enable CORS
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)(r)

	log.Printf("App running on http://localhost:%s", "8080")
	log.Fatal(http.ListenAndServe(":8080", corsHandler))
}
