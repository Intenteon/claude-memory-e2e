.PHONY: test test-fixture test-real test-search test-callgraph test-ollama setup build clean

# Default: run all tests
test: build
	go test ./tests/... -v -timeout 300s

# Run only fixture-dependent tests (no real projects, no ollama)
test-fixture: build
	go test ./tests/... -v -timeout 300s -run 'TestSearch|TestFuzzy|TestCallGraph|TestGetCallees|TestGetCompact|TestGetInterface|TestGetProjectSummary|TestValidate|TestGoal|TestPlan|TestLearning|TestExport'

# Run only real project smoke tests
test-real: build
	go test ./tests/... -v -timeout 120s -run 'TestRealProject'

# Run only search tests (quick sanity check)
test-search: build
	go test ./tests/... -v -timeout 60s -run 'TestSearch'

# Run only call graph tests
test-callgraph: build
	go test ./tests/... -v -timeout 60s -run 'TestGetCallers|TestGetCallees|TestGetCompact'

# Run only ollama tests
test-ollama: build
	go test ./tests/... -v -timeout 300s -run 'TestOllama'

# Run branch switch ghost function test
test-branch: build
	go test ./tests/... -v -run TestBranchSwitch -timeout 300s

# Setup the fixture project (clone, index)
setup:
	bash setup-fixture.sh

# Build (verify compilation)
build:
	go build ./...

# Clean fixture directory
clean:
	rm -rf /tmp/claude-memory-e2e-fixture
