# 環境確認チェックリスト

## 1. Go のインストール確認

### 1.1 Go 1.25 以降がインストールされているか

```bash
go version
```

**期待される出力例:**
```
go version go1.25.5 darwin/arm64
```

**確認ポイント:**
- [ ] `go version`が表示される
- [ ] バージョンが1.25.0以降である

**問題がある場合:**
- [Go公式サイト](https://go.dev/dl/)からインストール


## 2. リポジトリのクローン確認

### 2.1 リポジトリをクローン

```bash
git clone https://github.com/nnnkkk7/go-concurrency-workshop.git
cd go-concurrency-workshop
```

**確認ポイント:**
- [ ] クローンが成功した
- [ ] ディレクトリに移動できた


---

## 3. ログファイルの生成

### 3.1 ログ生成ツールの実行

```bash
make gen
```

**期待される出力例:**
```
Cleaned up 200 existing log file(s)
Generating log files...

Done! Generated 200 files (1645.1MB total) in 25.386s
```


**確認ポイント:**
- [ ] エラーなく完了した

### 3.2 ログファイルの確認

```bash
ls logs/*.json | wc -l
```

**期待される出力:**
```
200
```

**確認ポイント:**
- [ ] 200ファイルが存在する

