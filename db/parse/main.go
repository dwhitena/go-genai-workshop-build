package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"strings"
)

type Chunk struct {
	Text  string `json:"text"`
	Image string `json:"image"`
}

type Chunks []Chunk

func main() {

	// Read in the lines of the chess.txt book.
	file, err := os.Open("chess.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Read the contents of the file
	content, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	// Loop over the lines in the HTML.
	lines := strings.Split(string(content), "\n")
	image := ""
	textChunk := ""
	var chunks Chunks
	for _, line := range lines {

		// If "Diag." is in the line, reset the image.
		lineCheck := strings.Replace(line, " ", "", -1)
		if len(lineCheck) > 5 && lineCheck[0:5] == "Diag." {
			chunks = append(chunks, Chunk{Text: textChunk, Image: image})
			line = strings.Replace(line, " ", "", -1)
			image = strings.Split(line, ".")[1]
			image = strings.Replace(image, "\r", "", -1)
			image = strings.Replace(image, "#", "", -1)
			if len(image) == 1 {
				image = "0" + image
			}
			image = "https://www.gutenberg.org/cache/epub/5614/images/diag" + image + ".jpg"
			textChunk = ""
			continue
		}

		// Append the line to the HTML chunk.
		textChunk += line + "\n"
	}

	// Marshall the chunks into JSON file.
	json, err := json.MarshalIndent(chunks, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	// Write the JSON file.
	err = os.WriteFile("chunks.json", json, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
