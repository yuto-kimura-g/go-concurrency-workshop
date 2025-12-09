// Phase 3: Worker Pool with WaitGroup.Go() (Go 1.25+)
//
// This solution uses a worker pool pattern to limit concurrent file processing.
// It leverages Go 1.25's new WaitGroup.Go() method for cleaner goroutine management.
// The worker pool prevents resource exhaustion when processing many files.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
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

	// Determine optimal worker count (typically number of CPU cores)
	numWorkers := runtime.NumCPU()
	fmt.Printf("Processing %d log files with %d workers...\n", len(files), numWorkers)

	// Start timing
	start := time.Now()

	// Process files using worker pool
	results := processFilesWithWorkerPool(files, numWorkers)

	// Merge results
	total := logparser.MergeResults(results)

	// Calculate processing time
	elapsed := time.Since(start)

	// Display results
	printResults(total, elapsed)
}

// processFilesWithWorkerPool uses a worker pool pattern with Go 1.25's WaitGroup.Go().
func processFilesWithWorkerPool(files []string, numWorkers int) []*logparser.Result {
	// Create buffered channels for job distribution and result collection
	jobs := make(chan string, len(files))
	results := make(chan *logparser.Result, len(files))

	// Create WaitGroup for worker coordination
	var wg sync.WaitGroup

	// Start worker goroutines using WaitGroup.Go() (Go 1.25+)
	for i := 0; i < numWorkers; i++ {
		wg.Go(func() {
			// Each worker processes files from the jobs channel
			for filename := range jobs {
				result, err := processFile(filename)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error processing %s: %v\n", filename, err)
					continue
				}
				results <- result
			}
		})
	}

	// Send all files to the jobs channel
	for _, filename := range files {
		jobs <- filename
	}
	close(jobs) // Signal that no more jobs will be sent

	// Wait for all workers to finish, then close results channel
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect all results
	resultList := make([]*logparser.Result, 0, len(files))
	for result := range results {
		resultList = append(resultList, result)
	}

	return resultList
}

// processFile reads a single log file and counts status codes.
// This function uses optimized buffering for better performance.
func processFile(filename string) (*logparser.Result, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	result := logparser.NewResult(filename)
	decoder := json.NewDecoder(file)

	// Process entries in batches for better performance
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
