# Goで学ぶ並行処理：アクセスログ解析チャレンジ

実務でよくある「大量のログファイルを急いで解析してほしい」という状況を題材に、Goの並行処理を学ぶハンズオンワークショップです。

##  必要な環境

- **Go 1.25以降**（WaitGroup.Go()を使用）
- エディタ（VS Code推奨）
- ターミナル

> **📋 環境構築の詳細確認:** ワークショップ当日までに [docs/SETUP_CHECK.md](docs/SETUP_CHECK.md) のチェックリストを完了してください。

##  クイックスタート

### 1. リポジトリをクローン

```bash
git clone https://github.com/nnnkkk7/go-concurrency-workshop.git
cd go-concurrency-workshop
```

### 2. ログファイルを生成

```bash
go run cmd/loggen/main.go
```

これで `logs/` ディレクトリに50個のログファイル（合計約500MB）が生成されます。

> **✅ 環境確認:** 正常に動作するか確認したい場合は [docs/SETUP_CHECK.md](docs/SETUP_CHECK.md) を参照してください。

### 3. ワークショップを開始

```bash
# workshop/ で実装に挑戦
go run ./workshop/phase1/main.go

```

##  ワークショップの進め方

**[workshop/README.md](workshop/README.md)** を参照してください。


##  ドキュメント

- [workshop/README.md](workshop/README.md) - ワークショップ実施ガイド（各Phaseの詳細な説明）
- [docs/SETUP_CHECK.md](docs/SETUP_CHECK.md) - 環境構築方法とチェックリスト
- [docs/HINTS.md](docs/HINTS.md) - 実装のヒント集


##  プロジェクト構成

```
go-concurrency-workshop/
├── cmd/
│   └── loggen/          # ログ生成ツール
├── pkg/
│   └── logparser/       # ログパース共通処理
├── workshop/            # 実装用: シンプルなTODOのみ
│   ├── phase1/
│   ├── phase2/
│   ├── phase3/
│   └── phase4/
├── solutions/           # 模範解答
│   ├── phase1/
│   ├── phase2/
│   ├── phase3/
│   └── phase4/
├── docs/                # ドキュメント
│   ├── HINTS.md         # レベル別ヒント集
│   └── SLIDES.md
└── logs/                # 生成されたログファイル
```


### ログファイルのサイズを変更したい

以下のオプションが使えます。

```bash
go run cmd/loggen/main.go --files=100 --lines=50000
```

##  ライセンス

MIT License
