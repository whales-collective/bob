package agents

import "github.com/sea-monkeys/robby"

type AgentConfig struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Agent       *robby.Agent `json:"agent"`
	ToolAgent bool		 `json:"tool_agent,omitempty"` // Indicates if the agent has a tool agent
}
