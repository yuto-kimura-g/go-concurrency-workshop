package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bytedance/sonic"
)

// OptimizedResult はPhase4専用の最適化版Result構造体
// StatusCountsをmapではなく配列にすることでメモリアロケーションとアクセス速度を改善
type OptimizedResult struct {
	FileName     string
	TotalCount   int
	StatusCounts [600]int // 固定配列でステータスコード0-599をカバー
}

// MinimalLogEntry は最小限のフィールドのみをパースする構造体
// Phase4では status フィールドのみが必要なため、他のフィールドをパースしない
type MinimalLogEntry struct {
	Status int `json:"status"`
}

func main() {
	startTime := time.Now()

	logRoot, err := os.OpenRoot("./logs")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening log directory: %v\n", err)
		os.Exit(1)
	}
	defer logRoot.Close()

	entries, err := fs.ReadDir(logRoot.FS(), ".")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading log directory: %v\n", err)
		os.Exit(1)
	}

	files := make([]string, 0, len(entries))
	for _, entry := range entries {
		name := entry.Name()
		if !entry.IsDir() && strings.HasPrefix(name, "access_") && strings.HasSuffix(name, ".json") {
			files = append(files, name)
		}
	}

	numWorkers := runtime.NumCPU()
	results := processFiles(logRoot, files, numWorkers)

	elapsed := time.Since(startTime)
	printResults(results, elapsed)
	recordResult("phase4", elapsed)
}

// processFiles は最適化されたワーカープールパターンでファイルを処理します
func processFiles(root *os.Root, files []string, numWorkers int) []*OptimizedResult {
	//  ジョブチャネルは小さいバッファで十分
	jobs := make(chan string, numWorkers)
	//  ワーカーごとの集計結果を受け取る（ファイル数ではなくワーカー数）
	results := make(chan *OptimizedResult, numWorkers)

	var wg sync.WaitGroup

	// ワーカーを起動
	for range numWorkers {
		wg.Go(func() {
			//  ワーカーごとにローカル集計（複数ファイルを1つの結果にまとめる）
			localResult := &OptimizedResult{
				FileName: "worker-aggregate",
			}

			for filename := range jobs {
				if err := processFileInto(root, filename, localResult); err != nil {
					fmt.Fprintf(os.Stderr, "Error processing %s: %v\n", filename, err)
				}
			}

			// ワーカーが処理した全ファイルの集計結果を送信
			results <- localResult
		})
	}

	//  不要なgoroutineを削除して直接ジョブを送信
	for _, filename := range files {
		jobs <- filename
	}
	close(jobs)

	// 全てのワーカーが完了したら結果チャネルを閉じる
	go func() {
		wg.Wait()
		close(results)
	}()

	// ワーカーごとの集計結果を収集（ファイル数ではなくワーカー数分）
	resultList := make([]*OptimizedResult, 0, numWorkers)
	for result := range results {
		resultList = append(resultList, result)
	}

	return resultList
}

// processFileInto は1つのログファイルを解析して既存の結果に集計します
// ワーカーごとのローカル集計に使用され、Result作成のオーバーヘッドを削減
func processFileInto(root *os.Root, filename string, result *OptimizedResult) error {
	file, err := root.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	//  256KBのバッファでI/O効率を向上
	bufferedReader := bufio.NewReaderSize(file, 256*1024)
	//  sonic JSONデコーダーを使用（標準ライブラリより2-5倍高速）
	decoder := sonic.ConfigDefault.NewDecoder(bufferedReader)

	//  最小フィールドのみパースしてパース時間とメモリを削減
	// 必要なのは status のみなので、他の7フィールド(timestamp, method, path, etc.)は無視
	var entry MinimalLogEntry
	for decoder.More() {
		if err := decoder.Decode(&entry); err != nil {
			continue
		}
		result.TotalCount++
		result.StatusCounts[entry.Status]++
	}

	return nil
}

// printResults は処理結果を表示します
func printResults(results []*OptimizedResult, elapsed time.Duration) {
	totalRequests := 0
	var totalStatusCounts [600]int

	for _, result := range results {
		totalRequests += result.TotalCount
		for status := 0; status < 600; status++ {
			if count := result.StatusCounts[status]; count > 0 {
				totalStatusCounts[status] += count
			}
		}
	}

	fmt.Printf("\n=== 処理結果 ===\n")
	fmt.Printf("処理時間: %.2f秒\n", elapsed.Seconds())
	fmt.Printf("総リクエスト数: %s件\n", formatNumber(totalRequests))
	fmt.Printf("\nステータスコード別:\n")
	for status := 200; status <= 599; status += 100 {
		for s := status; s < status+100; s++ {
			if count := totalStatusCounts[s]; count > 0 {
				percentage := float64(count) / float64(totalRequests) * 100
				fmt.Printf("  %d: %s件 (%.2f%%)\n", s, formatNumber(count), percentage)
			}
		}
	}

	errorCount := 0
	for status := 400; status < 600; status++ {
		errorCount += totalStatusCounts[status]
	}
	errorRate := float64(errorCount) / float64(totalRequests) * 100
	fmt.Printf("\nエラー率 (4xx, 5xx): %.2f%%\n", errorRate)
}

// formatNumber は数値を3桁カンマ区切りでフォーマットします
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

// recordResult は実行時間をsolutions/results.txtに記録します
func recordResult(phase string, elapsed time.Duration) {
	// 既存の結果を読み込む（なければ空のマップ）
	results := loadResults("./solutions/results.txt")

	// 現在のフェーズの結果を更新（冪等操作）
	results[phase] = elapsed.Seconds()

	// ファイルに書き戻す（全体を上書き）
	if err := saveResults("./solutions/results.txt", results); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to save results: %v\n", err)
	}
}

// loadResults はresults.txtを読み込む（ファイルがなければ空のマップ）
func loadResults(path string) map[string]float64 {
	results := make(map[string]float64)

	data, err := os.ReadFile(path)
	if err != nil {
		return results // ファイルがなければ空
	}

	for _, line := range strings.Split(string(data), "\n") {
		if line = strings.TrimSpace(line); line == "" {
			continue
		}
		// "phase1=10.00" or "phase2=2.00 (phase1から5.00倍高速, 80.0%改善)" の形式をパース
		parts := strings.Split(line, "=")
		if len(parts) == 2 {
			// スペースで分割して数値部分だけを取得
			fields := strings.Fields(parts[1])
			if len(fields) > 0 {
				if val, err := strconv.ParseFloat(fields[0], 64); err == nil {
					results[parts[0]] = val
				}
			}
		}
	}

	return results
}

// saveResults は結果をresults.txtに保存（冪等）
func saveResults(path string, results map[string]float64) error {
	var lines []string

	// Phase 1の基準値を取得
	baseline, hasBaseline := results["phase1"]

	// 安定したソート順（phase1, phase2, phase3...）
	for _, phase := range []string{"phase1", "phase2", "phase3", "phase4"} {
		if val, ok := results[phase]; ok {
			line := fmt.Sprintf("%s=%.2f", phase, val)

			// Phase 1以外で基準値があれば改善率を追加
			if phase != "phase1" && hasBaseline && baseline > 0 {
				improvement := (baseline - val) / baseline * 100
				speedup := baseline / val
				line += fmt.Sprintf(" (phase1から%.2f倍高速, %.1f%%改善)", speedup, improvement)
			}

			lines = append(lines, line)
		}
	}

	content := strings.Join(lines, "\n")
	if len(content) > 0 {
		content += "\n"
	}
	return os.WriteFile(path, []byte(content), 0644)
}
