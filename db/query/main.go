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

var host = "https://api.predictionguard.com"
var apiKey = os.Getenv("PREDICTIONGUARD_API_KEY")

// VectorizedChunk is a struct that holds a vectorized chunk.
type VectorizedChunk struct {
	Id       int       `json:"id"`
	Chunk    string    `json:"chunk"`
	Vector   []float64 `json:"vector"`
	Metadata string    `json:"metadata"`
}

// VectorizedChunks is a slice of vectorized chunks.
type VectorizedChunks []VectorizedChunk

func embed(imageLink string, text string) (*VectorizedChunk, error) {

	logger := func(ctx context.Context, msg string, v ...any) {
		s := fmt.Sprintf("msg: %s", msg)
		log.Println(s)
	}

	cln := client.New(logger, host, apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var image client.ImageNetwork
	if imageLink != "" {
		imageParsed, err := client.NewImageNetwork(imageLink)
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
	if imageLink != "" {
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

func main() {

	query := "What is unique about an end game without pawns?"
	image := "https://www.ragchess.com/wp-content/uploads/2020/08/word-image-261.png"

	// Embed the query.
	chunk, err := embed(image, query)
	if err != nil {
		log.Fatalf("ERROR: %v", err)
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
		log.Fatal(err)
	}
	defer db.Close()

	// Query the database for the nearest neighbors.
	rows, err := db.Query("SELECT id, chunk FROM items ORDER BY embedding <=> $1 LIMIT 5", vectorStr)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	i := 0
	for rows.Next() {
		var id int
		var chunk string
		err := rows.Scan(&id, &chunk)
		if err != nil {
			log.Fatal(err)
		}
		if i == 0 {
			fmt.Printf("id=%d, chunk=%s\n\n", id, chunk)
		}
		i++
	}
}
