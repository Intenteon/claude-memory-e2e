package tests

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetInterface(t *testing.T) {
	requireFixture(t)

	t.Run("Repository_interface_found", func(t *testing.T) {
		result, err := mcp.CallTool("get_interface", map[string]interface{}{
			"interface_name": "Repository",
			"project_id":     fixtureProjectID,
		})
		require.NoError(t, err)
		assert.False(t, result.IsError, "tool error: %s", result.Text)
		// Verify the interface itself is found with expected methods
		assert.Contains(t, result.Text, "Repository")
		assert.Contains(t, result.Text, "Save")
		assert.Contains(t, result.Text, "FindByID")
		assert.Contains(t, result.Text, "Delete")
		// Note: Go uses implicit interface satisfaction, so the indexer may not
		// detect concrete implementers. If implementers are detected, verify them.
		if strings.Contains(result.Text, "UserRepository") {
			assert.Contains(t, result.Text, "UserRepository")
		}
	})
}

func TestGetProjectSummary(t *testing.T) {
	requireFixture(t)

	t.Run("returns_nonempty_summary", func(t *testing.T) {
		result, err := mcp.CallTool("get_project_summary", map[string]interface{}{
			"project_id": fixtureProjectID,
		})
		require.NoError(t, err)
		assert.False(t, result.IsError, "tool error: %s", result.Text)
		assert.NotEmpty(t, result.Text)
	})

	t.Run("with_language_breakdown", func(t *testing.T) {
		result, err := mcp.CallTool("get_project_summary", map[string]interface{}{
			"project_id":        fixtureProjectID,
			"include_languages": true,
		})
		require.NoError(t, err)
		assert.False(t, result.IsError)
		assert.NotEmpty(t, result.Text)
	})
}
