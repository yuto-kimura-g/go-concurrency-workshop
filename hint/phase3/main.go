// Phase 3: Worker Pool Pattern with WaitGroup.Go()
//
// このPhaseでは、ワーカープールパターンを使って、
// 同時実行数を制御しながら効率的に処理します。
//
// Go 1.25の新機能 WaitGroup.Go() を活用します！
//
// 目標: Phase 2よりさらに速く（または安定した性能）
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"time"

	"github.com/nnnkkk7/go-concurrency-workshop/pkg/logparser"
	// "sync" // TODO: WaitGroupを使う場合はこのimportを有効にしてください
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

	// ワーカー数を決定（CPU数が最適）
	numWorkers := runtime.NumCPU()
	fmt.Printf("Processing %d log files with %d workers...\n", len(files), numWorkers)

	// 処理時間の計測開始
	start := time.Now()

	// ワーカープールで処理
	results := processFiles(files, numWorkers)

	// 結果をマージ
	total := logparser.MergeResults(results)

	// 処理時間を計算
	elapsed := time.Since(start)

	// 結果を表示
	printResults(total, elapsed)
}

// processFilesWithWorkerPool はワーカープールパターンで処理します
//
// TODO: この関数を実装してください
//
// ワーカープールの仕組み:
// 1. 固定数のワーカーgoroutineを起動
// 2. ジョブ（ファイル名）をchannelで配布
// 3. ワーカーがジョブを取得して処理
// 4. 結果をchannelで収集
//
// HINT:
//
//  1. channelを作成:
//     jobs := make(chan string, len(files))
//     results := make(chan *logparser.Result, len(files))
//
//  2. ワーカーを起動（Go 1.25の新機能！）:
//     var wg sync.WaitGroup
//     for i := 0; i < numWorkers; i++ {
//     wg.Go(func() {  // ← wg.Add(1) + go func() が1行に！
//     for filename := range jobs {
//     result, _ := processFile(filename)
//     results <- result
//     }
//     })
//     }
//
//  3. ジョブを投入:
//     for _, filename := range files {
//     jobs <- filename
//     }
//     close(jobs)  // ← 重要！ワーカーに「もうジョブはない」と伝える
//
//  4. Wait後にresultsをclose:
//     go func() {
//     wg.Wait()
//     close(results)
//     }()
//
//  5. 結果を収集:
//     resultList := make([]*logparser.Result, 0, len(files))
//     for result := range results {
//     resultList = append(resultList, result)
//     }
func processFiles(files []string, numWorkers int) []*logparser.Result {
	// TODO: ここにワーカープールパターンを実装してください

	// HINT: channelを作成
	// jobs := make(chan string, len(files))
	// results := make(chan *logparser.Result, len(files))

	// HINT: ワーカーを起動（Go 1.25の新機能）
	// var wg sync.WaitGroup
	// for i := 0; i < numWorkers; i++ {
	//     wg.Go(func() {
	//         for filename := range jobs {
	//             result, err := processFile(filename)
	//             if err != nil {
	//                 fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	//                 continue
	//             }
	//             results <- result
	//         }
	//     })
	// }

	// HINT: ジョブを投入
	// for _, filename := range files {
	//     jobs <- filename
	// }
	// close(jobs)

	// HINT: Wait後にresultsをclose
	// go func() {
	//     wg.Wait()
	//     close(results)
	// }()

	// HINT: 結果を収集
	// resultList := make([]*logparser.Result, 0, len(files))
	// for result := range results {
	//     resultList = append(resultList, result)
	// }

	// return resultList

	return nil // TODO: 実装後は実際のresultListを返す
}

// === Go 1.24以前の書き方（参考） ===
//
// Go 1.25未満の環境では、以下のように書きます：
//
// for i := 0; i < numWorkers; i++ {
//     wg.Add(1)
//     go func() {
//         defer wg.Done()
//         for filename := range jobs {
//             result, _ := processFile(filename)
//             results <- result
//         }
//     }()
// }
//
// Go 1.25では wg.Go() を使うことで、Add(1) + go func() + defer wg.Done()
// がシンプルに書けるようになりました！

// processFile は1つのログファイルを読み込んでステータスコードを集計します
// （Phase 1/2と同じ）
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

// printResults は集計結果を見やすく表示します
// （Phase 1/2と同じ）
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

// formatNumber は数値にカンマを入れて見やすくします
// （Phase 1/2と同じ）
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
