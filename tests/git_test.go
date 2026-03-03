package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetRecentChanges(t *testing.T) {
	requireFixture(t)

	t.Run("recent_changes_last_30_days", func(t *testing.T) {
		result, err := mcp.CallTool("get_recent_changes", map[string]interface{}{
			"project_id": fixtureProjectID,
			"days":       30,
		})
		if skipIfUnavailable(t, err, result, "git context service") {
			return
		}
		require.NoError(t, err)
		assert.False(t, result.IsError, "tool error: %s", result.Text)
		assert.NotEmpty(t, result.Text)
	})
}

func TestGetFileHistory(t *testing.T) {
	requireFixture(t)

	t.Run("history_for_user_service", func(t *testing.T) {
		result, err := mcp.CallTool("get_file_history", map[string]interface{}{
			"file_path":  "user/service.go",
			"project_id": fixtureProjectID,
			"limit":      5,
		})
		if skipIfUnavailable(t, err, result, "git context service") {
			return
		}
		require.NoError(t, err)
		assert.False(t, result.IsError, "tool error: %s", result.Text)
		assert.NotEmpty(t, result.Text)
	})
}

func TestGetBlameContext(t *testing.T) {
	requireFixture(t)

	t.Run("blame_for_user_service_lines", func(t *testing.T) {
		result, err := mcp.CallTool("get_blame_context", map[string]interface{}{
			"file_path":  "user/service.go",
			"project_id": fixtureProjectID,
			"line_start": 1,
			"line_end":   10,
		})
		if skipIfUnavailable(t, err, result, "git context service") {
			return
		}
		require.NoError(t, err)
		assert.False(t, result.IsError, "tool error: %s", result.Text)
		assert.NotEmpty(t, result.Text)
	})
}

func TestGetChunkBlame(t *testing.T) {
	requireFixture(t)

	t.Run("chunk_blame_for_chunk_1", func(t *testing.T) {
		result, err := mcp.CallTool("get_chunk_blame", map[string]interface{}{
			"chunk_id":   1,
			"project_id": fixtureProjectID,
		})
		if skipIfUnavailable(t, err, result, "git context service") {
			return
		}
		require.NoError(t, err)
		assert.False(t, result.IsError, "tool error: %s", result.Text)
		assert.NotEmpty(t, result.Text)
	})
}

func TestGenerateCommit(t *testing.T) {
	requireFixture(t)

	t.Run("generate_commit_no_error", func(t *testing.T) {
		result, err := mcp.CallTool("generate_commit", map[string]interface{}{
			"project_id": fixtureProjectID,
		})
		if skipIfUnavailable(t, err, result, "smart commit service") {
			return
		}
		require.NoError(t, err)
		// May error if no staged changes — that's expected
		assert.NotEmpty(t, result.Text)
	})
}

func TestGetHotspots(t *testing.T) {
	requireFixture(t)

	t.Run("list_hotspots", func(t *testing.T) {
		result, err := mcp.CallTool("get_hotspots", map[string]interface{}{
			"project_id": fixtureProjectID,
			"limit":      5,
		})
		if skipIfUnavailable(t, err, result, "behavioral tools") {
			return
		}
		require.NoError(t, err)
		assert.False(t, result.IsError, "tool error: %s", result.Text)
	})
}

func TestGetCodeHealth(t *testing.T) {
	requireFixture(t)

	t.Run("code_health_for_user_service", func(t *testing.T) {
		result, err := mcp.CallTool("get_code_health", map[string]interface{}{
			"file_path":  "user/service.go",
			"project_id": fixtureProjectID,
		})
		if skipIfUnavailable(t, err, result, "behavioral tools") {
			return
		}
		require.NoError(t, err)
		assert.False(t, result.IsError, "tool error: %s", result.Text)
	})
}

func TestGetChangeCouplings(t *testing.T) {
	requireFixture(t)

	t.Run("couplings_for_user_service", func(t *testing.T) {
		result, err := mcp.CallTool("get_change_couplings", map[string]interface{}{
			"file_path":  "user/service.go",
			"project_id": fixtureProjectID,
		})
		if skipIfUnavailable(t, err, result, "behavioral tools") {
			return
		}
		require.NoError(t, err)
		assert.False(t, result.IsError, "tool error: %s", result.Text)
	})
}

func TestGetChurnTrend(t *testing.T) {
	requireFixture(t)

	t.Run("churn_trend_for_user_service", func(t *testing.T) {
		result, err := mcp.CallTool("get_churn_trend", map[string]interface{}{
			"file_path":  "user/service.go",
			"project_id": fixtureProjectID,
		})
		if skipIfUnavailable(t, err, result, "behavioral tools") {
			return
		}
		require.NoError(t, err)
		assert.False(t, result.IsError, "tool error: %s", result.Text)
	})
}
