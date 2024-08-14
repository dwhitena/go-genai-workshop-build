package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/predictionguard/go-client"
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

func embed(imageFile string, text string) (*VectorizedChunk, error) {

	logger := func(ctx context.Context, msg string, v ...any) {
		s := fmt.Sprintf("msg: %s", msg)
		log.Println(s)
	}

	cln := client.New(logger, host, apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var image client.ImageFile
	if imageFile != "" {
		imageParsed, err := client.NewImageFile(imageFile)
		if err != nil {
			return nil, fmt.Errorf("ERROR: %w", err)
		}
		image = imageParsed
	}

	input := []client.EmbeddingInput{
		{
			Text: text,
		},
	}
	if imageFile != "" {
		input[0].Image = image
	}

	resp, err := cln.Embedding(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("ERROR: %w", err)
	}

	return &VectorizedChunk{
		Chunk:  text,
		Vector: resp.Data[0].Embedding,
	}, nil
}

func vectorDBSearch(image, query string) (*VectorizedChunks, error) {

	// Embed the query.
	chunk, err := embed(image, query)
	if err != nil {
		return nil, fmt.Errorf("ERROR: %w", err)
	}
	vectorStr := "["
	for idx, val := range chunk.Vector {
		vectorStr += fmt.Sprintf("%f", val)
		if idx < len(chunk.Vector)-1 {
			vectorStr += ", "
		}
	}
	vectorStr += "]"

	// Get the DB connection string from env var.
	connStr := os.Getenv("DB_CONN_STR")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("ERROR: %w", err)
	}
	defer db.Close()

	// Query the database for the nearest neighbors.
	rows, err := db.Query("SELECT id, chunk FROM items ORDER BY embedding <=> $1 LIMIT 5", vectorStr)
	if err != nil {
		return nil, fmt.Errorf("ERROR: %w", err)
	}
	defer rows.Close()

	var vectorizedChunks VectorizedChunks
	for rows.Next() {
		var id int
		var chunk string
		err := rows.Scan(&id, &chunk)
		if err != nil {
			return nil, fmt.Errorf("ERROR: %w", err)
		}
		vectorizedChunks = append(vectorizedChunks, VectorizedChunk{
			Id:    id,
			Chunk: chunk,
		})
	}

	return &vectorizedChunks, nil
}
