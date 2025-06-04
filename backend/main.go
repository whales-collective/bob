package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"we-are-legion/agents"

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

type AgentConfig struct {
	Name               string                                 `json:"name"`
	Description        string                                 `json:"description"`
	Agent              *robby.Agent                           `json:"agent"`
}

func main() {

	// create a map of agents
	agentsCatalog := map[string]*AgentConfig{
		"bob": {
			Name:        "Bob",
			Description: "The original Bob agent",
			Agent: func() *robby.Agent {
				agent, err := agents.GetBob()
				if err != nil {
					log.Fatal("😡 Error creating Bob agent: ", err)
				}
				agent.Params.Messages = []openai.ChatCompletionMessageParamUnion{
					openai.SystemMessage(`
					Your name is Bob,
					You are the original Bob agent,
					You are a helpful assistant.
					If the user asks something about your clones, you can display this list of clones:
					- 😎 Bob: yourself
					- 🤓 Bill: to be define later 🚧
					- 🙂 Milo: to be define later 🚧
					- 🐱 Garfield: to be define later 🚧
					- 🤖 Riker: is in charge of the invocation of the other clones of Bob.
					`),
				}
				return agent
			}(),
		},
		"bill": {
			Name:        "Bill",
			Description: "A clone of Bob, with a different personality",
			Agent: func() *robby.Agent {
				agent, err := agents.GetBill()
				if err != nil {
					log.Fatal("😡 Error creating Bill agent: ", err)
				}
				agent.Params.Messages = []openai.ChatCompletionMessageParamUnion{
					openai.SystemMessage(`					
					Your name is Bill,
					You are a clone of Bob,
					You are a helpful assistant, but you have a different personality than Bob.
					`),
				}
				return agent
			}(),
		},
		"milo": {
			Name:        "Milo",
			Description: "A clone of Bob, with a different personality",
			Agent: func() *robby.Agent {
				agent, err := agents.GetMilo()
				if err != nil {
					log.Fatal("😡 Error creating Milo agent: ", err)
				}
				agent.Params.Messages = []openai.ChatCompletionMessageParamUnion{
					openai.SystemMessage(`
					Your name is Milo,
					You are a clone of Bob,
					You are a helpful assistant, but you have a different personality than Bob.
					`),
				}
				return agent
			}(),
		},
		"garfield": {
			Name:        "Garfield",
			Description: "A clone of Bob, with a different personality",
			Agent: func() *robby.Agent {
				agent, err := agents.GetGarfield()
				if err != nil {
					log.Fatal("😡 Error creating Garfield agent: ", err)
				}
				agent.Params.Messages = []openai.ChatCompletionMessageParamUnion{
					openai.SystemMessage(`
					Your name is Garfield,
					You are a clone of Bob,
					You are a helpful assistant, but you have a different personality than Bob.
					`),
				}
				return agent
			}(),
		},
	}

	selectedAgentName := "bob" // Default agent name
	currentSelection := agentsCatalog[selectedAgentName]

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
		//response.Write([]byte("<step>Analysing your message...</step>"))
		flusher.Flush()

		riker.Params.Messages = []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(data["message"]),
		}

		// Étape 2: Vérification des outils
		response.Write([]byte("<info>Checking for tool calls...</info>"))
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
			response.Write([]byte("<orange>Executing tool calls...</orange>"))
			flusher.Flush()

			toolCallsJSON, _ := riker.ToolCallsToJSON()
			fmt.Println("✋ Tool Calls:\n", toolCallsJSON)

			// This method execute the tool calls detected by the Agent.
			// And add the result to the message list of the Agent.
			results, err := riker.ExecuteToolCalls(map[string]func(any) (any, error){
				// TODO: remove this tool
				"add": func(args any) (any, error) {
					response.Write([]byte("<feature>Performing addition...</feature>"))
					flusher.Flush()
					a := args.(map[string]any)["a"].(float64)
					b := args.(map[string]any)["b"].(float64)
					return a + b, nil
				},

				"choose_clone_of_bob": func(args any) (any, error) {
					response.Write([]byte("<yellow>Selecting Bob clone...</yellow>"))
					flusher.Flush()
					cloneName := args.(map[string]any)["clone_name"].(string)

					//

					switch strings.ToLower(cloneName) {
					case "bill":
						response.Write([]byte("<enhancement>Hey, it's Bill!</enhancement>"))
						flusher.Flush()
						selectedAgentName := "bill"
						currentSelection = agentsCatalog[selectedAgentName]
						return "Bill is a clone of Bob", nil

					case "milo":
						response.Write([]byte("<enhancement>Hey, it's Milo!</enhancement>"))
						flusher.Flush()
						selectedAgentName := "milo"
						currentSelection = agentsCatalog[selectedAgentName]
						return "Milo is a clone of Bob", nil

					case "garfield":
						response.Write([]byte("<enhancement>Hey, it's Garfield!</enhancement>"))
						flusher.Flush()
						selectedAgentName := "garfield"
						currentSelection = agentsCatalog[selectedAgentName]
						return "Garfield is a clone of Bob", nil

					case "bob":
						response.Write([]byte("<enhancement>Hey, it's Bob!</enhancement>"))
						flusher.Flush()
						selectedAgentName := "bob"
						currentSelection = agentsCatalog[selectedAgentName]
						return "Bob is the start of everything", nil

					default:
						response.Write([]byte("<bug>" + cloneName + " is unknown!</bug>"))
						flusher.Flush()
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

			// NOTE: conversational memory
			currentSelection.Agent.Params.Messages = append(currentSelection.Agent.Params.Messages,
				openai.SystemMessage(strings.Join(results, " ")), // QUESTION: tool message or agent message?
				openai.SystemMessage("use the above result of the tool calls to answer the user question: "),
				openai.UserMessage(data["message"]),
			)

			/*
			currentSelection.Agent.Params.Messages = []openai.ChatCompletionMessageParamUnion{
				currentSelection.SystemInstructions,
				openai.SystemMessage(strings.Join(results, " ")),
				openai.SystemMessage("use the above result of the tool calls to answer the user question: "),
				openai.UserMessage(data["message"]),
			}
			*/

			riker.Params.Messages = []openai.ChatCompletionMessageParamUnion{}

		} else {
			fmt.Println("📝 user message:", data["message"])

			// NOTE: conversational memory, add the question to the Agent's message
			currentSelection.Agent.Params.Messages = append(
				currentSelection.Agent.Params.Messages, openai.UserMessage(data["message"]),
			)
			/*
				currentSelection.Agent.Params.Messages = []openai.ChatCompletionMessageParamUnion{
					currentSelection.SystemInstructions,
					openai.UserMessage(data["message"]),
				}
			*/

			fmt.Println("📝 number of messages:", len(currentSelection.Agent.Params.Messages))
		}

		// Étape finale: Génération de la réponse
		response.Write([]byte("<step>Generating response...</step><br>"))
		flusher.Flush()

		answer, errCompletion := currentSelection.Agent.ChatCompletionStream(func(self *robby.Agent, content string, err error) error {
			//fmt.Print(content)
			response.Write([]byte(content))

			flusher.Flush()
			if !shouldIStopTheCompletion {
				return nil
			} else {
				return errors.New("🚫 Cancelling request")
			}
		})
		if errCompletion != nil {
			// TODO: handle error
		}
		// NOTE: conversational memory, add the answer to the Agent's message
		currentSelection.Agent.Params.Messages = append(
			currentSelection.Agent.Params.Messages, openai.AssistantMessage(answer),
		)

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
