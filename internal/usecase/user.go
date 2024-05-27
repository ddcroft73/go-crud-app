package usecase

import (
	"database/sql"
	"fmt"
	"simple-crud-app/internal/datastore"
	"simple-crud-app/internal/util"
	"strings"
)

func GetUserByID(db *sql.DB, userID int64) (*datastore.User, error) {
	// Get a single uaet by ID.
	row := db.QueryRow(`SELECT id, username, email, fullname, message FROM users WHERE id = ?`, userID)

	var u datastore.User

	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.Fullname, &u.Message)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with ID %d not found", userID)
		}
		return nil, err
	}

	return &u, nil
}

func GetAllUsers(db *sql.DB) ([]datastore.User, error) {
	// Fetch all users from the database

	rows, err := db.Query(`SELECT id, username, email, fullname, message FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []datastore.User

	for rows.Next() {
		var u datastore.User
		err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.Fullname, &u.Message)
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

func CreateUser(db *sql.DB, username, email, fullname, message string) (*datastore.User, error) {
	// Create a new user in the database
	// ...

	result, err := db.Exec(
		`INSERT INTO users (username, email, fullname, message) VALUES (?, ?, ?, ?)`,
		username, email, fullname, message)

	if err != nil {
        util.WriteLog(err.Error())
		return nil, err
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	userID64 := int64(userID)

	user, err := GetUserByID(db, userID64)
	if err != nil {
		return nil, err
	}
	fmt.Printf("New user created, with ID: %d/n", userID64)
	return user, nil
}

func DeleteUser(db *sql.DB, userID int64) error {
	// Delete a user from the database based on ID
    
	result, err := db.Exec(`DELETE FROM users WHERE id = ?`, userID)
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
    
    // If there is no data being sent in updatData, set it back to the currUser.
    // THis way you only need to include the fields that are being updated.
	if currUser.Username != updateData.Username {
		updates = append(updates, "username = ?")
        
        if  updateData.Username == "" {
            args = append(args, currUser.Username)
        } else {
            args = append(args, updateData.Username)
        }		
	}

	if currUser.Email != updateData.Email {
		updates = append(updates, "email = ?")

        if updateData.Email == "" {
            args = append(args, currUser.Email)
        } else {
            args = append(args, updateData.Email)
        }		
	}

	if currUser.Fullname != updateData.Fullname {
		updates = append(updates, "fullname = ?")

        if updateData.Fullname == "" {
            args = append(args, currUser.Fullname)
        } else {
            args = append(args, updateData.Fullname)
        }		
	}

	if currUser.Message != updateData.Message {
		updates = append(updates, "message = ?")

        if updateData.Message == "" {            
	    	args = append(args, currUser.Message)
        } else {
            args = append(args, updateData.Message)
        }
	}

	// Construct the SQL statement with proper placeholders
	sqlStatement := "UPDATE users SET " + strings.Join(updates, ", ") + " WHERE id = ?"
	args = append(args, userID)
    
	_, err = db.Exec(sqlStatement, args...)
	if err != nil {
		util.WriteLog("Error executing SQL query: %s", err.Error())
		return nil, err
	}
	// Fetch the updated user to confirm changes
	updatedUser, err := GetUserByID(db, userID)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}
