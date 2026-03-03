# claude-memory E2E Test Suite

Comprehensive end-to-end tests for all [claude-memory](https://github.com/stevehockey/claude-memory) MCP tools.

## Prerequisites

- `claude-memory` CLI installed and on PATH
- MCP server running: `claude-memory serve`
- Ollama running (optional, for ollama_* tests): `ollama serve`
- Go 1.21+

## Quick Start

```bash
# 1. Setup the fixture project (clone, init, index)
make setup

# 2. Run all tests
make test
```

## Test Categories

| Command | What it tests |
|---|---|
| `make test` | All tests |
| `make test-fixture` | Fixture-dependent tests (search, callgraph, goals, plans, etc.) |
| `make test-real` | Smoke tests against real projects (auth-service, auth-sdk-go) |
| `make test-search` | Quick sanity check — just search_code and fuzzy_search |
| `make test-callgraph` | Call graph tools (get_callers, get_callees, get_compact_callers) |
| `make test-ollama` | All ollama_* tools (requires Ollama running) |

## Architecture

```
├── config.go              — Constants: fixture SHA, MCP URL, project IDs
├── client/
│   └── mcp.go             — JSON-RPC 2.0 client for MCP tools/call
├── report/
│   └── reporter.go        — Structured test result reporter
├── tests/
│   ├── main_test.go       — TestMain: health check, fixture discovery
│   ├── search_test.go     — search_code, fuzzy_search_functions, find_similar_code
│   ├── callgraph_test.go  — get_callers, get_callees, get_compact_callers
│   ├── interface_test.go  — get_interface, get_project_summary
│   ├── validate_test.go   — validate_changes
│   ├── memory_test.go     — record_memory, get_requirements
│   ├── goals_test.go      — Goal lifecycle (create → update → complete)
│   ├── plans_test.go      — Plan lifecycle (create → add step → complete)
│   ├── learning_test.go   — Session learning tools
│   ├── codegen_test.go    — generate_stub, generate_docs, generate_test, suggest_refactoring
│   ├── security_test.go   — security_scan, get_vulnerabilities
│   ├── git_test.go        — Git context tools (history, blame, hotspots, churn)
│   ├── patterns_test.go   — Pattern library tools
│   ├── export_test.go     — export_graph (mermaid, dot, json)
│   ├── ollama_test.go     — All ollama_* tools
│   └── real_projects_test.go — Smoke tests against real indexed projects
├── setup-fixture.sh       — Clone fixture at pinned SHA and index it
└── Makefile               — Build and test targets
```

## Fixture Project

The test suite uses [claude-memory-fixture](https://github.com/Intenteon/claude-memory-fixture), a deterministic Go project with known call graphs, interfaces, and search targets. Tests assert exact expected results against this fixture.

## Environment Variables

| Variable | Description |
|---|---|
| `FIXTURE_PROJECT_ID` | Override auto-discovery of fixture project ID |
| `SKIP_FIXTURE_SETUP` | Skip auto-setup in TestMain (use pre-indexed fixture) |
