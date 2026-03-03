package tests

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGoalLifecycle(t *testing.T) {
	requireFixture(t)

	// Step 1: Create a goal
	var goalID int64
	t.Run("set_goal", func(t *testing.T) {
		result, err := mcp.CallTool("set_goal", map[string]interface{}{
			"content":    "E2E test goal: verify goal lifecycle",
			"priority":   7,
			"project_id": fixtureProjectID,
		})
		require.NoError(t, err)
		assert.False(t, result.IsError, "tool error: %s", result.Text)
		assert.Contains(t, result.Text, "Goal created")

		goalID = extractGoalID(t, result.Text)
		require.NotZero(t, goalID, "failed to extract goal ID from: %s", result.Text)
	})

	if goalID == 0 {
		t.Fatal("cannot continue lifecycle test without a goal ID")
	}

	// Step 2: Verify it appears in active goals
	t.Run("get_active_goals_includes_new_goal", func(t *testing.T) {
		result, err := mcp.CallTool("get_active_goals", map[string]interface{}{
			"project_id": fixtureProjectID,
		})
		require.NoError(t, err)
		assert.False(t, result.IsError)
		assert.Contains(t, result.Text, "E2E test goal")
	})

	// Step 3: Update progress
	t.Run("update_goal_progress", func(t *testing.T) {
		result, err := mcp.CallTool("update_goal_progress", map[string]interface{}{
			"goal_id":  goalID,
			"progress": 50,
		})
		require.NoError(t, err)
		assert.False(t, result.IsError)
		assert.Contains(t, result.Text, "50%")
	})

	// Step 4: Complete the goal
	t.Run("complete_goal", func(t *testing.T) {
		result, err := mcp.CallTool("complete_goal", map[string]interface{}{
			"goal_id": goalID,
		})
		require.NoError(t, err)
		assert.False(t, result.IsError)
		assert.Contains(t, result.Text, "completed")
	})

	// Step 5: Verify it's no longer active
	t.Run("completed_goal_not_in_active_list", func(t *testing.T) {
		result, err := mcp.CallTool("get_active_goals", map[string]interface{}{
			"project_id": fixtureProjectID,
		})
		require.NoError(t, err)
		assert.False(t, result.IsError)
		assert.NotContains(t, result.Text, "E2E test goal")
	})
}

// extractGoalID parses the goal ID from the set_goal response text.
// The response contains JSON like: {"id":123,"content":"..."}
func extractGoalID(t *testing.T, text string) int64 {
	t.Helper()

	// The response text is "Goal created successfully:\n{...json...}"
	// Find the JSON part
	idx := 0
	for i, ch := range text {
		if ch == '{' {
			idx = i
			break
		}
	}
	if idx == 0 && text[0] != '{' {
		return 0
	}

	var data struct {
		ID int64 `json:"id"`
	}
	if err := json.Unmarshal([]byte(text[idx:]), &data); err != nil {
		return 0
	}
	return data.ID
}
