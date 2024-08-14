package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
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

	//messageContent := strings.Replace(board+"\n\n"+pgn+"\n\n"+"Chess expert move: ", "\n", "\\n", -1)
	//board = strings.Replace(board, "\n", "\\n", -1)
	messageContent := "Current placement of non-captured pieces on the board:\n" + board
	messageContent += "\n\nHistory of moves (PGN format):\n" + pgn
	if len(invalid) > 0 {
		messageContent += "\n\nInvalid moves (Do NOT respond with one of the following listed invalid moves.): " + strings.Join(invalid, ", ")
		messageContent += "\n\nNext chess move (different from the invalid moves): "
	} else {
		messageContent += "\n\nNext black chess move: "
	}

	fmt.Println(messageContent)

	input := client.ChatInput{
		Model: client.Models.Hermes2ProLlama38B,
		Messages: []client.ChatInputMessage{
			{
				Role:    client.Roles.System,
				Content: "You are an expert chess player. Given information about a chess game, you respond with a next chess move. You are playing the black chess pieces at the top of the board. Respond in standard Algebraic notation with only a single chess move and no other text.",
			},
			{
				Role:    client.Roles.User,
				Content: messageContent,
			},
		},
		MaxTokens:   200,
		Temperature: 0.1,
		TopP:        0.1,
		TopK:        50.0,
	}

	resp, err := cln.Chat(ctx, input)
	if err != nil {
		return "", fmt.Errorf("ERROR: %w", err)
	}

	move := resp.Choices[0].Message.Content
	fmt.Println(move)

	// if "..." in move, split on this substring.
	if strings.Contains(move, "...") {
		move = strings.Split(move, "...")[1]
		move = strings.TrimPrefix(move, " ")
	}

	// if "." in move, split on this substring.
	if strings.Contains(move, ".") {
		move = strings.Split(move, ".")[1]
		move = strings.TrimPrefix(move, " ")
	}

	// if " " in move, split on this substring.
	if strings.Contains(move, " ") {
		move = strings.Split(move, " ")[len(strings.Split(move, " "))-1]
	}

	return move, nil
}

func generateGameDescWithLLM(game string) (string, error) {
	cln := client.New(llmLogger, host, apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// TODO: Call the LLM to generate a description of the current game for use in searching
	// a knowledge base of reference game info.

	return resp.Choices[0].Message.Content, nil
}

func generateQAWithLLM(content, game, question string) (string, error) {
	cln := client.New(llmLogger, host, apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// TODO: Call the LLM using the content, game, and question strings to generate
	// helpful advice for the user.

	return resp.Choices[0].Message.Content, nil
}
