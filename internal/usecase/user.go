package usecase

import (
	"database/sql"
	"fmt"
	"simple-crud-app/internal/datastore"
)

func GetUser(db *sql.DB, userID int) (datastore.User, error) {
	// Get a single uaet by ID.
	row := db.QueryRow(`SELECT id, username, email, password, message FROM users WHERE id = ?`, userID)

	var u datastore.User

	err := row.Scan(&u.ID, &u.Username, &u.Password, &u.Email, &u.Message)
	if err != nil {
		if err == sql.ErrNoRows {
			return u,
				fmt.Errorf("user with ID %d not found", userID)
		}
		return u, err
	}

	return u, nil
}

func GetAllUsers(db *sql.DB) ([]datastore.User, error) {
	// Fetch all users from the database

	rows, err := db.Query(`SELECT id, username, email, password, message FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []datastore.User

	for rows.Next() {
		var u datastore.User
		err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.Message)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func CreateUser(db *sql.DB, username, email, password, message string) error {
	// Create a new user in the database
	// ...
	result, err := db.Exec(
		`INSERT INTO users (username, email, password,message) VALUES (?, ?, ?, ?)`,
		username, password, email, message)

	if err != nil {
		return err
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	fmt.Printf("New user created, with ID: %d/n", userID)
	return nil
}

func DeleteUser(db *sql.DB, userID int) error {
	// Delete a user from thr database based on ID

	result, err := db.Exec(`DELETE FROM users WHERE id = ?`, userID)
	defer db.Close()

	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	fmt.Println("Rows affected: ", rowsAffected)
	return nil
}

func DeleteAll(db *sql.DB) {
	// Delete all usersbin DB

}

func UpdateUser(db *sql.DB) {

}
