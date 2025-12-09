// Phase 4: Advanced Optimizations
//
// This solution applies various optimization techniques beyond basic worker pools:
// - Optimized buffer sizes for channels
// - Pre-allocated result slices
// - Result pooling to reduce allocations
// - Efficient channel communication patterns
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

	numWorkers := runtime.NumCPU()
	fmt.Printf("Processing %d log files with %d workers (optimized)...\n", len(files), numWorkers)

	start := time.Now()

	results := processFilesOptimized(files, numWorkers)

	total := logparser.MergeResults(results)

	elapsed := time.Since(start)

	printResults(total, elapsed)
}

// processFilesOptimized uses advanced optimization techniques
func processFilesOptimized(files []string, numWorkers int) []*logparser.Result {
	fileCount := len(files)

	// Use smaller buffer for jobs (only needs to hold pending work)
	jobs := make(chan string, numWorkers*2)
	// Use larger buffer for results (holds all results)
	results := make(chan *logparser.Result, fileCount)

	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Go(func() {
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

	// Feed jobs (non-blocking with buffered channel)
	go func() {
		for _, filename := range files {
			jobs <- filename
		}
		close(jobs)
	}()

	// Wait for all workers to finish
	go func() {
		wg.Wait()
		close(results)
	}()

	// Pre-allocate result slice with exact capacity
	resultList := make([]*logparser.Result, 0, fileCount)
	for result := range results {
		resultList = append(resultList, result)
	}

	return resultList
}

// processFile processes a single file with optimized I/O
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
			continue
		}
		result.AddEntry(&entry)
	}

	return result, nil
}

func printResults(total *logparser.TotalResult, elapsed time.Duration) {
	fmt.Println("\n=== Access Log Analysis Results ===")
	fmt.Printf("Total files: %d\n", total.FileCount)
	fmt.Printf("Total requests: %s\n", formatNumber(total.TotalCount))
	fmt.Printf("Processing time: %v\n", elapsed.Round(time.Millisecond))
	fmt.Println()

	fmt.Println("Status Code Distribution:")

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

func formatNumber(n int) string {
	s := fmt.Sprintf("%d", n)
	result := ""
	for i, c := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			result += ","
		}
		result += string(c)
	}
	return result
}
