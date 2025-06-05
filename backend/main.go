package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"we-are-legion/helpers"
	"we-are-legion/workflow"

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

	// Get a map of the agents
	agentsCatalog := workflow.InitializeAgents()
	// Select the current agent to use
	selectedAgent := agentsCatalog["bob"]

	// NOTE: we need separate agents for the tool completions
	// Riker is the agent in charge of detecting if the user wants to change the current Agent,
	// and to execute the tool calls.
	// Khan is the agent in charge of detecting if the user wants to use the MCP tools,
	// and to execute the MCP tool calls.
	riker := agentsCatalog["riker"].Agent
	khan := agentsCatalog["khan"].Agent

	var httpPort = os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "5050"
	}

	mux := http.NewServeMux()
	shouldIStopTheCompletion := false // TODO: implement a way to stop the completion cleanly

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
			helpers.ResponseLabel(response, flusher, "error", "Error parsing JSON: "+err.Error())
		}

		// NOTE: this is the message typed by the user
		userQuestion := data["message"]

		riker.Params.Messages = []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(userQuestion),
		}

		khan.Params.Messages = []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(userQuestion),
		}

		// STEP 1: TOOLCALLS: check if there are tool calls to detect in the user message
		// This is done by the Riker agent for the tool calls
		// and the Khan agent for the MCP tool calls.
		helpers.ResponseLabel(response, flusher, "info", "Checking for tool calls...")

		toolCalls, _ := workflow.DetectToolCalls(response, flusher, riker)
		mcpTooCalls, _ := workflow.DetectMCPToolCalls(response, flusher, khan)

		// STEP 2: MCP TOOLS EXECUTION if there are any MCP tool calls
		var mcpResults []string
		if len(mcpTooCalls) > 0 {
			mcpResults, _ = workflow.ExecuteMCPToolCalls(response, flusher, khan)
		}

		// STEP 3: TOOL CALLS EXECUTION
		if len(toolCalls) > 0 {
			workflow.ExecuteToolCalls(response, flusher, agentsCatalog, riker, selectedAgent)
		} else {
			// NOTE: If there are no tool calls, 
			// just continue the conversation
			fmt.Println("🤖 No tool calls detected, continuing the conversation...")
			fmt.Println("📝 user message:", userQuestion)
		}

		// STEP 4: add context to the prompt
		if len(mcpTooCalls) > 0 && len(mcpResults) > 0 { // OPTION 1: add the result of the MCP tool calls execution to the Agent's message
			selectedAgent.Agent.Params.Messages = append(
				selectedAgent.Agent.Params.Messages,
				openai.SystemMessage("Here are some relevant documents found in the MCP memory:\n"+strings.Join(mcpResults, "\n")),
				openai.SystemMessage("Use the above documents to answer the user question: "),
				openai.UserMessage(userQuestion),
			)
		} else { // OPTION 2: make similarity search
			workflow.SearchSimilarities(selectedAgent, userQuestion)
		}

		fmt.Println("🧠 number of messages in memory:", len(selectedAgent.Agent.Params.Messages))

		// STEP 5: generate the response using the selected Agent
		helpers.ResponseLabelNewLine(response, flusher, "info", "Generating response...")

		answer, errCompletion := selectedAgent.Agent.ChatCompletionStream(func(self *robby.Agent, content string, err error) error {
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
		selectedAgent.Agent.Params.Messages = append(
			selectedAgent.Agent.Params.Messages, openai.AssistantMessage(answer),
		)

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
