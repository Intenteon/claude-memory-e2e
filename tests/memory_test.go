package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRecordMemory(t *testing.T) {
	requireFixture(t)

	t.Run("record_and_verify", func(t *testing.T) {
		result, err := mcp.CallTool("record_memory", map[string]interface{}{
			"content":    "E2E test memory: UserService uses UserRepository for persistence",
			"type":       "learning",
			"project_id": fixtureProjectID,
			"metadata": map[string]interface{}{
				"file_paths": []string{"user/service.go", "user/repository.go"},
			},
		})
		require.NoError(t, err)
		assert.False(t, result.IsError, "tool error: %s", result.Text)
		assert.NotEmpty(t, result.Text)
	})
}

func TestGetRequirements(t *testing.T) {
	requireFixture(t)

	t.Run("query_returns_no_error", func(t *testing.T) {
		result, err := mcp.CallTool("get_requirements", map[string]interface{}{
			"project_id": fixtureProjectID,
		})
		require.NoError(t, err)
		// May return empty list but shouldn't error
		assert.False(t, result.IsError, "tool error: %s", result.Text)
	})
}
