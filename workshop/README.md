# Workshop 実装ガイド

このディレクトリは、ワークショップ参加者が実際にコードを書いて学ぶための作業スペースです。

##  ワークショップの進め方

### 基本的な流れ

1. **workshop/ ディレクトリで実装** - Phase 1から順に実装
2. **詰まったら docs/HINTS.md を参照** - ヒント集
3. **solutions/ で答え合わせ** - 模範解答と比較


---

## 前提条件

ワークショップを始める前に、[docs/SETUP_CHECK.md](../docs/SETUP_CHECK.md) のチェックリストを完了してください。

特に以下を確認：
- Go 1.25以降がインストールされている
- ログファイル（200個）が生成済み

##  ディレクトリ構成

```
workshop/
├── phase1/
│   └── main.go    # Phase 1: 逐次処理版を実装
├── phase2/
│   └── main.go    # Phase 2: 並行処理版を実装
├── phase3/
│   └── main.go    # Phase 3: ワーカープール版を実装
└── phase4/
    └── main.go    # Phase 4: さらなる高速化に挑戦
```


##  使い方

### Phase 1: 逐次処理版

```bash
# リポジトリのルートから実行
go run ./workshop/phase1/main.go
```

### Phase 2: 並行処理版

```bash
# リポジトリのルートから実行
go run ./workshop/phase2/main.go
```


### Phase 3: ワーカープール版

```bash
# リポジトリのルートから実行
go run ./workshop/phase3/main.go
```


### Phase 4: さらなる高速化

```bash
# リポジトリのルートから実行
go run ./workshop/phase4/main.go
```

Phase 3よりもさらに高速化する。あらゆる最適化手法を試してください。

---

##  模範解答

各Phaseの模範解答は `solutions/` ディレクトリにあります。

- [solutions/phase1/main.go](../solutions/phase1/main.go) - 逐次処理版
- [solutions/phase2/main.go](../solutions/phase2/main.go) - 基本並行処理版
- [solutions/phase3/main.go](../solutions/phase3/main.go) - ワーカープール版
- [solutions/phase4/main.go](../solutions/phase4/main.go) - さらなる最適化版

**模範解答を見るタイミング:**
1. 自分で実装を試してから
2. ヒントを見ても解決できない場合
3. 実装完了後、他のアプローチを確認したい場合


---

##  ヒントが必要な場合

詰まったときは以下を参照してください。

- **[docs/HINTS.md](../docs/HINTS.md)** - ヒント
- **[solutions/](../solutions/)** - 模範解答


##  パフォーマンス測定

各Phaseの処理時間を記録して、改善を可視化しましょう。

### 測定シート

```
Phase 1: _____ 秒（基準値）
Phase 2: _____ 秒（改善率: Phase 1の _____倍）
Phase 3: _____ 秒（改善率: Phase 1の _____倍）
Phase 4: _____ 秒（改善率: Phase 1の _____倍）
```

※絶対的な処理時間はハードウェアに依存します。改善率（倍率）に注目してください。

---

頑張ってください！
