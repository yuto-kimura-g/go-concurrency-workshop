// Package logparser provides utilities for parsing JSON access log files.
package logparser

import (
	"encoding/json"
	"fmt"
)

// LogEntry represents a single access log entry in JSON format.
type LogEntry struct {
	Timestamp      string `json:"timestamp"`
	Method         string `json:"method"`
	Path           string `json:"path"`
	Status         int    `json:"status"`
	ResponseTimeMs int    `json:"response_time_ms"`
	Bytes          int    `json:"bytes"`
	UserID         string `json:"user_id"`
	IP             string `json:"ip"`
}

// ParseLine parses a single JSON log line into a LogEntry.
// Returns an error if the JSON is invalid or malformed.
func ParseLine(line string) (*LogEntry, error) {
	var entry LogEntry
	if err := json.Unmarshal([]byte(line), &entry); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}
	return &entry, nil
}

// Result represents the analysis result for a single log file.
type Result struct {
	FileName     string
	TotalCount   int
	StatusCounts map[int]int
}

// NewResult creates a new Result with initialized maps.
func NewResult(filename string) *Result {
	return &Result{
		FileName:     filename,
		StatusCounts: make(map[int]int),
	}
}

// AddEntry adds a log entry to the result, updating counters.
func (r *Result) AddEntry(entry *LogEntry) {
	r.TotalCount++
	r.StatusCounts[entry.Status]++
}

// TotalResult represents the aggregated result from all log files.
type TotalResult struct {
	FileCount    int
	TotalCount   int
	StatusCounts map[int]int
}

// NewTotalResult creates a new TotalResult with initialized maps.
func NewTotalResult() *TotalResult {
	return &TotalResult{
		StatusCounts: make(map[int]int),
	}
}

// MergeResults merges multiple Results into a single TotalResult.
func MergeResults(results []*Result) *TotalResult {
	total := NewTotalResult()
	total.FileCount = len(results)

	for _, r := range results {
		total.TotalCount += r.TotalCount
		for status, count := range r.StatusCounts {
			total.StatusCounts[status] += count
		}
	}

	return total
}

// ErrorRate calculates the percentage of 4xx and 5xx status codes.
func (tr *TotalResult) ErrorRate() float64 {
	if tr.TotalCount == 0 {
		return 0.0
	}

	errorCount := 0
	for status, count := range tr.StatusCounts {
		if status >= 400 && status < 600 {
			errorCount += count
		}
	}

	return float64(errorCount) / float64(tr.TotalCount) * 100
}
