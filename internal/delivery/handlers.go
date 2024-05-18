package delivery

import (
    "database/sql"
    "net/http"
    "github.com/gorilla/mux"
    _ "simple-crud-app/internal/datastore"
    "simple-crud-app/internal/usecase"
    "encoding/json"
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
        
        // Create a slice to hold the user data
        var userList []map[string]interface{}
        
        // Iterate over the users and create a map for each user
        for _, user := range users {
            userMap := map[string]interface{}{
                "username": user.Username,
                "email":    user.Email,
                "password": user.Password,
                "message":  user.Message,
            }
            userList = append(userList, userMap)
        }
        
        // Set the response content type to JSON
        w.Header().Set("Content-Type", "application/json")
        
        // Encode the user list as JSON and write it to the response
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
         
         
        // Create a new user in the database
        err = usecase.CreateUser(db, username, email, password, message)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        // Redirect back to the home page
        http.Redirect(w, r, "/", http.StatusSeeOther)
    }
}

// Define other handlers for updating, deleting, and deleting all users
// ...