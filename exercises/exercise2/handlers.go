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
	Game string `json:"game"` // PGN
	Move string `json:"move"`
}

type ParseMoveResponse struct {
	Move         string `json:"move"`
	GameOriginal string `json:"game_original"`
	GameUpdated  string `json:"game_updated"`
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
