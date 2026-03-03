package tests

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearchCode(t *testing.T) {
	requireFixture(t)

	t.Run("CreateUser_found", func(t *testing.T) {
		result, err := mcp.CallTool("search_code", map[string]interface{}{
			"query":      "CreateUser",
			"project_id": fixtureProjectID,
		})
		require.NoError(t, err)
		assert.False(t, result.IsError, "tool error: %s", result.Text)
		assert.Contains(t, result.Text, "CreateUser")
	})

	t.Run("place_order_returns_PlaceOrder", func(t *testing.T) {
		result, err := mcp.CallTool("search_code", map[string]interface{}{
			"query":      "place order",
			"project_id": fixtureProjectID,
		})
		require.NoError(t, err)
		assert.False(t, result.IsError)
		assert.Contains(t, result.Text, "PlaceOrder")
	})

	t.Run("validate_token_returns_ValidateToken", func(t *testing.T) {
		result, err := mcp.CallTool("search_code", map[string]interface{}{
			"query":      "validate token authentication",
			"project_id": fixtureProjectID,
		})
		require.NoError(t, err)
		assert.False(t, result.IsError)
		assert.Contains(t, result.Text, "ValidateToken")
	})

	t.Run("with_language_filter", func(t *testing.T) {
		result, err := mcp.CallTool("search_code", map[string]interface{}{
			"query":      "Save",
			"project_id": fixtureProjectID,
			"language":   "go",
		})
		require.NoError(t, err)
		assert.False(t, result.IsError)
		assert.Contains(t, result.Text, "Save")
	})

	t.Run("with_limit", func(t *testing.T) {
		result, err := mcp.CallTool("search_code", map[string]interface{}{
			"query":      "user",
			"project_id": fixtureProjectID,
			"limit":      3,
		})
		require.NoError(t, err)
		assert.False(t, result.IsError)
		assert.NotEmpty(t, result.Text)
	})
}

func TestFuzzySearchFunctions(t *testing.T) {
	requireFixture(t)

	t.Run("Validat_matches_ValidateToken", func(t *testing.T) {
		result, err := mcp.CallTool("fuzzy_search_functions", map[string]interface{}{
			"query":        "Validat",
			"project_id":   fixtureProjectID,
			"threshold":    0.1,
			"max_distance": 10,
		})
		require.NoError(t, err)
		assert.False(t, result.IsError, "tool error: %s", result.Text)
		assert.Contains(t, result.Text, "ValidateToken")
	})

	t.Run("CreateUse_matches_CreateUser", func(t *testing.T) {
		result, err := mcp.CallTool("fuzzy_search_functions", map[string]interface{}{
			"query":        "CreateUse",
			"project_id":   fixtureProjectID,
			"threshold":    0.1,
			"max_distance": 10,
		})
		require.NoError(t, err)
		assert.False(t, result.IsError)
		assert.Contains(t, result.Text, "CreateUser")
	})

	t.Run("PlaceOrde_matches_PlaceOrder", func(t *testing.T) {
		result, err := mcp.CallTool("fuzzy_search_functions", map[string]interface{}{
			"query":        "PlaceOrde",
			"project_id":   fixtureProjectID,
			"threshold":    0.1,
			"max_distance": 10,
		})
		require.NoError(t, err)
		assert.False(t, result.IsError)
		assert.Contains(t, result.Text, "PlaceOrder")
	})
}

func TestFindSimilarCode(t *testing.T) {
	requireFixture(t)

	t.Run("find_user_creation_code", func(t *testing.T) {
		result, err := mcp.CallTool("find_similar_code", map[string]interface{}{
			"query":      "function that creates a new user with name and email",
			"project_id": fixtureProjectID,
		})
		require.NoError(t, err)
		if result.IsError && strings.Contains(result.Text, "not available") {
			t.Skip("find_similar_code requires embedding service")
		}
		assert.NotEmpty(t, result.Text)
	})
}
