# Go GenAI App Build - Exercise 3

The goal of this exercise is to generate a chess move to play against the user. Add your code to the `TODO` secion in `handlers.go` and `llm.go` to accomplish this task. Note, feel free to "augment" the LLM call with info about the current game, notation information, etc. via additional functions, multiple LLM calls, utilization of the chess package, etc.

Execute the following steps:

1. Add your functionality to the `TODO` in `handlers.go` and `llm.go`
2. Build and run the updated API
3. Test the API health check with Postman (or similar) or via cURL:

```
curl --location 'http://localhost:8080/move' \
--header 'Content-Type: application/json' \
--data '{
    "game": "1. Nf3 d5 2. c4 Nc6 3. cxd5 Nf6 4. d4 Bg4 5. Bg5 e6 6. Bxf6 Qxd5 7. e4 Bb4+ 8. Ke2 Qd7 9. Nc3 O-O 10. a3 Be7 11. e5 Qd8 12. Ke1 *"
}'
```

Which should return something like:

```json
{
    "move": "Qd7",
    "game_original": "1. Nf3 d5 2. c4 Nc6 3. cxd5 Nf6 4. d4 Bg4 5. Bg5 e6 6. Bxf6 Qxd5 7. e4 Bb4+ 8. Ke2 Qd7 9. Nc3 O-O 10. a3 Be7 11. e5 Qd8 12. Ke1 *",
    "game_updated": "1. Nf3 d5 2. c4 Nc6 3. cxd5 Nf6 4. d4 Bg4 5. Bg5 e6 6. Bxf6 Qxd5 7. e4 Bb4+ 8. Ke2 Qd7 9. Nc3 O-O 10. a3 Be7 11. e5 Qd8 12. Ke1 Qd7  *"
}
```