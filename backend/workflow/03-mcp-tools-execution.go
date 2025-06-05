package workflow

import (
	"net/http"
	"we-are-legion/helpers"

	"github.com/sea-monkeys/robby"
)

func ExecuteMCPToolCalls(response http.ResponseWriter, flusher http.Flusher, khan *robby.Agent) ([]string, error) {
	helpers.ResponseLabel(response, flusher, "orange", "Executing MCP tool calls...")
	mcpResults, err := khan.ExecuteMCPToolCalls()
	if err != nil {
		helpers.ResponseLabel(response, flusher, "error", "MCP Tool execution failed: "+err.Error())
	} else {
		helpers.ResponseLabel(response, flusher, "success", "MCP Tool calls executed successfully")
	}
	return mcpResults, err
}
