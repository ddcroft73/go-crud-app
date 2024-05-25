package delivery

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"simple-crud-app/internal/datastore"
	"simple-crud-app/internal/usecase"
	"strconv"
)

func SetupRoutes(r *mux.Router, db *sql.DB) {
	r.HandleFunc("/", homeHandler(db)).Methods("GET")
	r.HandleFunc("/", createUserHandler(db)).Methods("POST")
	r.HandleFunc("/update/{id}", updateUserHandler(db)).Methods("POST")
	r.HandleFunc("/delete/{id}", deleteUserHandler(db)).Methods("POST")
	r.HandleFunc("/delete-all", deleteAllUsersHandler(db)).Methods("POST")
}

// THis is the root endpoint for the GET method. When a user accesses the site, all users are gathered and
// sent to the client to populate the table.
func homeHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Fetch all users from the database
		users, err := usecase.GetAllUsers(db)
		if err != nil {
			errorResponse := map[string]string{
				"error": "Failed to fetch users",
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errorResponse)
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

// THis Handler is for the Root POST method. it Creates a new user.
func createUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the form data
		err := r.ParseForm()
		if err != nil {
			errorResponse := map[string]string{
				"error": "Bad Request",
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errorResponse)
			return
		}

		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")
		message := r.FormValue("message")

		var user *datastore.User

		user, err = usecase.CreateUser(db, username, email, password, message)
		if err != nil {
			errorResponse := map[string]string{
				"error": "Error Creating New User.",
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errorResponse)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}

// this function returns anoyjet funxtion, which is the actual handler
func updateUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the id from thr url.
		vars := mux.Vars(r)
		idStr := vars["id"]

		// Convert the 'id' to an integer
		userID, err := strconv.Atoi(idStr)
		userID64 := int64(userID)

		if err != nil {
			errorResponse := map[string]string{
				"error": "Invalid User ID",
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errorResponse)
			return
		}

		// Parse the form data
		err = r.ParseForm()
		if err != nil {
			errorResponse := map[string]string{
				"error": "Failed to parse form data.",
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errorResponse)
			return
		}

		// Get the updated user data from the form
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")
		message := r.FormValue("message")

		// Create a new User instance with the updated data
		// yo be passed into the updateUser function
		// in delivery/usecase
		updateData := &datastore.User{
			Username: username,
			Email:    email,
			Password: password,
			Message:  message,
		}

		// Update the user using the UpdateUser function
		updatedUser, err := usecase.UpdateUser(db, userID64, updateData)
		if err != nil {
			errorResponse := map[string]string{
				"error": "Failed to update user.",
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errorResponse)
			return
		}

		// Write a success response
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(updatedUser)
	}
}

func deleteUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the id from thr url
		vars := mux.Vars(r)
		idStr := vars["id"]

		// Convert the 'id' to an integer
		userID, err := strconv.Atoi(idStr)
		userID64 := int64(userID)

		if err != nil {
			errorResponse := map[string]string{
				"error": "Invalid User ID",
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errorResponse)
			return
		}

		// attempt to delete the User
		err = usecase.DeleteUser(db, userID64)
		if err != nil {
			errorResponse := map[string]string{
				"error": "Error deleting user",
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errorResponse)
			return
		}

		successResponse := map[string]string{
			"success": "User deleted.",
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(successResponse)

	}
}

func deleteAllUsersHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the id from thr url.
		_, err := usecase.DeleteAll(db)
		if err != nil {
			errorResponse := map[string]string{
				"error": "Error deleting all users",
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errorResponse)
			return
		}

		successResponse := map[string]string{
			"success": "All Users deleted.",
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(successResponse)

	}
}
