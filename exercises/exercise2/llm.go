package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/predictionguard/go-client"
)

var llmLogger = func(ctx context.Context, msg string, v ...any) {
	s := fmt.Sprintf("msg: %s", msg)
	for i := 0; i < len(v); i = i + 2 {
		s = s + fmt.Sprintf(", %s: %v", v[i], v[i+1])
	}
	log.Println(s)
}

var host = "https://api.predictionguard.com"
var apiKey = os.Getenv("PREDICTIONGUARD_API_KEY")

func parseMoveWithLLM(moveRequest string) (string, error) {
	cln := client.New(llmLogger, host, apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// TODO: Add your code here to prompt the LLM. The goal is to give the moveRequest
	// to the LLM (along with any instructions) to parse the moce into standard
	// chess notation. The function should return the move in standard chess notation.
	// Feel free to "augment" the prompt with additional information about the game,
	// notation information, etc.

	return resp.Choices[0].Message.Content, nil
}
