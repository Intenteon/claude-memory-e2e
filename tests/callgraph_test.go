package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCallers(t *testing.T) {
	requireFixture(t)

	t.Run("ValidateToken_called_by_Authenticate", func(t *testing.T) {
		result, err := mcp.CallTool("get_callers", map[string]interface{}{
			"function_name": "ValidateToken",
			"file_path":     "auth/middleware.go",
			"project_id":    fixtureProjectID,
		})
		require.NoError(t, err)
		assert.False(t, result.IsError, "tool error: %s", result.Text)
		assert.Contains(t, result.Text, "Authenticate")
	})

	t.Run("GetUser_called_by_ValidateToken", func(t *testing.T) {
		result, err := mcp.CallTool("get_callers", map[string]interface{}{
			"function_name": "GetUser",
			"file_path":     "user/service.go",
			"project_id":    fixtureProjectID,
		})
		require.NoError(t, err)
		assert.False(t, result.IsError, "tool error: %s", result.Text)
		assert.Contains(t, result.Text, "ValidateToken")
	})

	t.Run("Save_called_by_CreateUser", func(t *testing.T) {
		result, err := mcp.CallTool("get_callers", map[string]interface{}{
			"function_name": "Save",
			"file_path":     "user/repository.go",
			"project_id":    fixtureProjectID,
		})
		require.NoError(t, err)
		assert.False(t, result.IsError, "tool error: %s", result.Text)
		assert.Contains(t, result.Text, "CreateUser")
	})

	t.Run("with_depth_2", func(t *testing.T) {
		result, err := mcp.CallTool("get_callers", map[string]interface{}{
			"function_name": "GetUser",
			"file_path":     "user/service.go",
			"project_id":    fixtureProjectID,
			"depth":         2,
		})
		require.NoError(t, err)
		assert.False(t, result.IsError)
		assert.NotEmpty(t, result.Text)
	})
}

func TestGetCallees(t *testing.T) {
	requireFixture(t)

	t.Run("CreateUser_calls_Save", func(t *testing.T) {
		result, err := mcp.CallTool("get_callees", map[string]interface{}{
			"function_name": "CreateUser",
			"file_path":     "user/service.go",
			"project_id":    fixtureProjectID,
		})
		require.NoError(t, err)
		assert.False(t, result.IsError, "tool error: %s", result.Text)
		assert.Contains(t, result.Text, "Save")
	})

	t.Run("ValidateToken_calls_GetUser", func(t *testing.T) {
		result, err := mcp.CallTool("get_callees", map[string]interface{}{
			"function_name": "ValidateToken",
			"file_path":     "auth/middleware.go",
			"project_id":    fixtureProjectID,
		})
		require.NoError(t, err)
		assert.False(t, result.IsError, "tool error: %s", result.Text)
		assert.Contains(t, result.Text, "GetUser")
	})

	t.Run("PlaceOrder_calls_GetUser_and_Save", func(t *testing.T) {
		result, err := mcp.CallTool("get_callees", map[string]interface{}{
			"function_name": "PlaceOrder",
			"file_path":     "order/service.go",
			"project_id":    fixtureProjectID,
		})
		require.NoError(t, err)
		assert.False(t, result.IsError, "tool error: %s", result.Text)
		assert.Contains(t, result.Text, "GetUser")
		assert.Contains(t, result.Text, "Save")
	})

	t.Run("CancelOrder_calls_FindByID_and_Delete", func(t *testing.T) {
		result, err := mcp.CallTool("get_callees", map[string]interface{}{
			"function_name": "CancelOrder",
			"file_path":     "order/service.go",
			"project_id":    fixtureProjectID,
		})
		require.NoError(t, err)
		assert.False(t, result.IsError, "tool error: %s", result.Text)
		assert.Contains(t, result.Text, "FindByID")
		assert.Contains(t, result.Text, "Delete")
	})
}

func TestGetCompactCallers(t *testing.T) {
	requireFixture(t)

	t.Run("compact_callers_for_GetUser", func(t *testing.T) {
		result, err := mcp.CallTool("get_compact_callers", map[string]interface{}{
			"function_name": "GetUser",
			"file_path":     "user/service.go",
			"project_id":    fixtureProjectID,
			"top_n":         5,
		})
		require.NoError(t, err)
		assert.False(t, result.IsError, "tool error: %s", result.Text)
		assert.NotEmpty(t, result.Text)
	})
}
