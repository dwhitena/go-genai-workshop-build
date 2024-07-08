# Go Generative AI Workshop (App Build) Materials

This repo includes the materials for Go GenAI workshops (to follow on from an introductory workshop in `dwhitena/go-genai-workshop`). The materials were prepared by Daniel Whitenack for live attendees. However, others might benefit from them as they build generative AI applications with Go. 

## Application build description

This example application pulls together LLM chaining, retrieval/augmentation, embedding, prompting, etc. concepts to create an application that allows the user to:
- Play chess with an LLM, by putting in natural language move requests
- Generate LLM chess moves
- Get advice about their current strategy and the tactics employed in the game

## Structure

- [db](db) - Scripts to prep a database with reference chess information
- [api](api) - The backend REST API supporting the main functionality
- [ui](ui) - A thin UI that calls the REST API and displays the games