package main

import (
	"we-are-legion/agents"
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

func main() {

	bob, err := agents.GetBob()
	if err != nil {
		log.Fatal("😡 Error creating Bob agent: ", err)
	}

	// NOTE: we need a separate agent for the tool completion
	riker, err := agents.GetRiker()
	if err != nil {
		log.Fatal("😡 Error creating Riker agent: ", err)
	}

	var httpPort = os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "5050"
	}

	mux := http.NewServeMux()
	shouldIStopTheCompletion := false

	// Exemple de modifications dans le handler /chat de main.go

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

		// Étape 1: Analyse du message
		response.Write([]byte("<step>Analysing your message...</step>"))
		flusher.Flush()

		riker.Params.Messages = []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(data["message"]),
		}

		// Étape 2: Vérification des outils
		response.Write([]byte("<hr><info>Checking for tool calls...</info>"))
		flusher.Flush()

		// Always check the TOOLCALLS:
		toolCalls, err := riker.ToolsCompletion()
		if err != nil {
			if len(toolCalls) > 0 {
				fmt.Println("😡 Error: ", err.Error())
				response.Write([]byte("<error>Tool call error detected</error>"))
			} else {
				fmt.Println("🙂 no tool calls detected.")
				response.Write([]byte("<success>No tool calls detected</success>"))
			}
			flusher.Flush()
		}
		fmt.Println("✋ Number of Tool Calls:\n", len(toolCalls))

		if len(toolCalls) > 0 {
			response.Write([]byte("<step>Executing tool calls...</step>"))
			flusher.Flush()

			toolCallsJSON, _ := riker.ToolCallsToJSON()
			fmt.Println("✋ Tool Calls:\n", toolCallsJSON)

			// This method execute the tool calls detected by the Agent.
			// And add the result to the message list of the Agent.
			results, err := riker.ExecuteToolCalls(map[string]func(any) (any, error){
				"add": func(args any) (any, error) {
					response.Write([]byte("<info>Performing addition...</info>"))
					flusher.Flush()
					a := args.(map[string]any)["a"].(float64)
					b := args.(map[string]any)["b"].(float64)
					return a + b, nil
				},
				"choose_clone_of_bob": func(args any) (any, error) {
					response.Write([]byte("<info>Selecting Bob clone...</info>"))
					flusher.Flush()
					cloneName := args.(map[string]any)["clone_name"].(string)
					if cloneName == "Bill" {
						return "Bill is a clone of Bob", nil
					} else if cloneName == "Milo" {
						return "Milo is a clone of Bob", nil
					} else if cloneName == "Garfield" {
						return "Garfield is a clone of Bob", nil
					} else {
						return fmt.Sprintf("Unknown clone of Bob: %s", cloneName), nil
					}
				},
			})

			if err != nil {
				response.Write([]byte("<error>Tool execution failed: " + err.Error() + "</error>"))
				flusher.Flush()
			} else {
				response.Write([]byte("<success>Tool calls executed successfully</success>"))
				flusher.Flush()
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

		// Étape finale: Génération de la réponse
		response.Write([]byte("<hr><step>Generating response...</step><hr>"))
		flusher.Flush()

		bob.ChatCompletionStream(func(self *robby.Agent, content string, err error) error {
			//fmt.Print(content)
			response.Write([]byte(content))

			flusher.Flush()
			if !shouldIStopTheCompletion {
				return nil
			} else {
				return errors.New("🚫 Cancelling request")
			}
		})
	})
	/*
		if err != nil {
			shouldIStopTheCompletion = false
			response.Write([]byte("bye: " + err.Error()))
		}
	*/

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
