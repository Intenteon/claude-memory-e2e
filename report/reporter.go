// Package report provides structured test result reporting.
package report

import (
	"fmt"
	"strings"
	"time"
)

// TestResult holds the outcome of a single test.
type TestResult struct {
	Name     string
	Passed   bool
	Duration time.Duration
	Error    string
	Skipped  bool
}

// Report aggregates test results and prints a summary.
type Report struct {
	Results []TestResult
	Start   time.Time
}

// NewReport creates a new report with the start time set to now.
func NewReport() *Report {
	return &Report{Start: time.Now()}
}

// Add records a test result.
func (r *Report) Add(result TestResult) {
	r.Results = append(r.Results, result)
}

// PrintSummary prints a formatted pass/fail report.
func (r *Report) PrintSummary() {
	totalDuration := time.Since(r.Start)

	passed := 0
	failed := 0
	skipped := 0
	for _, res := range r.Results {
		if res.Skipped {
			skipped++
		} else if res.Passed {
			passed++
		} else {
			failed++
		}
	}

	fmt.Println()
	fmt.Println(strings.Repeat("=", 70))
	fmt.Printf("  claude-memory E2E Test Report\n")
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println()

	// Print failures first
	if failed > 0 {
		fmt.Println("FAILURES:")
		for _, res := range r.Results {
			if !res.Passed && !res.Skipped {
				fmt.Printf("  FAIL  %-50s %s\n", res.Name, res.Duration.Round(time.Millisecond))
				if res.Error != "" {
					fmt.Printf("        %s\n", res.Error)
				}
			}
		}
		fmt.Println()
	}

	// Print skipped
	if skipped > 0 {
		fmt.Println("SKIPPED:")
		for _, res := range r.Results {
			if res.Skipped {
				fmt.Printf("  SKIP  %-50s %s\n", res.Name, res.Error)
			}
		}
		fmt.Println()
	}

	// Summary line
	fmt.Printf("Total: %d  |  Passed: %d  |  Failed: %d  |  Skipped: %d  |  Duration: %s\n",
		len(r.Results), passed, failed, skipped, totalDuration.Round(time.Millisecond))
	fmt.Println(strings.Repeat("=", 70))
}

// AllPassed returns true if all non-skipped tests passed.
func (r *Report) AllPassed() bool {
	for _, res := range r.Results {
		if !res.Passed && !res.Skipped {
			return false
		}
	}
	return true
}
