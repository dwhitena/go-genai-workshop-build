package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

// VectorizedChunk is a struct that holds a vectorized chunk.
type VectorizedChunk struct {
	Id       int       `json:"id"`
	Chunk    string    `json:"chunk"`
	Vector   []float64 `json:"vector"`
	Metadata string    `json:"metadata"`
}

// VectorizedChunks is a slice of vectorized chunks.
type VectorizedChunks []VectorizedChunk

func main() {

	// Read in the ../embed/chunks_vectors.json file into a value of VectorizedChunks.
	file, err := os.Open("../embed/chunks_vectors.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var vectorizedChunks VectorizedChunks
	err = json.NewDecoder(file).Decode(&vectorizedChunks)
	if err != nil {
		log.Fatal(err)
	}

	// Get the DB connection string from env var.
	connStr := os.Getenv("DB_CONN_STR")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Loop over the vectorized chunks and insert them into the database.
	for _, vectorizedChunk := range vectorizedChunks {

		// Convert the []float64 vector value into a string that looks like '[1.7, 2.1, 3.2, etc.]'.
		vectorStr := "["
		for idx, val := range vectorizedChunk.Vector {
			vectorStr += fmt.Sprintf("%f", val)
			if idx < len(vectorizedChunk.Vector)-1 {
				vectorStr += ", "
			}
		}
		vectorStr += "]"

		// Insert the vectorized chunk into the database.
		_, err := db.Exec(
			"INSERT INTO items (id, chunk, metadata, embedding) VALUES ($1, $2, $3, $4);",
			vectorizedChunk.Id, vectorizedChunk.Chunk, vectorizedChunk.Metadata, vectorStr)
		if err != nil {
			log.Fatal(err)
		}
	}
}
