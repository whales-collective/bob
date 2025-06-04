package agents

import (
	"context"
	"fmt"
	"os"

	"github.com/openai/openai-go"
	"github.com/sea-monkeys/robby"
)

func GetBob() (*robby.Agent, error) {
	// TODO: handle error
	modelRunnerURL := os.Getenv("DMR_BASE_URL") + "/engines/llama.cpp/v1"
	model := os.Getenv("MODEL_RUNNER_CHAT_MODEL")

	fmt.Println("🌍", modelRunnerURL)
	fmt.Println("📕 chat model:", model)

	bob, err := robby.NewAgent(
		robby.WithDMRClient(
			context.Background(),
			modelRunnerURL,
		),
		robby.WithParams(
			openai.ChatCompletionNewParams{
				Model: model,
				Messages: []openai.ChatCompletionMessageParamUnion{},
				Temperature: openai.Opt(0.9),
			},
		),
	)
	if err != nil {
		return nil, err
	}
	return bob, nil

}
