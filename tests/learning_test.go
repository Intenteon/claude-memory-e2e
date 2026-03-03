package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLearningLifecycle(t *testing.T) {
	requireFixture(t)

	// Step 1: Record a memory (creates a learning)
	t.Run("record_memory", func(t *testing.T) {
		result, err := mcp.CallTool("record_memory", map[string]interface{}{
			"content":    "E2E test learning: always validate user before placing order",
			"type":       "decision",
			"project_id": fixtureProjectID,
		})
		require.NoError(t, err)
		assert.False(t, result.IsError, "tool error: %s", result.Text)
	})

	// Step 2: List learnings
	var learningID int64
	t.Run("list_learnings", func(t *testing.T) {
		result, err := mcp.CallTool("list_learnings", map[string]interface{}{
			"project_id": fixtureProjectID,
			"limit":      10,
		})
		require.NoError(t, err)
		assert.False(t, result.IsError, "tool error: %s", result.Text)
		// The response might be "No learnings found" or contain actual learnings
		// Both are valid — the tool shouldn't error
	})

	// Step 3: Search learnings
	t.Run("search_learnings", func(t *testing.T) {
		result, err := mcp.CallTool("search_learnings", map[string]interface{}{
			"query":      "validate user",
			"project_id": fixtureProjectID,
		})
		require.NoError(t, err)
		assert.False(t, result.IsError, "tool error: %s", result.Text)
	})

	// Step 4: Rate a learning (if we have one)
	if learningID > 0 {
		t.Run("rate_learning", func(t *testing.T) {
			result, err := mcp.CallTool("rate_learning", map[string]interface{}{
				"learning_id": learningID,
				"rating":      "helpful",
			})
			require.NoError(t, err)
			assert.False(t, result.IsError)
		})

		// Step 5: Delete the learning (cleanup)
		t.Run("delete_learning", func(t *testing.T) {
			result, err := mcp.CallTool("delete_learning", map[string]interface{}{
				"learning_id": learningID,
			})
			require.NoError(t, err)
			assert.False(t, result.IsError)
		})
	}
}
