// Package client provides a JSON-RPC client for the claude-memory MCP server.
package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync/atomic"
	"time"
)

// MCPClient is a JSON-RPC 2.0 client for the MCP server.
type MCPClient struct {
	endpoint string
	client   *http.Client
	nextID   atomic.Int64
}

// NewMCPClient creates a new MCP client targeting the given endpoint.
func NewMCPClient(endpoint string) *MCPClient {
	return &MCPClient{
		endpoint: endpoint,
		client: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

// JSONRPCRequest is a JSON-RPC 2.0 request.
type JSONRPCRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      int64       `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
}

// ToolCallParams wraps tool name and arguments for a tools/call request.
type ToolCallParams struct {
	Name      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments"`
}

// JSONRPCResponse is a JSON-RPC 2.0 response.
type JSONRPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      int64           `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *JSONRPCError   `json:"error,omitempty"`
}

// JSONRPCError represents a JSON-RPC error.
type JSONRPCError struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}

func (e *JSONRPCError) Error() string {
	return fmt.Sprintf("RPC error %d: %s", e.Code, e.Message)
}

// ToolResult is the MCP CallToolResult structure.
type ToolResult struct {
	Content []ToolContent `json:"content"`
	IsError bool          `json:"isError,omitempty"`
}

// ToolContent is a single content item in a tool result.
type ToolContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// CallToolResult holds the parsed result of a tool call.
type CallToolResult struct {
	Text    string // First text content
	IsError bool   // Whether the tool returned an error
	Raw     json.RawMessage
}

// CallTool invokes an MCP tool via JSON-RPC and returns the text result.
func (c *MCPClient) CallTool(toolName string, args interface{}) (*CallToolResult, error) {
	argsJSON, err := json.Marshal(args)
	if err != nil {
		return nil, fmt.Errorf("marshal args for %s: %w", toolName, err)
	}

	req := JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      c.nextID.Add(1),
		Method:  "tools/call",
		Params: ToolCallParams{
			Name:      toolName,
			Arguments: argsJSON,
		},
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	httpResp, err := c.client.Post(c.endpoint, "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("POST %s: %w", c.endpoint, err)
	}
	defer httpResp.Body.Close()

	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	var rpcResp JSONRPCResponse
	if err := json.Unmarshal(respBody, &rpcResp); err != nil {
		return nil, fmt.Errorf("parse response for %s: %w (body: %s)", toolName, err, string(respBody))
	}

	if rpcResp.Error != nil {
		return nil, rpcResp.Error
	}

	var toolResult ToolResult
	if err := json.Unmarshal(rpcResp.Result, &toolResult); err != nil {
		return nil, fmt.Errorf("parse tool result for %s: %w", toolName, err)
	}

	result := &CallToolResult{
		IsError: toolResult.IsError,
		Raw:     rpcResp.Result,
	}
	if len(toolResult.Content) > 0 {
		result.Text = toolResult.Content[0].Text
	}

	return result, nil
}

// HealthCheck verifies the MCP server is reachable.
func (c *MCPClient) HealthCheck() error {
	healthURL := c.endpoint[:len(c.endpoint)-4] + "/health" // replace /mcp with /health
	resp, err := c.client.Get(healthURL)
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check returned status %d", resp.StatusCode)
	}
	return nil
}
