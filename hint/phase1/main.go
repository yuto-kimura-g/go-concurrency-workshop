// Phase 1: Sequential Processing
//
// このPhaseでは、goroutineやchannelを使わずに、
// シンプルなfor文で全ファイルを順番に処理します。
//
// 制約:
// - ❌ goroutine使用禁止
// - ❌ channel使用禁止
// - ❌ syncパッケージ使用禁止
// - ✅ 普通のfor文のみ使用
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/nnnkkk7/go-concurrency-workshop/pkg/logparser"
)

func main() {
	// ログファイル一覧を取得
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

	fmt.Printf("Processing %d log files sequentially...\n", len(files))

	// 処理時間の計測開始
	start := time.Now()

	// ファイルを逐次処理
	results := processFilesSequentially(files)

	// 結果をマージ
	total := logparser.MergeResults(results)

	// 処理時間を計算
	elapsed := time.Since(start)

	// 結果を表示
	printResults(total, elapsed)
}

// processFilesSequentially は全ファイルを順番に処理します
//
// TODO: この関数を実装してください
//
// HINT:
// 1. 結果を格納するスライスを作成: results := make([]*logparser.Result, 0, len(files))
// 2. for文でfilesをループ
// 3. 各ファイルをprocessFile関数で処理
// 4. 結果をresultsスライスに追加
// 5. resultsを返す
func processFilesSequentially(files []string) []*logparser.Result {
	results := make([]*logparser.Result, 0, len(files))

	// TODO: ここにfor文を書いて、各ファイルを処理してください
	// for _, filename := range files {
	//     result, err := processFile(filename)
	//     if err != nil {
	//         fmt.Fprintf(os.Stderr, "Error processing %s: %v\n", filename, err)
	//         continue
	//     }
	//     results = append(results, result)
	// }

	return results
}

// processFile は1つのログファイルを読み込んでステータスコードを集計します
// （この関数は完成形で提供されています）
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
			// 不正な行はスキップ
			continue
		}
		result.AddEntry(&entry)
	}

	return result, nil
}

// printResults は集計結果を見やすく表示します
// （この関数は完成形で提供されています）
func printResults(total *logparser.TotalResult, elapsed time.Duration) {
	fmt.Println("\n=== Access Log Analysis Results ===")
	fmt.Printf("Total files: %d\n", total.FileCount)
	fmt.Printf("Total requests: %s\n", formatNumber(total.TotalCount))
	fmt.Printf("Processing time: %v\n", elapsed.Round(time.Millisecond))
	fmt.Println()

	fmt.Println("Status Code Distribution:")

	// ステータスコードをソート
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

// formatNumber は数値にカンマを入れて見やすくします
// （この関数は完成形で提供されています）
func formatNumber(n int) string {
	if n < 1000 {
		return fmt.Sprintf("%d", n)
	}

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
