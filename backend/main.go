package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/openai/openai-go"
	"github.com/sea-monkeys/robby"
)

/*
GetBytesBody returns the body of an HTTP request as a []byte.
  - It takes a pointer to an http.Request as a parameter.
  - It returns a []byte.
*/
func GetBytesBody(request *http.Request) []byte {
	body := make([]byte, request.ContentLength)
	request.Body.Read(body)
	return body
}

func main() {

	modelRunnerURL := os.Getenv("DMR_BASE_URL") + "/engines/llama.cpp/v1"
	model := os.Getenv("MODEL_RUNNER_MODEL")

	bob, _ := robby.NewAgent(
		robby.WithDMRClient(
			context.Background(),
			modelRunnerURL,
		),
		robby.WithParams(
			openai.ChatCompletionNewParams{
				Model: model,
				Messages: []openai.ChatCompletionMessageParamUnion{
					openai.SystemMessage("Your name is Bob, You are a pizza expert"),
					openai.UserMessage("What is the best pizza in the world?"),
				},
				Temperature: openai.Opt(0.9),
			},
		),
	)

	var httpPort = os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "5050"
	}

	fmt.Println("🌍", modelRunnerURL, "📕", model)

	mux := http.NewServeMux()
	shouldIStopTheCompletion := false

	mux.HandleFunc("POST /chat", func(response http.ResponseWriter, request *http.Request) {
		// add a flusher
		flusher, ok := response.(http.Flusher)
		if !ok {
			response.Write([]byte("😡 Error: expected http.ResponseWriter to be an http.Flusher"))
		}
		body := GetBytesBody(request)
		// unmarshal the json data
		var data map[string]string

		err := json.Unmarshal(body, &data)
		if err != nil {
			response.Write([]byte("😡 Error: " + err.Error()))
		}

		bob.ChatCompletionStream(func(self *robby.Agent, content string, err error) error {
			//fmt.Print(content)
			response.Write([]byte(content))

			flusher.Flush()
			if !shouldIStopTheCompletion {
				return nil
			} else {
				return errors.New("🚫 Cancelling request")
			}
			//return nil
		})

		if err != nil {
			shouldIStopTheCompletion = false
			response.Write([]byte("bye: " + err.Error()))
		}

	})

	// Cancel/Stop the generation of the completion
	mux.HandleFunc("DELETE /cancel", func(response http.ResponseWriter, request *http.Request) {
		shouldIStopTheCompletion = true
		response.Write([]byte("🚫 Cancelling request..."))
	})

	var errListening error
	log.Println("🌍 http server is listening on: " + httpPort)
	errListening = http.ListenAndServe(":"+httpPort, mux)

	log.Fatal(errListening)

}
