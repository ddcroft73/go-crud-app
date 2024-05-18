package datastore

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)


type User struct {
    Username string
    Email string
    Password string
    Message string
}


func ConnectDB() (*sql.DB, error) {
    // Connect to the database
    db, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/dbname")
    if err != nil {
        return nil, err
    }
    
    // new
    if err := db.Ping(); err != nil {
        return nil, err
    }

    return db, nil
}

func CreateTables(db *sql.DB) error {
    // Create the users table
    query := `
            CREATE TABLE users (
                id INT AUTO_INCREMENT,
                username TEXT NOT NULL,
                password TEXT NOT NULL,
                created_at DATETIME,
                PRIMARY KEY (id)
            );`
            
    if _, err := db.Exec(query); err != nil {
        return err
    }        

    return nil
}


//_, err := db.Exec(`
//CREATE TABLE IF NOT EXISTS users (
//    id INT AUTO_INCREMENT PRIMARY KEY,
//    username VARCHAR(255) NOT NULL,
//    email VARCHAR(255) NOT NULL,
//    password VARCHAR(255) NOT NULL,
//    message TEXT
//)
//`)

//if err != nil {
//return err
//}
//return nil