package delivery

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	_ "simple-crud-app/internal/datastore"
	"simple-crud-app/internal/usecase"
)

func SetupRoutes(r *mux.Router, db *sql.DB) {
	r.HandleFunc("/", homeHandler(db)).Methods("GET")
	r.HandleFunc("/", createUserHandler(db)).Methods("POST")
	//r.HandleFunc("/update/{id}", updateUserHandler(db)).Methods("POST")
	//r.HandleFunc("/delete/{id}", deleteUserHandler(db)).Methods("POST")
	// r.HandleFunc("/delete-all", deleteAllUsersHandler(db)).Methods("POST")
}

// THis is the root endpoint for the GET method. When a user accesses the site, all users are gathered and
// sent to the client to populate the table.
func homeHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Fetch all users from the database
		users, err := usecase.GetAllUsers(db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var userList []map[string]interface{}

		for _, user := range users {
			userMap := map[string]interface{}{
				"username": user.Username,
				"email":    user.Email,
				"password": user.Password,
				"message":  user.Message,
			}
			userList = append(userList, userMap)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(userList)
	}
}

// THis Handler is for the Roor POST method. it Creates a new user.
func createUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the form data
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")
		message := r.FormValue("message")

		err = usecase.CreateUser(db, username, email, password, message)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// Define other handlers for updating, deleting, and deleting all users
// ...
