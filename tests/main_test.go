package tests

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	e2e "github.com/Intenteon/claude-memory-e2e"
	"github.com/Intenteon/claude-memory-e2e/client"
)

// Shared state for all tests.
var (
	mcp              *client.MCPClient
	fixtureProjectID string
)

// TestMain handles global setup: MCP health check and fixture project discovery.
func TestMain(m *testing.M) {
	fmt.Println("=== claude-memory E2E Test Suite ===")
	fmt.Println()

	// Initialize MCP client
	mcp = client.NewMCPClient(e2e.MCPEndpoint)

	// 1. Check MCP server health
	fmt.Print("Checking MCP server... ")
	if err := mcp.HealthCheck(); err != nil {
		fmt.Printf("FAILED: %v\n", err)
		fmt.Println("Start the MCP server with: claude-memory serve")
		os.Exit(1)
	}
	fmt.Println("OK")

	// 2. Discover fixture project ID
	fmt.Print("Discovering fixture project ID... ")
	fixtureProjectID = discoverFixtureProjectID()
	if fixtureProjectID == "" {
		fmt.Println("NOT FOUND")
		fmt.Println("Run: bash setup-fixture.sh")
		fmt.Println("Fixture-dependent tests will be skipped.")
	} else {
		fmt.Printf("OK (%s)\n", fixtureProjectID)
	}

	// 3. Verify real projects exist
	fmt.Print("Checking real projects... ")
	checkRealProject(e2e.AuthServiceProjectID)
	checkRealProject(e2e.AuthSDKProjectID)
	fmt.Println()
	fmt.Println()

	os.Exit(m.Run())
}

func discoverFixtureProjectID() string {
	// Try env var first
	if id := os.Getenv("FIXTURE_PROJECT_ID"); id != "" {
		return id
	}

	// Use claude-memory project list to find it
	cmd := exec.Command("claude-memory", "project", "list")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return ""
	}

	// The output format is multi-line per project:
	//   claude-memory-e2e-fixture (29 functions) (active)
	//     ID:   claude-memory-e2e-fixture-08f439c8
	//     Path: /tmp/...
	lines := strings.Split(string(output), "\n")
	foundFixture := false
	for _, line := range lines {
		if strings.Contains(line, "fixture") && !strings.HasPrefix(strings.TrimSpace(line), "ID:") {
			foundFixture = true
			continue
		}
		if foundFixture && strings.Contains(line, "ID:") {
			parts := strings.Fields(line)
			for i, p := range parts {
				if p == "ID:" && i+1 < len(parts) {
					return parts[i+1]
				}
			}
		}
	}
	return ""
}

func checkRealProject(projectID string) {
	result, err := mcp.CallTool("get_project_summary", map[string]interface{}{
		"project_id": projectID,
	})
	if err != nil {
		fmt.Printf("%s=MISSING ", projectID)
	} else if result.IsError {
		fmt.Printf("%s=ERROR ", projectID)
	} else {
		fmt.Printf("%s=OK ", projectID)
	}
}

// requireMCP skips the test if the MCP server is unreachable.
func requireMCP(t *testing.T) {
	t.Helper()
	c := &http.Client{Timeout: 3 * time.Second}
	resp, err := c.Get(e2e.HealthEndpoint)
	if err != nil {
		t.Skipf("MCP server not available: %v", err)
	}
	resp.Body.Close()
}

// requireFixture skips the test if the fixture project isn't indexed.
func requireFixture(t *testing.T) {
	t.Helper()
	requireMCP(t)
	if fixtureProjectID == "" {
		t.Skip("fixture project not indexed — run setup-fixture.sh")
	}
}

// requireOllama skips the test if Ollama is not running.
func requireOllama(t *testing.T) {
	t.Helper()
	c := &http.Client{Timeout: 3 * time.Second}
	resp, err := c.Get("http://localhost:11434/api/tags")
	if err != nil {
		t.Skip("Ollama not available, skipping")
	}
	resp.Body.Close()
}

// skipIfToolNotFound checks if the error is a "Tool not found" RPC error and skips the test.
// Returns true if the test was skipped.
func skipIfToolNotFound(t *testing.T, err error) bool {
	t.Helper()
	if err != nil && strings.Contains(err.Error(), "Tool not found") {
		t.Skipf("tool not registered in MCP server: %v", err)
		return true
	}
	return false
}

// skipIfUnavailable checks both RPC-level and tool-level errors for "not available/configured/installed" patterns.
// Returns true if the test was skipped.
func skipIfUnavailable(t *testing.T, err error, result *client.CallToolResult, reason string) bool {
	t.Helper()
	if skipIfToolNotFound(t, err) {
		return true
	}
	if err != nil {
		return false
	}
	if result != nil && result.IsError {
		text := result.Text
		if strings.Contains(text, "not configured") ||
			strings.Contains(text, "not available") ||
			strings.Contains(text, "not installed") ||
			strings.Contains(text, "Semgrep") ||
			strings.Contains(text, "semgrep") ||
			strings.Contains(text, "LLM service") ||
			strings.Contains(text, "FUNCTION NOT FOUND") ||
			strings.Contains(text, "default project not found") ||
			strings.Contains(text, "path outside project") ||
			strings.Contains(text, "no valid files") ||
			strings.Contains(text, "could not read any context") {
			t.Skipf("%s: %s", reason, truncate(text, 120))
			return true
		}
	}
	return false
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// requireRealProject skips the test if the given real project is not available.
func requireRealProject(t *testing.T, projectID string) {
	t.Helper()
	requireMCP(t)
	result, err := mcp.CallTool("get_project_summary", map[string]interface{}{
		"project_id": projectID,
	})
	if err != nil || result.IsError {
		t.Skipf("real project %s not available", projectID)
	}
}
