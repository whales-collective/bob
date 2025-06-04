package agents

import (
	"context"
	"fmt"
	"os"

	"github.com/openai/openai-go"
	"github.com/sea-monkeys/robby"
)



func GetBill() (*robby.Agent, error) {
	// TODO: handle error
	modelRunnerURL := os.Getenv("DMR_BASE_URL") + "/engines/llama.cpp/v1"
	model := os.Getenv("MODEL_RUNNER_CHAT_MODEL_BILL")

	fmt.Println("🌍", modelRunnerURL)
	fmt.Println("📕 Bill, chat model:", model)

	bill, err := robby.NewAgent(
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
	return bill, nil

}
