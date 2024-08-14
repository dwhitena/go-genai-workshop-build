# Go GenAI App Build - Exercise 2

The goal of this exercise is to parse a natural language request for a chess move with the LLM. Add your code to the `TODO` secion in `llm.go` to accomplish this task. Note, feel free to "augment" the LLM call with info about the current game, notation information, etc. via additional functions, multiple LLM calls, utilization of the chess package, etc.

Execute the following steps:

1. Add your functionality to the `TODO` in `llm.go`
2. Build and run the updated API
3. Test the API health check with Postman (or similar) or via cURL:

```
curl --location 'http://localhost:8080/parse' \
--header 'Content-Type: application/json' \
--data '{
    "game": "1. g4 g6 2. h4 Nc6  *",
    "move": "Knight to h3"
}'
```

Which should return something like:

```json
{
    "move": "Nh3",
    "game_original": "1. g4 g6 2. h4 Nc6  *",
    "game_updated": "1. g4 g6 2. h4 Nc6 3. Nh3 *"
}
```