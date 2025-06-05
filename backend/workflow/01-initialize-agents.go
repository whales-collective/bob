package workflow

import "we-are-legion/agents"

func InitializeAgents() (map[string]*agents.AgentConfig) {
	// create a map of agents
	agentsCatalog := map[string]*agents.AgentConfig{
		"bob": func() *agents.AgentConfig {
			cfg, err := agents.InitializeBobAgent()
			if err != nil {
				panic("Error initializing Bob agent: " + err.Error())
			}
			return cfg
		}(),
		"bill": func() *agents.AgentConfig {
			cfg, err := agents.InitializeBillAgent()
			if err != nil {
				panic("Error initializing Bill agent: " + err.Error())
			}
			return cfg
		}(),
		"milo": func() *agents.AgentConfig {
			cfg, err := agents.InitializeMiloAgent()
			if err != nil {
				panic("Error initializing Milo agent: " + err.Error())
			}
			return cfg
		}(),
		"garfield": func() *agents.AgentConfig {
			cfg, err := agents.InitializeGarfieldAgent()
			if err != nil {
				panic("Error initializing Garfield agent: " + err.Error())
			}
			return cfg
		}(),
		"riker": func() *agents.AgentConfig {
			cfg, err := agents.InitializeRikerAgent()
			if err != nil {
				panic("Error initializing Riker agent: " + err.Error())
			}
			return cfg
		}(),
		"khan": func() *agents.AgentConfig {
			cfg, err := agents.InitializeKhanAgent()
			if err != nil {
				panic("Error initializing Khan agent: " + err.Error())
			}
			return cfg
		}(),
	}
	return agentsCatalog

}
