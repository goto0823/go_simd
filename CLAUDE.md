# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## このリポジトリの目的と進め方

Go 初心者のユーザーが、SIMD を題材に低レベル（生成アセンブリ、CPU キャッシュ、メモリ）を学ぶための学習用リポジトリ。

- **Claude はコードを書かない・修正しない。** コードはすべてユーザーが書く。Claude の役割は、概念の説明、仕様とヒント（API 名・イディオム・コマンド）の提示、書かれたコードのレビューと問題点の指摘。
- 「直したら次に何のエラーが出るか」などの確認は、ファイルをスクラッチ領域にコピーして試す。ユーザーのファイルには触れない。
- 説明は日本語で、コンパイラのメッセージや実測値など、このマシンで確かめた根拠を添える。
- 学習の計画と現在地は下の「計画」の節を参照。

## 計画

ステップが終わったら、この節の状態を更新する。

| # | 内容 | 状態 |
|---|---|---|
| 1 | スカラ版 `SumScalar`、テーブル駆動テスト、サイズ別ベンチ | 完了 (2026-09-10) |
| 2 | `simd` パッケージで SIMD 版を書き、スカラ版と比べる | **次** |
| 3 | 端数処理 (`LoadFloat32sPart` / `StorePart`) | |
| 4 | マスク（条件付きの集計、バイト検索など） | |
| 5 | `simd/archsimd` で AVX2 の型を直接使う、`ToArch()` | |
| 6 | 生成アセンブリを読み、スカラ版 (`ADDSS`) からの変化を見る | |
| 7 | FMA、メモリ帯域、アラインメント、アキュムレータを分けて依存チェーンを短くする | |
| 8 | 実用的な題材（内積、UTF-8 検証、JSON パーサなど） | |

### ステップ2 の予定

- ベンチのテーブルに実装そのもの（`func([]float32) float32` の値）を入れ、「実装 × サイズ」で比べる。
- ステップ1の結果をファイルに残し、`benchstat` で比較する。
- 予想: 小さいサイズでは大きく速くなるが、1e7 要素ではメモリ側が律速になって伸びが鈍る。

### スカラ版の基準値（2026-09-10、開発機）

1 要素あたり、1e3〜1e6 要素で約 0.86 ns、1e7 要素で約 0.99 ns。
解釈: 1e6 までは `ADDSS` のループ伝搬依存チェーン（1 回約 4 サイクル）が律速。1e7 で遅くなるのは、4KB ページ境界ごとにプリフェッチが途切れて DRAM を読みに行くことと TLB ミスによるという仮説（未検証）。

### 寄り道の候補

- `int32` 版の総和: 整数の加算は約 1 サイクルなので、スカラのままでも約 4 倍速くなるはず。
- ランダムな順番で読む総和: プリフェッチが効かないときのキャッシュミスの影響を見る。
- 2MB の huge page を使って 1e7 要素のケースが変わるか見る。
- `perf` でキャッシュミスを数えて上の仮説を確かめる（WSL では制約あり）。

### 未対応

- `.gitignore`（`.direnv/`, `.go/`）がまだない。
- `flake.nix` / `flake.lock` が git に add されていないため、`nix develop` は `path:` 付きでないと失敗する。

## 開発環境

- Nix flake + direnv (`.envrc` は `use flake`)。`flake.nix` が Go 1.27.0・gopls・gotools・go-tools・delve・air を提供し、`GOEXPERIMENT=simd` と `GOPATH=$PWD/.go` を設定する。
- 標準ライブラリの `simd` / `simd/archsimd` は `//go:build goexperiment.simd` 付きなので、`GOEXPERIMENT=simd` がないと import できない。
- Claude の Bash は direnv の環境に入らず、システムにある別バージョンの Go が使われる。Go のコマンドは dev shell 経由で実行する:

  ```bash
  nix develop "path:$PWD" -c go test ./...
  ```

  `path:` を付けるのは、git リポジトリ内の flake は git が追跡しているファイルしか見ず、`flake.nix` が未追跡だと `Path 'flake.nix' ... is not tracked by Git` で失敗するため。shellHook が毎回 `go version` を表示する。
- `.direnv/` と `.go/` はローカルで生成されるディレクトリ。

## よく使うコマンド（dev shell 内で実行）

```bash
go test -v ./...                                          # テスト
go test -v -run 'TestSumScalar/nil' ./...                 # サブテストを1つだけ
go test -run='^$' -bench=. -benchmem -count=3             # ベンチのみ（テストを飛ばす）
go test -run='^$' -bench='SumScalar/10000000$' -benchmem  # 1サイズだけ
go vet ./...
gofmt -l .
go build -gcflags='-S' sum.go 2>&1                        # 生成アセンブリ（出力は stderr）
go doc simd                                               # simd パッケージのドキュメント
```

## 構成

- ルートにある単一パッケージ `simdlab`。
- 実装は `func([]float32) float32` のシグネチャで揃える。`sum.go` の `SumScalar` がスカラの基準実装で、SIMD 版を同じテーブル駆動テストとベンチで比較していく。
- ベンチのサイズ（1e3〜1e7 要素、float32 は 4 バイト）は開発機のキャッシュ階層をまたぐように選んでいる。開発機は i7-10700K: L1d 32KB・L2 256KB（コアごと）、L3 16MB（全コア共有）、ページは 4KB。AVX2 まで対応で AVX-512 はないので、ベクタ幅は最大 256 ビット。
- `simd` パッケージはベクタ幅を実行時の CPU 機能検出で決める（Go ソースの `src/simd/midway_amd64.go`）ので、この開発機では 256 ビットになる。

## 注意点

- `b.Loop()` のループは 1 つの `*testing.B` につき 1 回だけ。2 回目は `B.Loop called with timer stopped` で失敗するので、サイズごとに `b.Run` で分ける。
- SIMD 版では加算の順序が変わるため、浮動小数点の結果はスカラ版とビット単位では一致しない。実装どうしを比べるテストは許容誤差で比較する。
- 1e7 要素のベンチは反復回数が約 100 回と少なく、L3 が全コア共有なこともあって結果がブレやすい。比べるときは `-count` を増やす。
