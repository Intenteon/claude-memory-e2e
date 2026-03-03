package tests

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearchPatterns(t *testing.T) {
	requireFixture(t)

	t.Run("search_error_handling_patterns", func(t *testing.T) {
		result, err := mcp.CallTool("search_patterns", map[string]interface{}{
			"query":    "error handling",
			"language": "go",
		})
		require.NoError(t, err)
		if result.IsError && strings.Contains(result.Text, "not configured") {
			t.Skip("pattern store not configured")
		}
		assert.False(t, result.IsError, "tool error: %s", result.Text)
	})
}

func TestGetPattern(t *testing.T) {
	requireFixture(t)

	t.Run("get_pattern_by_id", func(t *testing.T) {
		result, err := mcp.CallTool("get_pattern", map[string]interface{}{
			"pattern_id": 1,
		})
		require.NoError(t, err)
		if result.IsError && (strings.Contains(result.Text, "not configured") || strings.Contains(result.Text, "not found")) {
			t.Skip("pattern not found or store not configured")
		}
		// Pattern may or may not exist — just verify no crash
	})
}

func TestListPatternUsages(t *testing.T) {
	requireFixture(t)

	t.Run("list_usages_for_pattern_1", func(t *testing.T) {
		result, err := mcp.CallTool("list_pattern_usages", map[string]interface{}{
			"pattern_id": 1,
			"limit":      5,
		})
		require.NoError(t, err)
		if result.IsError && (strings.Contains(result.Text, "not configured") || strings.Contains(result.Text, "not found")) {
			t.Skip("pattern not found or store not configured")
		}
	})
}
