package config

import (
    "database/sql"
    "fmt"
    "log"
    "os"

    _ "github.com/go-sql-driver/mysql"
)

func ConnectDB() *sql.DB {
    // Read environment variable
    dbURL := os.Getenv("DB_URL")
    if dbURL == "" {
        dbURL = ""
    }

    db, err := sql.Open("mysql", dbURL)
    if err != nil {
        log.Fatalf("Error to open conexion: %v\n", err)
    }

    // Check the connection
    if err := db.Ping(); err != nil {
        log.Fatalf("Error to connect with DB: %v\n", err)
    }

    fmt.Println("Connection successful")
    return db
}
