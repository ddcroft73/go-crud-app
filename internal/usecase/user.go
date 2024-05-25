package usecase

import (
	"database/sql"
	"fmt"
	"simple-crud-app/internal/datastore"
	"strings"
)

func GetUserByID(db *sql.DB, userID int64) (*datastore.User, error) {
	// Get a single uaet by ID.
	row := db.QueryRow(`SELECT id, username, email, password, message FROM users WHERE id = ?`, userID)

	var u datastore.User

	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.Message)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("User with ID %d not found", userID)
		}
		return nil, err
	}

	return &u, nil
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

func CreateUser(db *sql.DB, username, email, password, message string) (*datastore.User, error) {
	// Create a new user in the database
	// ...

	result, err := db.Exec(
		`INSERT INTO users (username, email, password, message) VALUES (?, ?, ?, ?)`,
		username, password, email, message)

	if err != nil {
		return nil, err
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	user, err := GetUserByID(db, userID)
	if err != nil {
		return nil, err
	}
	fmt.Printf("New user created, with ID: %d/n", userID)
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

func DeleteAll(db *sql.DB) (bool, error) {

	currUsers, err := GetAllUsers(db)
	if err != nil {
		return false, err
	}

	for _, user := range currUsers {
		err = DeleteUser(db, user.ID)
		if err != nil {
			return false, err
		}
		fmt.Println("Deleted User ID: ", user.ID)
	}

	fmt.Println("Deleted all users")
	return true, nil
}

func UpdateUser(db *sql.DB, userID int64, updateData *datastore.User) (*datastore.User, error) {
	// Fetch the current user to compare fields
	currUser, err := GetUserByID(db, userID)
	if err != nil {
		return nil, err
	}

	var updates []string
	var args []interface{}
	idx := 1

	if currUser.Username != updateData.Username {
		updates = append(updates, fmt.Sprintf("username = $%d", idx))
		args = append(args, updateData.Username)
		idx++
	}

	if currUser.Email != updateData.Email {
		updates = append(updates, fmt.Sprintf("email = $%d", idx))
		args = append(args, updateData.Email)
		idx++
	}

	if currUser.Password != updateData.Password {
		updates = append(updates, fmt.Sprintf("password = $%d", idx))
		args = append(args, updateData.Password)
		idx++
	}

	if currUser.Message != updateData.Message {
		updates = append(updates, fmt.Sprintf("message = $%d", idx))
		args = append(args, updateData.Message)
		idx++
	}

	if len(updates) == 0 {
		return currUser, nil // No updates needed
	}

	sqlStatement := "UPDATE users SET " + strings.Join(updates, ", ") + fmt.Sprintf(" WHERE id = $%d", idx)
	args = append(args, userID)

	_, err = db.Exec(sqlStatement, args...)
	if err != nil {
		return nil, err
	}

	// Fetch the updated user to confirm changes
	updatedUser, err := GetUserByID(db, userID)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}
