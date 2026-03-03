package e2e

const (
	// FixtureRepo is the GitHub URL for the fixture project.
	FixtureRepo = "https://github.com/Intenteon/claude-memory-fixture"
	// FixtureSHA is the pinned commit SHA of the fixture project.
	FixtureSHA = "e4c954ce4783ca37bb1c5f2e1a45825fbed8e93f"
	// FixtureDir is the default local path for the cloned fixture.
	FixtureDir = "/tmp/claude-memory-e2e-fixture"

	// MCPServerURL is the MCP JSON-RPC endpoint.
	MCPServerURL = "http://localhost:7677"
	// MCPEndpoint is the full MCP endpoint path.
	MCPEndpoint = MCPServerURL + "/mcp"
	// HealthEndpoint is the health check endpoint.
	HealthEndpoint = MCPServerURL + "/health"
	// APIServerURL is the REST API server.
	APIServerURL = "http://localhost:7676"

	// Real project IDs for smoke testing.
	AuthServiceProjectID = "auth-service-eae34498"
	AuthSDKProjectID     = "auth-sdk-go-008e1dc8"
)
