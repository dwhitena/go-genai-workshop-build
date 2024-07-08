package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/predictionguard/go-client"
)

var host = "https://api.predictionguard.com"
var apiKey = os.Getenv("PREDICTIONGUARD_API_KEY")

// characterTextSplitter takes in a string and splits the string into
// chunks of a given size (split on whitespace) with an overlap of a
// given size of tokens (split on whitespace).
func characterTextSplitter(text string, splitSize int, overlapSize int) []string {

	// Create a slice to hold the chunks.
	chunks := []string{}

	// Split the text into tokens based on whitespace.
	tokens := strings.Split(text, " ")

	// Loop over the tokens creating chunks of size splitSize with an
	// overlap of overlapSize.
	for i := 0; i < len(tokens); i += splitSize - overlapSize {
		end := i + splitSize - overlapSize
		if end > len(tokens) {
			end = len(tokens)
		}
		chunks = append(chunks, strings.Join(tokens[i:end], " "))
	}
	return chunks
}

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

type Chunk struct {
	Text  string `json:"text"`
	Image string `json:"image"`
}

type Chunks []Chunk

func main() {

	// Read in the chunks from ../parse/chunks.json.
	file, err := os.Open("../parse/chunks.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Read the contents of the file into a value of Chunks.
	var chunks Chunks
	err = json.NewDecoder(file).Decode(&chunks)
	if err != nil {
		log.Fatal(err)
	}

	// Loop over the chunks and embed them.
	var vectorizedChunks VectorizedChunks
	chunkId := 0
	for idx, chunk := range chunks {

		fmt.Printf("Embedding chunk %d of %d\n", idx+1, len(chunks))

		// Use the characterTextSplitter to split the chunk into smaller bits.
		sectionChunks := characterTextSplitter(chunk.Text, 500, 50)

		// Loop over the section chunks and embed them.
		for _, sectionChunk := range sectionChunks {
			vectorizedChunk, err := embed(chunk.Image, sectionChunk)
			if err != nil {
				log.Fatal(err)
			}

			// Add the vectorized chunk to the vectorizedChunks slice.
			vectorizedChunk.Id = chunkId
			vectorizedChunk.Metadata = ""
			vectorizedChunks = append(vectorizedChunks, *vectorizedChunk)

			chunkId++
		}
	}

	// Output the vectorized chunks to a JSON file.
	file, err = os.Create("chunks_vectors.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Marshall the chunks into JSON file.
	json, err := json.MarshalIndent(vectorizedChunks, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	// Write the JSON file.
	err = os.WriteFile("chunks_vectors.json", json, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
