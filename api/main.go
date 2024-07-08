package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

func init() {

	// Get the DB connection string from env var.
	connStr := os.Getenv("DB_CONN_STR")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Check the connection.
	rows, err := db.Query("select version()")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var version string
	for rows.Next() {
		err := rows.Scan(&version)
		if err != nil {
			log.Fatal(err)
		}
	}
	log.Printf("Checked DB connection: version=%s\n", version)
}

func main() {

	// ListenAndServe starts an HTTP server with a given address and
	// handler defined in NewRouter.
	log.Println("🎧 Starting to listen on port 8080!")
	router := NewRouter()
	log.Fatal(http.ListenAndServe(":8080", router))
}
