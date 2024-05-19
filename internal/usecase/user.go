package usecase

import (
	"database/sql"
	"fmt"
	"simple-crud-app/internal/datastore"
	"strings"
)

func GetUserByID(db *sql.DB, userID int64) (datastore.User, error) {
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

func CreateUser(db *sql.DB, username, email, password, message string) (datastore.User, error) {
	// Create a new user in the database, on success return the users data
	// on error, return the error
	// ...

	var user datastore.User

	result, err := db.Exec(
		`INSERT INTO users (username, email, password,message) VALUES (?, ?, ?, ?)`,
		username, password, email, message)

	if err != nil {
		return datastore.User{}, err
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return datastore.User{}, err
	}

	fmt.Printf("New user created, with ID: %d/n", userID)

	user, err = GetUserByID(db, userID)
	if err != nil {
		return datastore.User{}, err
	}
	return user, nil
}

func DeleteUser(db *sql.DB, userID int64) error {
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

func UpdateUser(db *sql.DB, userID int64, updateData *datastore.User) (datastore.User, error) {
	// 5Update a user by the ID.

	var u datastore.User

	u, err := GetUserByID(db, userID)

	if err != nil {
		return datastore.User{}, err
	}

	// holds the new data to be updated
	var updates []string

	if u.Username != updateData.Username {
		updates = append(updates, `username = '${updateData.Username}'`)
	}

	if u.Email != updateData.Email {
		updates = append(updates, `email = '${updateData.Email}'`)
	}

	if u.Password != updateData.Password {
		updates = append(updates, `password = '${updateData.Password}'`)
	}

	if u.Message != updateData.Message {
		updates = append(updates, `message = '${updateData.Message}'`)
	}

	sql := "UPDATE users SET " + strings.Join(updates, ", ") + " WHERE id = ?"

	_, err = db.Exec(sql, userID)

	if err != nil {
		return datastore.User{}, err
	}

	u, err = GetUserByID(db, userID)

	if err != nil {
		return datastore.User{}, err
	}

	return u, nil
}
