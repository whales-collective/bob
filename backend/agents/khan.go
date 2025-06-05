package agents

import (
	"context"
	"fmt"
	"os"

	"github.com/openai/openai-go"
	"github.com/sea-monkeys/robby"
)

func GetKhan() (*robby.Agent, error) {
	// TODO: handle error
	modelRunnerURL := os.Getenv("DMR_BASE_URL") + "/engines/llama.cpp/v1"
	modelForTools := os.Getenv("MODEL_RUNNER_TOOLS_MODEL")

	fmt.Println("🌍", modelRunnerURL)
	fmt.Println("📘 Khan, tool model:", modelForTools)

	khan, err := robby.NewAgent(
		robby.WithDMRClient(
			context.Background(),
			modelRunnerURL,
		),
		robby.WithParams(
			openai.ChatCompletionNewParams{
				Model: modelForTools,
				Messages: []openai.ChatCompletionMessageParamUnion{
					openai.SystemMessage(`
					Your name is Khan, 
					Use the tool, only if the user specify he wants to use brae search.
					Otherwise, ignore the tool.
					`),
				},
				Temperature: openai.Opt(0.0),
				//ParallelToolCalls: openai.Bool(true),
			},
		),
		robby.WithMCPClient(robby.WithSocatMCPToolkit()),
		//robby.WithMCPClient(robby.WithDockerMCPToolkit()),
		robby.WithMCPTools([]string{"brave_web_search"}), 
		// NOTE: you must activate the fetch MCP server in Docker MCP Toolkit

	)
	if err != nil {
		return nil, err
	}
	return khan, nil
}

