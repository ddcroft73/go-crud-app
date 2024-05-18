package usecase

import (
    "database/sql"
    "simple-crud-app/internal/datastore"
)

func GetAllUsers(db *sql.DB) ([]datastore.User, error) {
    // Fetch all users from the database
    // ...
    return nil, nil
}

func CreateUser(db *sql.DB, username, email, password, message string) error {
    // Create a new user in the database
    // ...
    return nil
}

// Implement other functions for updating, deleting, and deleting all users
// ...