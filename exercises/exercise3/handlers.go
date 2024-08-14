package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/notnil/chess"
)

// Index is the handler for the root URL.
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "The API is healthy!\n")
}

type ParseMoveRequest struct {
	Game string `json:"game"`
	Move string `json:"move"`
}

type ParseMoveResponse struct {
	Move         string `json:"move"`
	GameOriginal string `json:"game_original"`
	GameUpdated  string `json:"game_updated"`
}

func formatBoard(game *chess.Game) string {

	pieceMap := map[string]string{
		"r": "rook",
		"n": "knight",
		"b": "bishop",
		"q": "queen",
		"k": "king",
		"p": "pawn",
	}

	// Loop over squares on the chess board.
	pieceList := ""
	for sq := chess.A1; sq <= chess.H8; sq++ {

		// If the square is empty, skip it.
		if game.Position().Board().Piece(sq).Type() == chess.NoPieceType {
			pieceList += fmt.Sprintf("%v - open\n", sq)
			continue
		}

		// For each square, print a single line.
		pieceList += fmt.Sprintf(
			"%v - %s %s\n", sq,
			game.Position().Board().Piece(sq).Color().Name(),
			pieceMap[game.Position().Board().Piece(sq).Type().String()],
		)
	}

	return pieceList
}

// ParseMove parses natural language moves.
func ParseMove(w http.ResponseWriter, r *http.Request) {

	// Parse the body into a value of ParseMoveRequest.
	var req ParseMoveRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Parse the game.
	pgnReader := bytes.NewReader([]byte(req.Game))
	pgn, err := chess.PGN(pgnReader)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	game := chess.NewGame(pgn)

	// Parse the move with an LLM.
	pieceList := formatBoard(game)
	move, err := parseMoveWithLLM(req.Move, pieceList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Move the piece.
	if err = game.MoveStr(move); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Prepare the response.
	resp := ParseMoveResponse{
		Move:         move,
		GameOriginal: req.Game,
		GameUpdated:  strings.TrimPrefix(game.String(), "\n"),
	}

	// Return the response.
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type MakeMoveRequest struct {
	Game string `json:"game"`
}

type MakeMoveResponse struct {
	Move         string `json:"move"`
	GameOriginal string `json:"game_original"`
	GameUpdated  string `json:"game_updated"`
}

// MakeMove take a game and uses an LLM to make a move.
func MakeMove(w http.ResponseWriter, r *http.Request) {

	// Parse the body into a value of MakeMoveRequest.
	var req MakeMoveRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Parse the game.
	pgnReader := bytes.NewReader([]byte(req.Game))
	pgn, err := chess.PGN(pgnReader)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	game := chess.NewGame(pgn)
	gameBoard := formatBoard(game)
	gamePGN := strings.TrimPrefix(game.String(), "\n")
	gamePGN = strings.TrimSuffix(gamePGN, " *")

	// TODO: Generate a move with an LLM. Chain LLM calls together, utilize the
	// gameBoard and/or gamePGN to augement the calls. This section should generate
	// a valid chess move in a string call "move".

	// Prep the response.
	resp := MakeMoveResponse{
		Move:         move,
		GameOriginal: req.Game,
		GameUpdated:  strings.TrimPrefix(game.String(), "\n"),
	}

	// Return the response.
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
