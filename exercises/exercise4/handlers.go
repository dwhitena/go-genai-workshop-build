package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image/color"
	"log"
	"net/http"
	"os"
	"os/exec"
	"slices"
	"strings"
	"time"

	"github.com/notnil/chess"
	"github.com/notnil/chess/image"
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
	//gameBoard := strings.TrimPrefix(game.Position().Board().Draw(), "\n")
	gameBoard := formatBoard(game)
	gamePGN := strings.TrimPrefix(game.String(), "\n")
	gamePGN = strings.TrimSuffix(gamePGN, " *")

	invalidMoves := []string{}
	var move string
	for i := 0; i < 10; i++ {

		// Generate a move with an LLM.
		move, err = generateMoveWithLLM(gameBoard, gamePGN, invalidMoves)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Move the piece.
		if err = game.MoveStr(move); err != nil {
			if i == 9 {
				fmt.Println(err.Error())
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// Add to invalid moves if the move isn't already in the list.
			if !slices.Contains(invalidMoves, move) {
				invalidMoves = append(invalidMoves, move)
			}
			continue
		}
		break
	}

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

type GenHelpRequest struct {
	Game string `json:"game"`
}

type GenHelpResponse struct {
	Message       string `json:"message"`
	ReferenceInfo string `json:"reference_info"`
}

// GenHelp generates help messages for a game.
func GenHelp(w http.ResponseWriter, r *http.Request) {

	// Parse the body into a value of MakeMoveRequest.
	var req GenHelpRequest
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

	// Create a temporary string based on the unix time.
	t := time.Now().Unix()
	tempFilename := fmt.Sprintf("/tmp/%d.svg", t)

	// Write the game to an SVG image.
	f, err := os.Create(tempFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	yellow := color.RGBA{255, 255, 0, 1}
	mark := image.MarkSquares(yellow, chess.D2, chess.D4)
	if err := image.SVG(f, game.Position().Board(), mark); err != nil {
		log.Fatal(err)
	}

	// Convert the SVG file to a JPG image using ffmpeg.
	cmd := exec.Command("ffmpeg", "-i", tempFilename, "-vf", "scale=800:-1", "-y", "/tmp/"+fmt.Sprintf("%d.jpg", t))
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	// Get a description of the game.
	description, err := generateGameDescWithLLM(req.Game)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Embed and search for relevant reference info.
	chunks, err := vectorDBSearch("/tmp/"+fmt.Sprintf("%d.jpg", t), description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Grab the text from the first chunk in chunks.
	var referenceInfo string
	for idx, chunk := range *chunks {
		if idx == 0 {
			referenceInfo = chunk.Chunk
		}
	}

	// Generate the response.
	responseMessage, err := generateQAWithLLM(description, req.Game, referenceInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Prep the response.
	resp := GenHelpResponse{
		Message:       responseMessage,
		ReferenceInfo: referenceInfo,
	}

	// Return the response.
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
