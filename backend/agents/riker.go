package agents

import (
	"context"
	"fmt"
	"os"

	"github.com/openai/openai-go"
	"github.com/sea-monkeys/robby"
)

func GetRiker() (*robby.Agent, error) {
	// TODO: handle error
	modelRunnerURL := os.Getenv("DMR_BASE_URL") + "/engines/llama.cpp/v1"
	modelForTools := os.Getenv("MODEL_RUNNER_TOOLS_MODEL")

	fmt.Println("🌍", modelRunnerURL)
	fmt.Println("📘 Riker, tool model:", modelForTools)

	riker, err := robby.NewAgent(
		robby.WithDMRClient(
			context.Background(),
			modelRunnerURL,
		),
		robby.WithParams(
			openai.ChatCompletionNewParams{
				Model: modelForTools,
				Messages: []openai.ChatCompletionMessageParamUnion{
					openai.SystemMessage(`
					Your name is Riker, 
					You know how to join the other clones of Bob, 
					and you can use tools to do so.
					`),
				},
				Temperature: openai.Opt(0.0),
				//ParallelToolCalls: openai.Bool(true),
			},
		),
		robby.WithTools(GetRikerToolsCatalog()),
		//robby.WithMCPClient(robby.WithSocatMCPToolkit()),
		//robby.WithMCPClient(robby.WithDockerMCPToolkit()),
		//robby.WithMCPTools([]string{"brave_web_search"}), 
		// NOTE: you must activate the fetch MCP server in Docker MCP Toolkit


	)
	if err != nil {
		return nil, err
	}
	return riker, nil
}

func GetRikerToolsCatalog() []openai.ChatCompletionToolParam {
	/*
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
	*/

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

	detectTheRealTopicInUserMessage := openai.ChatCompletionToolParam{
		Function: openai.FunctionDefinitionParam{
			Name:        "detect_the_real_topic_in_user_message",
			Description: openai.String(`select a topic in this list [docker, docker compose, docker bake, docker model runner] by saying I have questions on <topic_name>.`),
			Parameters: openai.FunctionParameters{
				"type": "object",
				"properties": map[string]interface{}{	
					"topic_name": map[string]string{
						"type":        "string",
						"description": "The topic name to detect in the user message. The topic can be one of the following: [docker, docker compose, docker bake, docker model runner].",
					},
				},
				"required": []string{"message"},
			},
		},
	}

	tools := []openai.ChatCompletionToolParam{chooseCloneOfBobTool, detectTheRealTopicInUserMessage}
	return tools
}
