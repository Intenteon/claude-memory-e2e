package tests

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPlanLifecycle(t *testing.T) {
	requireFixture(t)

	var planID int64
	var stepID int64

	// Step 1: Create plan with steps
	t.Run("create_plan", func(t *testing.T) {
		result, err := mcp.CallTool("create_plan", map[string]interface{}{
			"title":       "E2E test plan",
			"description": "Verify plan lifecycle tools",
			"priority":    5,
			"project_id":  fixtureProjectID,
			"steps": []map[string]interface{}{
				{"step_number": 1, "title": "First step", "description": "Do the first thing"},
				{"step_number": 2, "title": "Second step", "description": "Do the second thing", "blocked_by": []int{1}},
			},
		})
		require.NoError(t, err)
		assert.False(t, result.IsError, "tool error: %s", result.Text)
		assert.Contains(t, result.Text, "Plan created")

		planID = extractPlanID(t, result.Text)
		require.NotZero(t, planID)
	})

	if planID == 0 {
		t.Fatal("cannot continue plan lifecycle without plan ID")
	}

	// Step 2: Verify it appears in active plans
	t.Run("get_active_plans", func(t *testing.T) {
		result, err := mcp.CallTool("get_active_plans", map[string]interface{}{
			"project_id": fixtureProjectID,
		})
		require.NoError(t, err)
		assert.False(t, result.IsError)
		assert.Contains(t, result.Text, "E2E test plan")
	})

	// Step 3: Add another step
	t.Run("add_plan_step", func(t *testing.T) {
		result, err := mcp.CallTool("add_plan_step", map[string]interface{}{
			"plan_id":     planID,
			"step_number": 3,
			"title":       "Third step",
			"description": "Do the third thing",
		})
		require.NoError(t, err)
		assert.False(t, result.IsError, "tool error: %s", result.Text)
		assert.Contains(t, result.Text, "added")

		stepID = extractStepID(t, result.Text)
	})

	// Step 4: Get next unblocked step (should be step 1)
	t.Run("get_next_unblocked_step", func(t *testing.T) {
		result, err := mcp.CallTool("get_next_unblocked_step", map[string]interface{}{
			"plan_id": planID,
		})
		require.NoError(t, err)
		assert.False(t, result.IsError)
		assert.Contains(t, result.Text, "First step")
	})

	// Step 5: Update step status
	if stepID != 0 {
		t.Run("update_plan_step_status", func(t *testing.T) {
			result, err := mcp.CallTool("update_plan_step_status", map[string]interface{}{
				"step_id": stepID,
				"status":  "completed",
			})
			require.NoError(t, err)
			assert.False(t, result.IsError)
			assert.Contains(t, result.Text, "completed")
		})
	}

	// Step 6: Complete the plan (cleanup)
	t.Run("complete_plan", func(t *testing.T) {
		result, err := mcp.CallTool("complete_plan", map[string]interface{}{
			"plan_id": planID,
		})
		require.NoError(t, err)
		assert.False(t, result.IsError)
		assert.Contains(t, result.Text, "completed")
	})
}

func extractPlanID(t *testing.T, text string) int64 {
	t.Helper()
	idx := 0
	for i, ch := range text {
		if ch == '{' {
			idx = i
			break
		}
	}
	if idx == 0 && len(text) > 0 && text[0] != '{' {
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

func extractStepID(t *testing.T, text string) int64 {
	t.Helper()
	// Response: "Step 3 added successfully to plan X (step ID: Y)"
	var stepID int64
	// Try to parse "step ID: N" from the text
	for i := 0; i < len(text)-8; i++ {
		if text[i:i+8] == "step ID:" {
			rest := text[i+8:]
			for _, ch := range rest {
				if ch >= '0' && ch <= '9' {
					stepID = stepID*10 + int64(ch-'0')
				} else if stepID > 0 {
					break
				}
			}
			break
		}
	}
	return stepID
}
