package workflow

import (
	"fmt"
	"net/http"
	"strings"
	"we-are-legion/agents"
	"we-are-legion/helpers"

	"github.com/openai/openai-go"
	"github.com/sea-monkeys/robby"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func ExecuteToolCalls(response http.ResponseWriter, flusher http.Flusher,agentsCatalog map[string]*agents.AgentConfig, riker *robby.Agent, selectedAgent *agents.AgentConfig) ([]string, error) {
	helpers.ResponseLabel(response, flusher, "orange", "Executing tool calls...")

	// IMPORTANT: 
	// the job of Riker is only to detect if the user wants to change the current Agent,
	// and to execute the tool calls.
	// This method execute the tool calls detected by the Agent.
	// And add the result to the message list of the Agent.
	// BEGIN: execute the tool calls
	results, err := riker.ExecuteToolCalls(map[string]func(any) (any, error){

		"choose_clone_of_bob": func(args any) (any, error) {

			helpers.ResponseLabel(response, flusher, "yellow", "Selecting Bob clone...")
			cloneName := args.(map[string]any)["clone_name"].(string)
			cloneName = strings.ToLower(cloneName)

			switch cloneName {
			case "bill", "milo", "garfield", "bob":
				// NOTE: change the current selection to the selected clone
				selectedAgent = agentsCatalog[cloneName]

				caser := cases.Title(language.English)
				cloneName = caser.String(cloneName)

				txtLabel := "Hey, it's " + cloneName + ", " + selectedAgent.Agent.Params.Model
				helpers.ResponseLabel(response, flusher, "enhancement", txtLabel)

				// NOTE: conversational memory
				selectedAgent.Agent.Params.Messages = append(selectedAgent.Agent.Params.Messages,
					// IMPORTANT: QUESTION: should I use a system message or a agent message?
					openai.SystemMessage("You have been selected to speak with the user, your name is: "+selectedAgent.Name),
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
					selectedAgent = agentsCatalog["bob"]
					helpers.ResponseLabel(response, flusher, "pink", "You are speaking with Bob")
				case "docker compose":
					selectedAgent = agentsCatalog["bill"]
					helpers.ResponseLabel(response, flusher, "orange", "You are speaking with Bill")
				case "docker model runner":
					selectedAgent = agentsCatalog["garfield"]
					helpers.ResponseLabel(response, flusher, "red", "You are speaking with Garfield")
				case "docker bake":
					selectedAgent = agentsCatalog["milo"]
					helpers.ResponseLabel(response, flusher, "warning", "You are speaking with Milo")

				}
			}

			selectedAgent.Agent.Params.Messages = append(selectedAgent.Agent.Params.Messages,
				// IMPORTANT: QUESTION: should I use a system message or a agent message?
				openai.AssistantMessage("I understand that you want to talk about: "+topic),
				//openai.SystemMessage("You have been selected to speak with the user, your name is: "+currentSelection.Name),
				//openai.SystemMessage("use the above result of the tool calls to answer the user question: "),
				//openai.UserMessage(userQuestion),
			)

			fmt.Println("🤖 Detected topic in user message:", topic)
			return topic, nil
		},
	}) // END: execute the tool calls

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

	// NOTE: reset the Riker messages
	riker.Params.Messages = []openai.ChatCompletionMessageParamUnion{}

	return results, err
}
