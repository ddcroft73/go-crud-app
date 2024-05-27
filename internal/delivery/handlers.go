package delivery

import (
	"database/sql"
	"github.com/gorilla/mux"
	"net/http"
	"simple-crud-app/internal/datastore"
	"simple-crud-app/internal/usecase"
	"simple-crud-app/internal/util"
)

func SetupRoutes(r *mux.Router, db *sql.DB) {
	r.HandleFunc("/", homeHandler(db)).Methods("GET")
	r.HandleFunc("/create-user", createUserHandler(db)).Methods("POST")
	r.HandleFunc("/update/{id}", updateUserHandler(db)).Methods("PUT")
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
			util.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch users")
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

		util.RespondWithSuccess(w, userList)
	}
}

// THis Handler is for the Root POST method. it Creates a new user.
func createUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the form data
		err := r.ParseForm()
		if err != nil {
			util.RespondWithError(w, http.StatusBadRequest, "Bad Request")
			return
		}

		username := r.FormValue("username")
		email := r.FormValue("email")
		fullname := r.FormValue("fullname")
		message := r.FormValue("message")

		var user *datastore.User
        
		user, err = usecase.CreateUser(db, username, email, fullname, message)
		if err != nil {
			util.RespondWithError(w, http.StatusInternalServerError, "Error Creating New User.")
			return
		}

		util.RespondWithSuccess(w, user)
	}
}

// this function returns anoyjet funxtion, which is the actual handler
func updateUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the id from thr url.

		vars := mux.Vars(r)
		userID64, err := util.ParseUserID(vars)

		if err != nil {
			util.RespondWithError(w, http.StatusBadRequest, "Invalid User ID")
			return
		}

		// Parse the form data
		err = r.ParseForm()
		if err != nil {
			util.RespondWithError(w, http.StatusBadRequest, "Failed to parse form data.")
			return
		}

		username := r.FormValue("username")
		email := r.FormValue("email")
		fullname := r.FormValue("fullname")
		message := r.FormValue("message")
        

		updateData := &datastore.User{
			Username: username,
			Email:    email,
			Fullname: fullname,
			Message:  message,
		}

		// Update the user using the UpdateUser function
		updatedUser, err := usecase.UpdateUser(db, userID64, updateData)
		if err != nil {
			util.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		util.RespondWithSuccess(w, updatedUser)
	}
}

func deleteUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the id from thr url

		vars := mux.Vars(r)
		userID64, err := util.ParseUserID(vars)

		if err != nil {
			util.RespondWithError(w, http.StatusBadRequest, "Invalid User ID")
			return
		}

		// attempt to delete the User
		err = usecase.DeleteUser(db, userID64)
		if err != nil {
			util.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
        
		util.RespondWithSuccess(w, map[string]string{"success": "User deleted."})
	}
}

func deleteAllUsersHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the id from thr url.
		_, err := usecase.DeleteAll(db)
		if err != nil {
			util.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		util.RespondWithSuccess(w, map[string]string{"success": "All Users deleted."})
	}
}
