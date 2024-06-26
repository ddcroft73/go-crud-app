package datastore

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type User struct {
	ID       int64
	Username string
	Email    string
	Fullname string
	Message  string
}

func ConnectDB() (*sql.DB, error) {
	// Connect to the database, and return a pointer to the db to be used throughout
    // the application
	db, err := sql.Open("mysql", "new_user:password@tcp(localhost:3306)/goCRUD")
	if err != nil {
		return nil, err
	} else {
		log.Println("Database connected successfully.")
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func CreateTables(db *sql.DB) error {
	// Create the users table
	dbExists := checkTableExists(db, "users")
	if !dbExists {
		query := `
                CREATE TABLE users (
                    id INT AUTO_INCREMENT,
                    username TEXT NOT NULL,
                    email TEXT NOT NULL,
                    fullname TEXT NOT NULL,
                    message TEXT NOT NULL,
                    PRIMARY KEY (id)
                );`

		if _, err := db.Exec(query); err != nil {
			return err
		} else {
			log.Println("Tables created/verified successfully.")
		}
	}
	return nil
}

func checkTableExists(db *sql.DB, tableName string) bool {
	var tmp string
	query := `SELECT table_name FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?`
	err := db.QueryRow(query, tableName).Scan(&tmp)
	if err != nil && err != sql.ErrNoRows {
		log.Fatal("Query Error: ", err)
	}
	return err != sql.ErrNoRows
}
func checkUserExists(tx *sql.Tx) (bool, error) {
	// Check if any user exists, not going by ID because the ID could literally
	// be any number.
	var count int
	err := tx.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		log.Fatal(err)
	}

	if count == 0 {
		return false, nil
	} else {
		return true, nil
	}
}

func CreateInitialUserIfNeeded(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	userExists, err := checkUserExists(tx)
	if err != nil {
		return err
	}
	if userExists {
		return nil // User already exists, no need to create
	}

	if err := CreateInitialUser(tx); err != nil {
		return err
	}

	return tx.Commit()
}
func CreateInitialUser(tx *sql.Tx) error {
	_, err := tx.Exec("INSERT INTO users (username, email, fullname, message) VALUES (?, ?, ?, ?)",
		"Admin User", "admin@email.com", "John Harris", "This is the initial user on the system.")
	return err
}
