package workflow

import (
	"fmt"
	"strings"
	"we-are-legion/agents"

	"github.com/openai/openai-go"
)

func SearchSimilarities(selectedAgent *agents.AgentConfig, userQuestion string) []string {
	similarities, err := selectedAgent.Agent.RAGMemorySearchSimilaritiesWithText(userQuestion, 0.7)
	if err != nil {
		fmt.Println("Error when searching for similarities:", err)
		// NOTE: do nothing, just continue the conversation
	}
	fmt.Println("🎉 Similarities found:", len(similarities))
	//for _, similarity := range similarities {
	//	fmt.Println("-", similarity)
	//}
	if len(similarities) > 0 {
		// NOTE: conversational memory, add the similarities to the Agent's message
		selectedAgent.Agent.Params.Messages = append(
			selectedAgent.Agent.Params.Messages,
			openai.SystemMessage(
				"Here are some relevant documents found in the RAG memory:\n"+strings.Join(similarities, "\n"),
			),
			openai.SystemMessage("Use the above documents to answer the user question: "),
			openai.UserMessage(userQuestion),
		)
	} else {
		// NOTE: conversational memory, add the question to the Agent's message
		selectedAgent.Agent.Params.Messages = append(
			selectedAgent.Agent.Params.Messages, openai.UserMessage(userQuestion),
		)
	}
	return similarities
}
