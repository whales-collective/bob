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
	"we-are-legion/helpers"

	"github.com/openai/openai-go"
	"github.com/sea-monkeys/robby"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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

	// create a map of agents
	agentsCatalog := map[string]*agents.AgentConfig{
		"bob": func() *agents.AgentConfig {
			cfg, _ := agents.InitializeBobAgent()
			return cfg
		}(),
		"bill": func() *agents.AgentConfig {
			cfg, _ := agents.InitializeBillAgent()
			return cfg
		}(),
		"milo": func() *agents.AgentConfig {
			cfg, _ := agents.InitializeMiloAgent()
			return cfg
		}(),
		"garfield": func() *agents.AgentConfig {
			cfg, _ := agents.InitializeGarfieldAgent()
			return cfg
		}(),
	}

	currentSelection := agentsCatalog["bob"]

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
			//response.Write([]byte("😡 Error: expected http.ResponseWriter to be an http.Flusher"))
			helpers.ResponseLabel(response, flusher, "error", "Expected http.ResponseWriter to be an http.Flusher")
		}
		body := GetBytesBody(request)
		// unmarshal the json data
		var data map[string]string

		err := json.Unmarshal(body, &data)
		if err != nil {
			//response.Write([]byte("😡 Error: " + err.Error()))
			helpers.ResponseLabel(response, flusher, "error", "Error parsing JSON: "+err.Error())
		}
		fmt.Println("🤖 Bob is ready to chat!", data)
		userQuestion := data["message"]

		riker.Params.Messages = []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(userQuestion),
		}

		// STEP 1: check if there are tool calls to detect in the user message
		// This is done by the Riker agent, which is in charge of detecting tool calls.
		// It will return a list of tool calls to execute.
		helpers.ResponseLabel(response, flusher, "info", "Checking for tool calls...")

		// Always check the TOOLCALLS:
		toolCalls, err := riker.ToolsCompletion()
		if err != nil {
			if len(toolCalls) > 0 {
				fmt.Println("😡 Error: ", err.Error())
				helpers.ResponseLabel(response, flusher, "error", "Tool call error detected: "+err.Error())
			} else {
				fmt.Println("🙂 no tool calls detected.")
				helpers.ResponseLabel(response, flusher, "success", "No tool calls detected")
			}
		}
		fmt.Println("🤖 Number of Tool Calls:\n", len(toolCalls))

		// OPTION 1: if there are tool calls, execute them
		if len(toolCalls) > 0 {
			helpers.ResponseLabel(response, flusher, "orange", "Executing tool calls...")

			toolCallsJSON, _ := riker.ToolCallsToJSON()
			fmt.Println("🤖 Tool Calls:\n", toolCallsJSON)

			// IMPORTANT: the job of Riker is only to detect if the user wants to change the current Agent,
			// and to execute the tool calls.
			// This method execute the tool calls detected by the Agent.
			// And add the result to the message list of the Agent.
			results, err := riker.ExecuteToolCalls(map[string]func(any) (any, error){
				// TODO: remove this tool
				//"add": func(args any) (any, error) {
				//	response.Write([]byte("<feature>Performing addition...</feature>"))
				//	flusher.Flush()
				//	a := args.(map[string]any)["a"].(float64)
				//	b := args.(map[string]any)["b"].(float64)
				//	return a + b, nil
				//},

				"choose_clone_of_bob": func(args any) (any, error) {

					helpers.ResponseLabel(response, flusher, "yellow", "Selecting Bob clone...")
					cloneName := args.(map[string]any)["clone_name"].(string)
					cloneName = strings.ToLower(cloneName)

					switch cloneName {
					case "bill", "milo", "garfield", "bob":
						// NOTE: change the current selection to the selected clone
						currentSelection = agentsCatalog[cloneName]

						caser := cases.Title(language.English)
						cloneName = caser.String(cloneName)

						txtLabel := "Hey, it's " + cloneName + ", " + currentSelection.Agent.Params.Model
						helpers.ResponseLabel(response, flusher, "enhancement", txtLabel)

						// NOTE: conversational memory
						currentSelection.Agent.Params.Messages = append(currentSelection.Agent.Params.Messages,
							// IMPORTANT: QUESTION: should I use a system message or a agent message?
							openai.SystemMessage("You have been selected to speak with the user, your name is: "+currentSelection.Name),
							//openai.SystemMessage("use the above result of the tool calls to answer the user question: "),
							//openai.UserMessage(userQuestion),
						)

						return cloneName, nil

					default:
						helpers.ResponseLabel(response, flusher, "bug", "Unknown clone of Bob: "+cloneName)

						return fmt.Sprintf("Unknown clone of Bob: %s", cloneName), nil

					}

				},
				// IMPORTANT: TODO: check if it could be better to delegate this tool to another tool agent?
				"detect_the_real_topic_in_user_message": func(args any) (any, error) {
					helpers.ResponseLabel(response, flusher, "step", "Detecting the real topic in user message...")

					topic := args.(map[string]any)["topic_name"].(string)
					// Here you can implement your logic to detect the real topic in the user message
					// For now, we will just return the user message as the detected topic
					helpers.ResponseLabel(response, flusher, "white", "Topic: "+topic)

					switch topic {
					case "docker", "docker compose", "docker model runner", "docker bake":
						// NOTE: change the current selection to the selected clone
						switch topic {
						case "docker":
							currentSelection = agentsCatalog["bob"]
							helpers.ResponseLabel(response, flusher, "pink", "You are speaking with Bob")
						case "docker compose":
							currentSelection = agentsCatalog["bill"]
							helpers.ResponseLabel(response, flusher, "orange", "You are speaking with Bill")
						case "docker model runner":
							currentSelection = agentsCatalog["garfield"]
							helpers.ResponseLabel(response, flusher, "red", "You are speaking with Garfield")
						case "docker bake":
							currentSelection = agentsCatalog["milo"]
							helpers.ResponseLabel(response, flusher, "warning", "You are speaking with Milo")

						}
					}

					currentSelection.Agent.Params.Messages = append(currentSelection.Agent.Params.Messages,
						// IMPORTANT: QUESTION: should I use a system message or a agent message?
						openai.AssistantMessage("I understand that you want to talk about: "+topic),
						//openai.SystemMessage("You have been selected to speak with the user, your name is: "+currentSelection.Name),
						//openai.SystemMessage("use the above result of the tool calls to answer the user question: "),
						//openai.UserMessage(userQuestion),
					)

					fmt.Println("🤖 Detected topic in user message:", topic)
					return topic, nil
				},
			})

			if err != nil {
				helpers.ResponseLabel(response, flusher, "error", "Tool execution failed: "+err.Error())
			} else {
				helpers.ResponseLabel(response, flusher, "success", "Tool calls executed successfully")
			}

			fmt.Println("")

			// Print the results of the tool calls
			fmt.Println("🤖 Results of the tool calls execution:")
			for _, result := range results {
				fmt.Println(result)
			}

			// NOTE: conversational memory
			//currentSelection.Agent.Params.Messages = append(currentSelection.Agent.Params.Messages,
			//	openai.SystemMessage("You have been selected to speak with the user, your name is: "+currentSelection.Name),
			//	//openai.SystemMessage("use the above result of the tool calls to answer the user question: "),
			//	//openai.UserMessage(userQuestion),
			//)
			// NOTE: reset the Riker messages
			riker.Params.Messages = []openai.ChatCompletionMessageParamUnion{}

		} else {
			// OPTION 2: if there are no tool calls, just continue the conversation
			fmt.Println("🤖 No tool calls detected, continuing the conversation...")
			fmt.Println("📝 user message:", userQuestion)
		}

		// STEP 2: SIMILARITY SEARCH:
		similarities, err := currentSelection.Agent.RAGMemorySearchSimilaritiesWithText(userQuestion, 0.7)
		if err != nil {
			fmt.Println("Error when searching for similarities:", err)
			// NOTE: do nothing, just continue the conversation
		}
		fmt.Println("🎉 Similarities found:", len(similarities))
		for _, similarity := range similarities {
			fmt.Println("-", similarity)
		}

		if len(similarities) > 0 {
			// NOTE: conversational memory, add the similarities to the Agent's message
			currentSelection.Agent.Params.Messages = append(
				currentSelection.Agent.Params.Messages,
				openai.SystemMessage(
					"Here are some relevant documents found in the RAG memory:\n"+strings.Join(similarities, "\n"),
				),
				openai.SystemMessage("Use the above documents to answer the user question: "),
				openai.UserMessage(userQuestion),
			)
		} else {
			// NOTE: conversational memory, add the question to the Agent's message
			currentSelection.Agent.Params.Messages = append(
				currentSelection.Agent.Params.Messages, openai.UserMessage(userQuestion),
			)
		}

		fmt.Println("🤖 number of messages in memory:", len(currentSelection.Agent.Params.Messages))

		// STEP 3: generate the response using the selected Agent
		helpers.ResponseLabelNewLine(response, flusher, "info", "Generating response...")

		answer, errCompletion := currentSelection.Agent.ChatCompletionStream(func(self *robby.Agent, content string, err error) error {
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

		// 👋 TODO: IMPORTANT: make a toll calls detection to see if we need to change of agent

	})

	// Cancel/Stop the generation of the completion
	mux.HandleFunc("DELETE /cancel", func(response http.ResponseWriter, request *http.Request) {
		shouldIStopTheCompletion = true
		helpers.ResponseLabel(response, response.(http.Flusher), "info", "Cancelling request...")
	})

	var errListening error
	log.Println("🌍 http server is listening on: " + httpPort)
	errListening = http.ListenAndServe(":"+httpPort, mux)

	log.Fatal(errListening)

}
