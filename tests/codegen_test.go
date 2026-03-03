package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateStub(t *testing.T) {
	requireFixture(t)
	requireOllama(t)

	t.Run("generate_go_function_stub", func(t *testing.T) {
		result, err := mcp.CallTool("generate_stub", map[string]interface{}{
			"function_name": "ValidateEmail",
			"description":   "Validates that an email address has correct format",
			"language":      "go",
			"project_id":    fixtureProjectID,
			"parameters": []map[string]string{
				{"name": "email", "type": "string"},
			},
			"return_type": "error",
		})
		if skipIfUnavailable(t, err, result, "codegen service") {
			return
		}
		require.NoError(t, err)
		assert.False(t, result.IsError, "tool error: %s", result.Text)
		assert.NotEmpty(t, result.Text)
	})
}

func TestGenerateDocs(t *testing.T) {
	requireFixture(t)
	requireOllama(t)

	t.Run("generate_docs_for_CreateUser", func(t *testing.T) {
		result, err := mcp.CallTool("generate_docs", map[string]interface{}{
			"function_name": "CreateUser",
			"file_path":     "user/service.go",
			"project_id":    fixtureProjectID,
		})
		if skipIfUnavailable(t, err, result, "codegen service") {
			return
		}
		require.NoError(t, err)
		assert.False(t, result.IsError, "tool error: %s", result.Text)
		assert.NotEmpty(t, result.Text)
	})
}

func TestGenerateTest(t *testing.T) {
	requireFixture(t)
	requireOllama(t)

	t.Run("generate_tests_for_CreateUser", func(t *testing.T) {
		result, err := mcp.CallTool("generate_test", map[string]interface{}{
			"function_name": "CreateUser",
			"file_path":     "user/service.go",
			"project_id":    fixtureProjectID,
		})
		if skipIfUnavailable(t, err, result, "codegen service") {
			return
		}
		require.NoError(t, err)
		assert.False(t, result.IsError, "tool error: %s", result.Text)
		assert.NotEmpty(t, result.Text)
	})
}

func TestSuggestRefactoring(t *testing.T) {
	requireFixture(t)
	requireOllama(t)

	t.Run("suggest_refactoring_for_PlaceOrder", func(t *testing.T) {
		result, err := mcp.CallTool("suggest_refactoring", map[string]interface{}{
			"function_name": "PlaceOrder",
			"file_path":     "order/service.go",
			"project_id":    fixtureProjectID,
			"focus":         "readability",
		})
		if skipIfUnavailable(t, err, result, "codegen service") {
			return
		}
		require.NoError(t, err)
		assert.False(t, result.IsError, "tool error: %s", result.Text)
		assert.NotEmpty(t, result.Text)
	})
}
