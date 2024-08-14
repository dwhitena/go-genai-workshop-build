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

func parseMoveWithLLM(moveRequest string, pieceList string) (string, error) {
	cln := client.New(llmLogger, host, apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	input := client.ChatInput{
		Model: client.Models.Hermes2ProLlama38B,
		Messages: []client.ChatInputMessage{
			{
				Role:    client.Roles.System,
				Content: "You are an chess game assistant. Given a request for a chess move, you parse that chess move into standard chess Algebraic Notation. Respond with only the Algebraic notation of the requested chess move and no other text.\n\n-- if parsing a move related to a pawn, do not use any abbreviation such as \"P\" or \"N\" for the pawn. Instead, respond with the square that it is moving to (c6, e4, a5, etc.) and no other text.\n-- When a piece makes a capture (or \"takes\" or \"kills\" another piece), an \"x\" is inserted immediately before the destination square. For example, Bxe5 (bishop captures the piece on e5). When a pawn makes a capture, the file from which the pawn departed is used to identify the pawn. For example, exd5 (pawn on the e-file captures the piece on d5).\n-- For moves with pieces other than pawns, the King is abbreviated to K, the Queen is abbreviated to Q, Rooks are abbreviated to R, Knights are abbreviated to N, Bishops are abbreviated to B.\n\nHere is the current placement of pieces on the board for reference:\n" + pieceList,
			},
			{
				Role:    client.Roles.User,
				Content: "Move by white: " + moveRequest,
			},
		},
		MaxTokens:   10,
		Temperature: 0.1,
	}

	resp, err := cln.Chat(ctx, input)
	if err != nil {
		return "", fmt.Errorf("ERROR: %w", err)
	}

	return resp.Choices[0].Message.Content, nil
}

func generateMoveWithLLM(board string, pgn string, invalid []string) (string, error) {
	cln := client.New(llmLogger, host, apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// TODO: Relevant LLM functionality to make a move.

	return move, nil
}
