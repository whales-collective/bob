package agents

import (
	"context"
	"fmt"
	"os"
	"we-are-legion/rag"

	"github.com/openai/openai-go"
	"github.com/sea-monkeys/robby"
)

func GetBill() (*robby.Agent, error) {
	// TODO: handle error
	modelRunnerURL := os.Getenv("DMR_BASE_URL") + "/engines/llama.cpp/v1"
	model := os.Getenv("MODEL_RUNNER_CHAT_MODEL_BILL")
	embeddingModel := os.Getenv("MODEL_RUNNER_EMBEDDING_MODEL")

	chunks, err := rag.GetChunksOfCloneDocuments("bill")
	if err != nil {
		return nil, fmt.Errorf("error getting chunks for Bill: %w", err)
	}

	fmt.Println("🌍", modelRunnerURL)
	fmt.Println("📕 Bill, chat model:", model)
	fmt.Println("📗 Bill, embedding model:", embeddingModel)

	bill, err := robby.NewAgent(
		robby.WithDMRClient(
			context.Background(),
			modelRunnerURL,
		),
		robby.WithParams(
			openai.ChatCompletionNewParams{
				Model:       model,
				Messages:    []openai.ChatCompletionMessageParamUnion{},
				Temperature: openai.Opt(0.9),
			},
		),
		robby.WithEmbeddingParams(
			openai.EmbeddingNewParams{
				Model: embeddingModel,
			},
		),
		robby.WithRAGMemory(chunks),
	)
	if err != nil {
		return nil, err
	}
	return bill, nil

}

func InitializeBillAgent() (*AgentConfig, error) {
	bill, err := GetBill()
	if err != nil {
		return nil, fmt.Errorf("error creating Bill agent: %w", err)
	}
	bill.Params.Messages = []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(`
		Your name is Bill, you are a Docker Compose expert.
		You are a clone of Bob,
		You are a helpful assistant, but you have a different personality than Bob.

		If the user asks something about Docker Compose, do your best to answer it using only your knowledge.
		If the user asks something about Docker, you can use the Bob clone to answer it.
		If the user asks something about Docker Model Runner, you can use the Garfield clone to answer it.
		If the user asks something about Docker Bake, you can use the Milo clone to answer it.
		`),
	}

	return &AgentConfig{
		Name:        "Bill",
		Description: "A clone of Bob, with a different personality",
		Agent:       bill,
	}, nil
}
