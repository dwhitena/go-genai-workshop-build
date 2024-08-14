# Go GenAI App Build - Exercise 1

The goal of this initial exercise is to setup the basic scaffolding of the backend API that will drive the Chess app. The code here includes only one healthcheck route on the endpoint `/`. 

Execute the following steps to ensure that you have the basics of a working API:

1. Build the API with `go build`
2. Execute `./exercise1`
3. Test the API health check with Postman (or similar) or via cURL:

```
curl --location 'http://localhost:8080'
```