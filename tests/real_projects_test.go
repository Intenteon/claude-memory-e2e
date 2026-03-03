package tests

import (
	"testing"

	e2e "github.com/Intenteon/claude-memory-e2e"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRealProjectAuthService runs smoke tests against the auth-service project.
func TestRealProjectAuthService(t *testing.T) {
	requireRealProject(t, e2e.AuthServiceProjectID)

	t.Run("search_code", func(t *testing.T) {
		result, err := mcp.CallTool("search_code", map[string]interface{}{
			"query":      "authenticate",
			"project_id": e2e.AuthServiceProjectID,
		})
		require.NoError(t, err)
		assert.False(t, result.IsError, "tool error: %s", result.Text)
		assert.NotEmpty(t, result.Text)
	})

	t.Run("get_project_summary", func(t *testing.T) {
		result, err := mcp.CallTool("get_project_summary", map[string]interface{}{
			"project_id": e2e.AuthServiceProjectID,
		})
		require.NoError(t, err)
		assert.False(t, result.IsError)
		assert.NotEmpty(t, result.Text)
	})

	t.Run("get_active_goals", func(t *testing.T) {
		result, err := mcp.CallTool("get_active_goals", map[string]interface{}{
			"project_id": e2e.AuthServiceProjectID,
		})
		require.NoError(t, err)
		assert.False(t, result.IsError)
	})

	t.Run("fuzzy_search", func(t *testing.T) {
		result, err := mcp.CallTool("fuzzy_search_functions", map[string]interface{}{
			"query":      "Login",
			"project_id": e2e.AuthServiceProjectID,
			"threshold":  0.3,
		})
		require.NoError(t, err)
		assert.False(t, result.IsError)
	})
}

// TestRealProjectAuthSDK runs smoke tests against the auth-sdk-go project.
func TestRealProjectAuthSDK(t *testing.T) {
	requireRealProject(t, e2e.AuthSDKProjectID)

	t.Run("search_code", func(t *testing.T) {
		result, err := mcp.CallTool("search_code", map[string]interface{}{
			"query":      "authenticate",
			"project_id": e2e.AuthSDKProjectID,
		})
		require.NoError(t, err)
		assert.False(t, result.IsError, "tool error: %s", result.Text)
		assert.NotEmpty(t, result.Text)
	})

	t.Run("get_project_summary", func(t *testing.T) {
		result, err := mcp.CallTool("get_project_summary", map[string]interface{}{
			"project_id": e2e.AuthSDKProjectID,
		})
		require.NoError(t, err)
		assert.False(t, result.IsError)
		assert.NotEmpty(t, result.Text)
	})

	t.Run("get_active_goals", func(t *testing.T) {
		result, err := mcp.CallTool("get_active_goals", map[string]interface{}{
			"project_id": e2e.AuthSDKProjectID,
		})
		require.NoError(t, err)
		assert.False(t, result.IsError)
	})
}
