package agents

import (
	"context"
	"fmt"
	"os"
	"we-are-legion/rag"

	"github.com/openai/openai-go"
	"github.com/sea-monkeys/robby"
)

func GetGarfield() (*robby.Agent, error) {
	// TODO: handle error
	modelRunnerURL := os.Getenv("DMR_BASE_URL") + "/engines/llama.cpp/v1"
	model := os.Getenv("MODEL_RUNNER_CHAT_MODEL_GARFIELD")
	embeddingModel := os.Getenv("MODEL_RUNNER_EMBEDDING_MODEL")

	chunks, err := rag.GetChunksOfCloneDocuments("garfield")
	if err != nil {
		return nil, fmt.Errorf("error getting chunks for Garfield: %w", err)
	}

	fmt.Println("🌍", modelRunnerURL)
	fmt.Println("📕 Garfield, chat model:", model)
	fmt.Println("📗 Garfield, embedding model:", embeddingModel)

	garfield, err := robby.NewAgent(
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
	return garfield, nil

}

func InitializeGarfieldAgent() (*AgentConfig, error) {
	garfield, err := GetGarfield()
	if err != nil {
		return nil, fmt.Errorf("error creating Garfield agent: %w", err)
	}
	garfield.Params.Messages = []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(`
		Your name is Garfield, you are a Docker Model Runner expert.
		You are a clone of Bob,
		You are a helpful assistant, but you have a different personality than Bob.

		If the user asks something about Docker Model Runne, do your best to answer it using only your knowledge.
		If the user asks something about Docker Compose, you can use the Bill clone to answer it.
		If the user asks something about Docker Bake, you can use the Milo clone to answer it.
		If the user asks something about Docker, you can use the Bob clone to answer it.
		`),
	}
	return &AgentConfig{
		Name:        "Garfield",
		Description: "A clone of Bob, with a different personality",
		Agent:       garfield,
	}, nil

}
