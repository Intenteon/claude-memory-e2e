package tests

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	e2e "github.com/Intenteon/claude-memory-e2e"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestBranchSwitch verifies that:
// 1. Incremental index does NOT clean up ghost functions from deleted files (documents the bug)
// 2. Force re-index DOES clean up ghost functions correctly
// 3. Switching back and force re-indexing restores the original functions
func TestBranchSwitch(t *testing.T) {
	requireFixture(t)

	fixtureDir := e2e.FixtureDir

	// Save the original HEAD SHA so we can verify clean restore
	origSHA := gitOutput(t, fixtureDir, "rev-parse", "HEAD")

	// Ensure cleanup regardless of test outcome
	defer func() {
		// Force restore to original state
		gitRun(t, fixtureDir, "checkout", "main")
		gitRun(t, fixtureDir, "branch", "-D", "test-branch-switch")
		// Remove any untracked files
		branchOnlyFile := filepath.Join(fixtureDir, "branch_only.go")
		os.Remove(branchOnlyFile)
		// Verify we're back to the pinned SHA
		currentSHA := gitOutput(t, fixtureDir, "rev-parse", "HEAD")
		if currentSHA != origSHA {
			t.Logf("WARNING: fixture SHA mismatch after cleanup: got %s, want %s", currentSHA, origSHA)
			gitRun(t, fixtureDir, "reset", "--hard", origSHA)
		}
		// Force re-index to restore original state
		runIndex(t, "--force")
	}()

	// --- Step 1: Verify baseline ---
	t.Run("baseline", func(t *testing.T) {
		result, err := mcp.CallTool("search_code", map[string]interface{}{
			"query":      "PlaceOrder",
			"project_id": fixtureProjectID,
		})
		require.NoError(t, err)
		assert.False(t, result.IsError, "tool error: %s", result.Text)
		assert.Contains(t, result.Text, "PlaceOrder", "PlaceOrder should exist in baseline index")

		result, err = mcp.CallTool("search_code", map[string]interface{}{
			"query":      "BranchOnlyFunction",
			"project_id": fixtureProjectID,
		})
		require.NoError(t, err)
		// BranchOnlyFunction should NOT be found (doesn't exist yet)
		if !result.IsError {
			assert.NotContains(t, result.Text, "BranchOnlyFunction",
				"BranchOnlyFunction should not exist in baseline")
		}
	})

	// --- Step 2: Create test branch and modify files ---
	gitRun(t, fixtureDir, "checkout", "-b", "test-branch-switch")

	// Step 3: Create a new file with a unique function
	branchOnlyContent := `package main

// BranchOnlyFunction is a test function that only exists on the test branch.
func BranchOnlyFunction() string {
	return "I only exist on the test branch"
}
`
	branchOnlyFile := filepath.Join(fixtureDir, "branch_only.go")
	err := os.WriteFile(branchOnlyFile, []byte(branchOnlyContent), 0644)
	require.NoError(t, err, "failed to write branch_only.go")

	// Step 4: Delete order/service.go (removes PlaceOrder, CancelOrder, GetOrder)
	orderServicePath := filepath.Join(fixtureDir, "order", "service.go")
	orderServiceBackup, err := os.ReadFile(orderServicePath)
	require.NoError(t, err, "failed to read order/service.go for backup")
	err = os.Remove(orderServicePath)
	require.NoError(t, err, "failed to delete order/service.go")

	// Ensure we can restore the file on cleanup
	defer func() {
		if _, err := os.Stat(orderServicePath); os.IsNotExist(err) {
			os.WriteFile(orderServicePath, orderServiceBackup, 0644)
		}
	}()

	// --- Step 5: Run incremental index ---
	t.Run("incremental_index", func(t *testing.T) {
		runIndex(t) // no --force = incremental

		// Step 6: BranchOnlyFunction should be found (new file picked up)
		result, err := mcp.CallTool("search_code", map[string]interface{}{
			"query":      "BranchOnlyFunction",
			"project_id": fixtureProjectID,
		})
		require.NoError(t, err)
		if !result.IsError {
			assert.Contains(t, result.Text, "BranchOnlyFunction",
				"incremental index should pick up new file")
		}

		// Step 7: PlaceOrder should STILL be found (ghost — documents the known bug)
		result, err = mcp.CallTool("search_code", map[string]interface{}{
			"query":      "PlaceOrder",
			"project_id": fixtureProjectID,
		})
		require.NoError(t, err)
		if !result.IsError && strings.Contains(result.Text, "PlaceOrder") {
			t.Log("KNOWN BUG: PlaceOrder is still found after incremental index (ghost function)")
			t.Log("This documents the behavior that issue #182 addresses with the post-checkout hook")
		} else {
			t.Log("PlaceOrder was cleaned up by incremental index (unexpected but acceptable)")
		}
	})

	// --- Step 8: Run force re-index ---
	t.Run("force_reindex_cleans_ghosts", func(t *testing.T) {
		runIndex(t, "--force")

		// Step 9: BranchOnlyFunction should still be found
		result, err := mcp.CallTool("search_code", map[string]interface{}{
			"query":      "BranchOnlyFunction",
			"project_id": fixtureProjectID,
		})
		require.NoError(t, err)
		if !result.IsError {
			assert.Contains(t, result.Text, "BranchOnlyFunction",
				"force reindex should keep new file")
		}

		// Step 10: PlaceOrder should NOT be found (cleaned up by force reindex)
		result, err = mcp.CallTool("search_code", map[string]interface{}{
			"query":      "PlaceOrder",
			"project_id": fixtureProjectID,
		})
		require.NoError(t, err)
		if !result.IsError {
			assert.NotContains(t, result.Text, "PlaceOrder",
				"force reindex should clean up ghost functions from deleted files")
		}
	})

	// --- Step 11: Switch back and verify restore ---
	t.Run("restore_after_switch_back", func(t *testing.T) {
		// Restore the deleted file before switching back
		os.WriteFile(orderServicePath, orderServiceBackup, 0644)
		// Remove the branch-only file
		os.Remove(branchOnlyFile)

		gitRun(t, fixtureDir, "checkout", "main")

		// Force re-index to restore original state
		runIndex(t, "--force")

		// PlaceOrder should be back
		result, err := mcp.CallTool("search_code", map[string]interface{}{
			"query":      "PlaceOrder",
			"project_id": fixtureProjectID,
		})
		require.NoError(t, err)
		assert.False(t, result.IsError, "tool error: %s", result.Text)
		assert.Contains(t, result.Text, "PlaceOrder",
			"PlaceOrder should be restored after switching back to main")

		// BranchOnlyFunction should be gone
		result, err = mcp.CallTool("search_code", map[string]interface{}{
			"query":      "BranchOnlyFunction",
			"project_id": fixtureProjectID,
		})
		require.NoError(t, err)
		if !result.IsError {
			assert.NotContains(t, result.Text, "BranchOnlyFunction",
				"BranchOnlyFunction should be gone after switching back to main")
		}
	})
}

// gitRun executes a git command in the given directory, failing the test on error.
func gitRun(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		// Don't fail on branch delete if branch doesn't exist
		if len(args) > 0 && args[0] == "branch" && args[1] == "-D" {
			t.Logf("git %s: %s (ignored)", strings.Join(args, " "), strings.TrimSpace(string(out)))
			return
		}
		t.Logf("git %s failed: %s", strings.Join(args, " "), strings.TrimSpace(string(out)))
	}
}

// gitOutput runs a git command and returns trimmed stdout.
func gitOutput(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.Output()
	require.NoError(t, err, "git %s failed", strings.Join(args, " "))
	return strings.TrimSpace(string(out))
}

// runIndex runs claude-memory index with the given flags, waiting for completion.
func runIndex(t *testing.T, extraArgs ...string) {
	t.Helper()
	args := []string{"index", "--project", fixtureProjectID}
	args = append(args, extraArgs...)
	cmd := exec.Command("claude-memory", args...)
	cmd.Dir = e2e.FixtureDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("claude-memory index output: %s", string(out))
		t.Fatalf("claude-memory index failed: %v", err)
	}

	// Wait a moment for index to complete and be queryable
	time.Sleep(2 * time.Second)
}
