package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSecurityScan(t *testing.T) {
	requireFixture(t)

	t.Run("scan_middleware_file", func(t *testing.T) {
		result, err := mcp.CallTool("security_scan", map[string]interface{}{
			"file_path":  "auth/middleware.go",
			"project_id": fixtureProjectID,
		})
		if skipIfUnavailable(t, err, result, "security scanner") {
			return
		}
		require.NoError(t, err)
		assert.False(t, result.IsError, "tool error: %s", result.Text)
		assert.NotEmpty(t, result.Text)
	})
}

func TestGetVulnerabilities(t *testing.T) {
	requireFixture(t)

	t.Run("list_vulnerabilities_no_error", func(t *testing.T) {
		result, err := mcp.CallTool("get_vulnerabilities", map[string]interface{}{
			"project_id": fixtureProjectID,
		})
		if skipIfUnavailable(t, err, result, "security store") {
			return
		}
		require.NoError(t, err)
		assert.False(t, result.IsError, "tool error: %s", result.Text)
	})
}
