# Workshop 実装ガイド

このディレクトリは、ワークショップ参加者が実際にコードを書いて学ぶための作業スペースです。

##  ワークショップの進め方

### 基本的な流れ

1. workshop/ ディレクトリで実装
2. 詰まったら docs/HINTS.md を参照
3. solutions/ で答え合わせ


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
    └── main.go    # Phase 4: さらなる高速化
```


##  使い方

### Phase 1: 逐次処理版

```bash
# リポジトリのルートから実行
make w1
```

### Phase 2: 並行処理版

```bash
# リポジトリのルートから実行
make w2
```


### Phase 3: ワーカープール版

```bash
# リポジトリのルートから実行
make w3
```


### Phase 4: さらなる高速化

```bash
# リポジトリのルートから実行
make w4
```

**利用可能な Make コマンド:**

- `make w1 w2 w3 w4` - Workshop Phase 1-4 を実行
- `make s1 s2 s3 s4` - Solution Phase 1-4 を実行
- `make gen` - ログファイルを生成
- `make help` - コマンド一覧を表示

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
- **[docs/SLIDES.md](../docs/SLIDES.md)** - スライド資料
- **[solutions/](../solutions/)** - 模範解答


##  パフォーマンス測定

各Phaseの処理時間を記録して、改善を可視化しましょう。

各コマンドを実行すると、実行時間が自動的に `workshop/results.txt` に記録されます。

```bash
make w1  # Phase 1を実行 → workshop/results.txtに記録
make w2  # Phase 2を実行 → workshop/results.txtに記録（Phase 1からの改善率も計算）
make w3  # Phase 3を実行 → workshop/results.txtに記録（Phase 1からの改善率も計算）
make w4  # Phase 4を実行 → workshop/results.txtに記録（Phase 1からの改善率も計算）
```

**結果ファイルの例:**
```
phase1=16.44
phase2=3.32 (phase1から4.95倍高速, 79.8%改善)
phase3=3.02 (phase1から5.44倍高速, 81.6%改善)
phase4=3.03 (phase1から5.43倍高速, 81.6%改善)
```

同様に、模範解答（solutions）を実行した場合は `solutions/results.txt` に記録されます。


※絶対的な処理時間はハードウェアに依存します。改善率（倍率）に注目してください。


頑張ってください！
