package agents

import "github.com/sea-monkeys/robby"

type AgentConfig struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Agent       *robby.Agent `json:"agent"`
}
