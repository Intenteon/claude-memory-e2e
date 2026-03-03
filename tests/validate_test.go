package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateChanges(t *testing.T) {
	requireFixture(t)

	t.Run("validate_edit_to_service", func(t *testing.T) {
		result, err := mcp.CallTool("validate_changes", map[string]interface{}{
			"file_path":  "user/service.go",
			"project_id": fixtureProjectID,
			"changes": []map[string]interface{}{
				{
					"old_content": "func (s *UserService) CreateUser(name, email string) (*User, error) {",
					"new_content": "func (s *UserService) CreateUser(name, email, role string) (*User, error) {",
					"start_line":  18,
					"end_line":    18,
				},
			},
			"change_type": "edit",
		})
		require.NoError(t, err)
		assert.False(t, result.IsError, "tool error: %s", result.Text)
		assert.NotEmpty(t, result.Text)
	})
}
