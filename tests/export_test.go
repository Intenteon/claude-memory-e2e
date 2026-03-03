package tests

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExportGraph(t *testing.T) {
	requireFixture(t)

	t.Run("export_mermaid_format", func(t *testing.T) {
		result, err := mcp.CallTool("export_graph", map[string]interface{}{
			"project_id": fixtureProjectID,
			"format":     "mermaid",
			"depth":      2,
		})
		require.NoError(t, err)
		if result.IsError && strings.Contains(result.Text, "not configured") {
			t.Skip("graph export service not configured")
		}
		assert.False(t, result.IsError, "tool error: %s", result.Text)
		assert.NotEmpty(t, result.Text)
	})

	t.Run("export_focused_on_function", func(t *testing.T) {
		result, err := mcp.CallTool("export_graph", map[string]interface{}{
			"project_id":    fixtureProjectID,
			"format":        "mermaid",
			"function_name": "CreateUser",
			"direction":     "both",
			"depth":         2,
		})
		require.NoError(t, err)
		if result.IsError && strings.Contains(result.Text, "not configured") {
			t.Skip("graph export service not configured")
		}
		assert.False(t, result.IsError, "tool error: %s", result.Text)
		assert.NotEmpty(t, result.Text)
	})

	t.Run("export_dot_format", func(t *testing.T) {
		result, err := mcp.CallTool("export_graph", map[string]interface{}{
			"project_id": fixtureProjectID,
			"format":     "dot",
		})
		require.NoError(t, err)
		if result.IsError && strings.Contains(result.Text, "not configured") {
			t.Skip("graph export service not configured")
		}
		assert.False(t, result.IsError)
	})

	t.Run("export_json_format", func(t *testing.T) {
		result, err := mcp.CallTool("export_graph", map[string]interface{}{
			"project_id": fixtureProjectID,
			"format":     "json",
		})
		require.NoError(t, err)
		if result.IsError && strings.Contains(result.Text, "not configured") {
			t.Skip("graph export service not configured")
		}
		assert.False(t, result.IsError)
	})
}
