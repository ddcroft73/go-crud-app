package delivery

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"simple-crud-app/internal/datastore"
	"simple-crud-app/internal/usecase"
	"strconv"

	"github.com/gorilla/mux"
)

func SetupRoutes(r *mux.Router, db *sql.DB) {
	r.HandleFunc("/", homeHandler(db)).Methods("GET")
	r.HandleFunc("/", createUserHandler(db)).Methods("POST")
	r.HandleFunc("/update/{id}", updateUserHandler(db)).Methods("POST")
	r.HandleFunc("/delete/{id}", deleteUserHandler(db)).Methods("POST")
	r.HandleFunc("/delete-all", deleteAllUsersHandler(db)).Methods("POST")
}


func parseUserID(vars map[string]string) (int64, error) {
    idStr, ok := vars["id"]
    if !ok {
        return 0, errors.New("id is missing")
    }
    userID, err := strconv.Atoi(idStr)
    if err != nil {
        return 0, errors.New("invalid user ID")
    }
    return int64(userID), nil
}
func respondWithError(w http.ResponseWriter, code int, message string) {
    respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    json.NewEncoder(w).Encode(payload)
}

func respondWithSuccess(w http.ResponseWriter, payload interface{}) {
    respondWithJSON(w, http.StatusOK, payload)
}


// THis is the root endpoint for the GET method. When a user accesses the site, all users are gathered and
// sent to the client to populate the table.
func homeHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Fetch all users from the database
		users, err := usecase.GetAllUsers(db)
		if err != nil {
            respondWithError(w, http.StatusInternalServerError, "Failed to fetch users")
			return
		}

		var userList []map[string]interface{}
      
		for _, user := range users {

			userMap := map[string]interface{}{
				"username": user.Username,
				"email":    user.Email,
				"fullname": user.Fullname,
				"message":  user.Message,
			}
			userList = append(userList, userMap)
		}
        
        respondWithSuccess(w, userList)
	}
}

// THis Handler is for the Root POST method. it Creates a new user.
func createUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the form data
		err := r.ParseForm()
		if err != nil {
            respondWithError(w, http.StatusBadRequest, "Bad Request")
			return
		}

		username := r.FormValue("username")
		email := r.FormValue("email")
		fullname := r.FormValue("fullname")
		message := r.FormValue("message")

		var user *datastore.User

		user, err = usecase.CreateUser(db, username, email, fullname, message)
		if err != nil {
            respondWithError(w, http.StatusInternalServerError, "Error Creating New User.")			
			return
		}
        
        respondWithSuccess(w, user)
	}
}

// this function returns anoyjet funxtion, which is the actual handler
func updateUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the id from thr url.

        vars := mux.Vars(r)
        userID64, err := parseUserID(vars)
		

		if err != nil {
            respondWithError(w, http.StatusBadRequest, "Invalid User ID")			
			return
		}

		// Parse the form data
		err = r.ParseForm()
		if err != nil {
            respondWithError(w, http.StatusBadRequest, "Failed to parse form data.")
			return
		}
        
		// Get the updated user data from the form
		username := r.FormValue("username")
		email := r.FormValue("email")
		fullname := r.FormValue("fullname")
		message := r.FormValue("message")

		// Create a new User instance with the updated data
		// yo be passed into the updateUser function
		// in delivery/usecase
		updateData := &datastore.User{
			Username: username,
			Email:    email,
			Fullname: fullname,
			Message:  message,
		}

		// Update the user using the UpdateUser function
		updatedUser, err := usecase.UpdateUser(db, userID64, updateData)
		if err != nil {
            respondWithError(w, http.StatusBadRequest, err.Error())            		
			return
		}

        respondWithSuccess(w, updatedUser)
	}
}

func deleteUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the id from thr url

        vars := mux.Vars(r)
        userID64, err := parseUserID(vars)
		
		if err != nil {
            respondWithError(w, http.StatusBadRequest, "Invalid User ID")
			return
		}

		// attempt to delete the User
		err = usecase.DeleteUser(db, userID64)
		if err != nil {
            respondWithError(w, http.StatusInternalServerError, "Error deleting user")
			return
		}

        respondWithSuccess(w, map[string]string{"success": "User deleted.",})
	}
}

func deleteAllUsersHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the id from thr url.
		_, err := usecase.DeleteAll(db)
		if err != nil {
            respondWithError(w, http.StatusInternalServerError, "Error deleting all users")
			return
		}

        respondWithSuccess(w, map[string]string{"success": "All Users deleted.",})
	}
}

