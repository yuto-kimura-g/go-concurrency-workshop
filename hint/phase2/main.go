// Phase 2: Concurrent Processing with Goroutines and Channels
//
// このPhaseでは、goroutineとchannelを使って、
// 複数のファイルを並行に処理します。
//
// 目標: Phase 1より5〜10倍速く！
//
// 使用可能:
// - ✅ goroutine (go キーワード)
// - ✅ channel
// - ✅ sync.WaitGroup
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

	fmt.Printf("Processing %d log files concurrently...\n", len(files))

	// 処理時間の計測開始
	start := time.Now()

	// ファイルを並行処理
	results := processFilesConcurrently(files)

	// 結果をマージ
	total := logparser.MergeResults(results)

	// 処理時間を計算
	elapsed := time.Since(start)

	// 結果を表示
	printResults(total, elapsed)
}

// processFilesConcurrently は全ファイルを並行処理します
//
// TODO: この関数を実装してください
//
// HINT:
//
//  1. バッファ付きchannelを作成: resultCh := make(chan *logparser.Result, len(files))
//
//  2. 各ファイルに対してgoroutineを起動:
//     for _, filename := range files {
//     go func(name string) {
//     result, err := processFile(name)
//     if err == nil {
//     resultCh <- result
//     }
//     }()
//     }
//
//  3. ファイル数だけループして結果を受信:
//     for i := 0; i < len(files); i++ {
//     result := <-resultCh
//     results = append(results, result)
//     }
//
//     go func(name string) { ... }(filename) のように引数で渡す
func processFilesConcurrently(files []string) []*logparser.Result {
	// TODO: ここにgoroutineとchannelを使った並行処理を実装してください

	// ヒント:
	// 1. channelを作成
	// 2. 各ファイルに対してgoroutineを起動
	// 3. channelから結果を収集
	// return results

	return nil // TODO: 実装後は実際のresultsを返す
}

// === 代替実装案: WaitGroupを使う方法 ===
//
// 上記のchannel方式とは別に、sync.WaitGroupとMutexを使う方法もあります。
// 時間に余裕があれば、こちらも試してみてください。
//
// func processFilesConcurrentlyWithWaitGroup(files []string) []*logparser.Result {
//     var wg sync.WaitGroup
//     var mu sync.Mutex
//     results := make([]*logparser.Result, 0, len(files))
//
//     for _, filename := range files {
//         wg.Add(1)
//         go func(name string) {
//             defer wg.Done()
//             result, err := processFile(name)
//             if err != nil {
//                 return
//             }
//             mu.Lock()
//             results = append(results, result)
//             mu.Unlock()
//         }(filename)
//     }
//
//     wg.Wait()
//     return results
// }

// processFile は1つのログファイルを読み込んでステータスコードを集計します
// （Phase 1と同じ）
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
// （Phase 1と同じ）
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
// （Phase 1と同じ）
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
