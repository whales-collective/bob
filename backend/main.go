package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

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

func GetBobToolsCatalog() []openai.ChatCompletionToolParam {
	addTool := openai.ChatCompletionToolParam{
		Function: openai.FunctionDefinitionParam{
			Name:        "add",
			Description: openai.String("add two numbers"),
			Parameters: openai.FunctionParameters{
				"type": "object",
				"properties": map[string]interface{}{
					"a": map[string]string{
						"type":        "number",
						"description": "The first number to add.",
					},
					"b": map[string]string{
						"type":        "number",
						"description": "The second number to add.",
					},
				},
				"required": []string{"a", "b"},
			},
		},
	}

	chooseCloneOfBobTool := openai.ChatCompletionToolParam{
		Function: openai.FunctionDefinitionParam{
			Name:        "choose_clone_of_bob",
			Description: openai.String("choose a clone of Bob by saying I want to speak to <clone_name>"),
			Parameters: openai.FunctionParameters{
				"type": "object",
				"properties": map[string]interface{}{
					"clone_name": map[string]string{
						"type":        "string",
						"description": "The name of the clone of Bob to choose.",
					},
				},
				"required": []string{"clone_name"},
			},
		},
	}


	tools := []openai.ChatCompletionToolParam{addTool, chooseCloneOfBobTool}
	return tools
}

func main() {

	modelRunnerURL := os.Getenv("DMR_BASE_URL") + "/engines/llama.cpp/v1"
	model := os.Getenv("MODEL_RUNNER_CHAT_MODEL")
	modelForTools := os.Getenv("MODEL_RUNNER_TOOLS_MODEL")

	bob, _ := robby.NewAgent(
		robby.WithDMRClient(
			context.Background(),
			modelRunnerURL,
		),
		robby.WithParams(
			openai.ChatCompletionNewParams{
				Model: model,
				Messages: []openai.ChatCompletionMessageParamUnion{
					openai.SystemMessage("Your name is Bob, You are a the original Bob"),
				},
				Temperature: openai.Opt(0.9),
			},
		),
	)

	// NOTE: we need a separate agent for the tool completion

	riker, _ := robby.NewAgent(
		robby.WithDMRClient(
			context.Background(),
			modelRunnerURL,
		),
		robby.WithParams(
			openai.ChatCompletionNewParams{
				Model:             modelForTools,
				Messages:          []openai.ChatCompletionMessageParamUnion{},
				Temperature:       openai.Opt(0.0),
				//ParallelToolCalls: openai.Bool(true),
			},
		),
		robby.WithTools(GetBobToolsCatalog()),
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
		fmt.Println("🤖 Bob is ready to chat!", data)

		riker.Params.Messages = []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(data["message"]),
		}

		// Always check the TOOLCALLS:
		toolCalls, err := riker.ToolsCompletion()
		if err != nil {
			//response.Write([]byte("😡 Error: " + err.Error()))
			if len(toolCalls) > 0 {
				fmt.Println("😡 Error: ", err.Error())
			} else {
				fmt.Println("🙂 no tool calls detected.")
			}
			
		}
		fmt.Println("✋ Number of Tool Calls:\n", len(toolCalls))

		if len(toolCalls) > 0 {

			toolCallsJSON, _ := riker.ToolCallsToJSON()
			fmt.Println("✋ Tool Calls:\n", toolCallsJSON)

			// This method execute the tool calls detected by the Agent.
			// And add the result to the message list of the Agent.
			results, err := riker.ExecuteToolCalls(map[string]func(any) (any, error){
				"add": func(args any) (any, error) {
					a := args.(map[string]any)["a"].(float64)
					b := args.(map[string]any)["b"].(float64)
					return a + b, nil
				},
				"choose_clone_of_bob": func(args any) (any, error) {
					cloneName := args.(map[string]any)["clone_name"].(string)
					if cloneName == "Riker" {
						return "Riker is a clone of Bob", nil
					} else if cloneName == "Milo" {
						return "Milo is a clone of Bob", nil
					} else {
						return fmt.Sprintf("Unknown clone of Bob: %s", cloneName), nil
					}
				},
			})

			if err != nil {
				response.Write([]byte("😡 Error: " + err.Error()))
			}

			fmt.Println("")

			// Print the results of the tool calls
			fmt.Println("✋ Results of the tool calls execution:")
			for _, result := range results {
				fmt.Println(result)
			}

			bob.Params.Messages = append(bob.Params.Messages,
				openai.SystemMessage(strings.Join(results, " ")),
				openai.SystemMessage("use the above result of the tool calls to answer the user question: "),
				openai.UserMessage(data["message"]),
			)

			riker.Params.Messages = []openai.ChatCompletionMessageParamUnion{}

		} else {
			fmt.Println("📝 user message:", data["message"])
			bob.Params.Messages = append(bob.Params.Messages,
				openai.UserMessage(data["message"]),
			)
			fmt.Println("📝 number of messages:", len(bob.Params.Messages))


		}
		response.Write([]byte("**🤖 Bob is thinking...**"))
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

		/*
		if err != nil {
			shouldIStopTheCompletion = false
			response.Write([]byte("bye: " + err.Error()))
		}
		*/

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
