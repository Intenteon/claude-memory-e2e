package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestOllamaStringTools tests the string-based Ollama delegation tools.
func TestOllamaStringTools(t *testing.T) {
	requireFixture(t)
	requireOllama(t)

	t.Run("llm_review_code", func(t *testing.T) {
		result, err := mcp.CallTool("llm_review_code", map[string]interface{}{
			"code":       "func Add(a, b int) int { return a + b }",
			"focus":      "general code quality",
			"project_id": fixtureProjectID,
		})
		if skipIfUnavailable(t, err, result, "ollama") {
			return
		}
		require.NoError(t, err)
		assert.False(t, result.IsError, "tool error: %s", result.Text)
		assert.NotEmpty(t, result.Text)
	})

	t.Run("llm_explain_code", func(t *testing.T) {
		result, err := mcp.CallTool("llm_explain_code", map[string]interface{}{
			"code":       "func fibonacci(n int) int {\n  if n <= 1 { return n }\n  return fibonacci(n-1) + fibonacci(n-2)\n}",
			"project_id": fixtureProjectID,
		})
		if skipIfUnavailable(t, err, result, "ollama") {
			return
		}
		require.NoError(t, err)
		assert.False(t, result.IsError, "tool error: %s", result.Text)
		assert.NotEmpty(t, result.Text)
	})

	t.Run("llm_fix_code", func(t *testing.T) {
		result, err := mcp.CallTool("llm_fix_code", map[string]interface{}{
			"code":       "func divide(a, b int) int { return a / b }",
			"error":      "runtime error: integer divide by zero when b is 0",
			"project_id": fixtureProjectID,
		})
		if skipIfUnavailable(t, err, result, "ollama") {
			return
		}
		require.NoError(t, err)
		assert.False(t, result.IsError, "tool error: %s", result.Text)
		assert.NotEmpty(t, result.Text)
	})

	t.Run("llm_generate_code", func(t *testing.T) {
		result, err := mcp.CallTool("llm_generate_code", map[string]interface{}{
			"prompt":     "Write a function that reverses a string",
			"language":   "go",
			"project_id": fixtureProjectID,
		})
		if skipIfUnavailable(t, err, result, "ollama") {
			return
		}
		require.NoError(t, err)
		assert.False(t, result.IsError, "tool error: %s", result.Text)
		assert.NotEmpty(t, result.Text)
	})

	t.Run("llm_general_task", func(t *testing.T) {
		result, err := mcp.CallTool("llm_general_task", map[string]interface{}{
			"task":       "Explain the difference between a mutex and a channel in Go concurrency",
			"project_id": fixtureProjectID,
		})
		if skipIfUnavailable(t, err, result, "ollama") {
			return
		}
		require.NoError(t, err)
		assert.False(t, result.IsError, "tool error: %s", result.Text)
		assert.NotEmpty(t, result.Text)
	})
}

// TestOllamaFileTools tests the file-aware Ollama delegation tools.
func TestOllamaFileTools(t *testing.T) {
	requireFixture(t)
	requireOllama(t)

	t.Run("llm_review_file", func(t *testing.T) {
		result, err := mcp.CallTool("llm_review_file", map[string]interface{}{
			"file_path":  "user/service.go",
			"focus":      "bugs",
			"project_id": fixtureProjectID,
		})
		if skipIfUnavailable(t, err, result, "ollama file tools") {
			return
		}
		require.NoError(t, err)
		assert.False(t, result.IsError, "tool error: %s", result.Text)
		assert.NotEmpty(t, result.Text)
	})

	t.Run("llm_explain_file", func(t *testing.T) {
		result, err := mcp.CallTool("llm_explain_file", map[string]interface{}{
			"file_path":  "auth/middleware.go",
			"project_id": fixtureProjectID,
		})
		if skipIfUnavailable(t, err, result, "ollama file tools") {
			return
		}
		require.NoError(t, err)
		assert.False(t, result.IsError, "tool error: %s", result.Text)
		assert.NotEmpty(t, result.Text)
	})

	t.Run("llm_analyze_files", func(t *testing.T) {
		result, err := mcp.CallTool("llm_analyze_files", map[string]interface{}{
			"file_paths": []string{"user/service.go", "user/repository.go"},
			"task":       "find dependencies between these files",
			"project_id": fixtureProjectID,
		})
		if skipIfUnavailable(t, err, result, "ollama file tools") {
			return
		}
		require.NoError(t, err)
		assert.False(t, result.IsError, "tool error: %s", result.Text)
		assert.NotEmpty(t, result.Text)
	})

	t.Run("llm_generate_code_with_context", func(t *testing.T) {
		result, err := mcp.CallTool("llm_generate_code_with_context", map[string]interface{}{
			"prompt":        "Write a function that validates an email address",
			"language":      "go",
			"context_files": []string{"user/service.go"},
			"project_id":    fixtureProjectID,
		})
		if skipIfUnavailable(t, err, result, "ollama file tools") {
			return
		}
		require.NoError(t, err)
		assert.False(t, result.IsError, "tool error: %s", result.Text)
		assert.NotEmpty(t, result.Text)
	})
}

// TestOllamaWriteTests tests llm_write_tests (alias for generate_test).
func TestOllamaWriteTests(t *testing.T) {
	requireFixture(t)
	requireOllama(t)

	t.Run("write_tests_for_GetUser", func(t *testing.T) {
		result, err := mcp.CallTool("llm_write_tests", map[string]interface{}{
			"function_name": "GetUser",
			"file_path":     "user/service.go",
			"project_id":    fixtureProjectID,
		})
		if skipIfUnavailable(t, err, result, "ollama write tests") {
			return
		}
		require.NoError(t, err)
		assert.False(t, result.IsError, "tool error: %s", result.Text)
		assert.NotEmpty(t, result.Text)
	})
}

// TestOllamaRefactorCode tests llm_refactor_code (alias for suggest_refactoring).
func TestOllamaRefactorCode(t *testing.T) {
	requireFixture(t)
	requireOllama(t)

	t.Run("refactor_CreateUser", func(t *testing.T) {
		result, err := mcp.CallTool("llm_refactor_code", map[string]interface{}{
			"function_name": "CreateUser",
			"file_path":     "user/service.go",
			"project_id":    fixtureProjectID,
		})
		if skipIfUnavailable(t, err, result, "ollama refactor") {
			return
		}
		require.NoError(t, err)
		assert.False(t, result.IsError, "tool error: %s", result.Text)
		assert.NotEmpty(t, result.Text)
	})
}
