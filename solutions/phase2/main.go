// Phase 2: Basic Concurrent Processing
//
// This solution processes log files concurrently using goroutines and channels.
// Each file is processed in its own goroutine, and results are collected via a channel.
// This demonstrates the basic pattern of concurrent processing in Go.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/nnnkkk7/go-concurrency-workshop/pkg/logparser"
)

func main() {
	// Get list of log files
	files, err := filepath.Glob("logs/*.log")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding log files: %v\n", err)
		os.Exit(1)
	}

	if len(files) == 0 {
		fmt.Println("No log files found in ./logs directory")
		fmt.Println("Please run: go run cmd/loggen/main.go")
		os.Exit(1)
	}

	fmt.Printf("Processing %d log files concurrently...\n", len(files))

	// Start timing
	start := time.Now()

	// Process files concurrently
	results := processFilesConcurrently(files)

	// Merge results
	total := logparser.MergeResults(results)

	// Calculate processing time
	elapsed := time.Since(start)

	// Display results
	printResults(total, elapsed)
}

// processFilesConcurrently spawns a goroutine for each file and collects results via a channel.
func processFilesConcurrently(files []string) []*logparser.Result {
	// Create channel to collect results
	resultCh := make(chan *logparser.Result, len(files))

	// Use WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Launch a goroutine for each file
	for _, filename := range files {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()

			result, err := processFile(name)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error processing %s: %v\n", name, err)
				return
			}

			// Send result to channel
			resultCh <- result
		}(filename)
	}

	// Close the channel when all goroutines are done
	go func() {
		wg.Wait()
		close(resultCh)
	}()

	// Collect all results from channel
	results := make([]*logparser.Result, 0, len(files))
	for result := range resultCh {
		results = append(results, result)
	}

	return results
}

// processFile reads a single log file and counts status codes.
func processFile(filename string) (*logparser.Result, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	result := logparser.NewResult(filename)
	decoder := json.NewDecoder(file)

	for decoder.More() {
		var entry logparser.LogEntry
		if err := decoder.Decode(&entry); err != nil {
			// Skip invalid lines
			continue
		}
		result.AddEntry(&entry)
	}

	return result, nil
}

// printResults displays the analysis results in a formatted way.
func printResults(total *logparser.TotalResult, elapsed time.Duration) {
	fmt.Println("\n=== Access Log Analysis Results ===")
	fmt.Printf("Total files: %d\n", total.FileCount)
	fmt.Printf("Total requests: %s\n", formatNumber(total.TotalCount))
	fmt.Printf("Processing time: %v\n", elapsed.Round(time.Millisecond))
	fmt.Println()

	fmt.Println("Status Code Distribution:")

	// Sort status codes for consistent output
	statuses := make([]int, 0, len(total.StatusCounts))
	for status := range total.StatusCounts {
		statuses = append(statuses, status)
	}
	sort.Ints(statuses)

	for _, status := range statuses {
		count := total.StatusCounts[status]
		percentage := float64(count) / float64(total.TotalCount) * 100
		fmt.Printf("  %d: %s (%.2f%%)\n", status, formatNumber(count), percentage)
	}

	fmt.Println()
	fmt.Printf("Error Rate (4xx + 5xx): %.2f%%\n", total.ErrorRate())
}

// formatNumber adds commas to numbers for readability.
func formatNumber(n int) string {
	if n < 1000 {
		return fmt.Sprintf("%d", n)
	}

	// Convert to string with commas
	str := fmt.Sprintf("%d", n)
	var result []byte
	for i, c := range str {
		if i > 0 && (len(str)-i)%3 == 0 {
			result = append(result, ',')
		}
		result = append(result, byte(c))
	}
	return string(result)
}
