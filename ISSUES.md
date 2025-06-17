# Go プロジェクト構造の改善点

このドキュメントは、blog4 プロジェクトの Go プロジェクトとしての構造的な問題点と改善案をまとめたものです。

## 1. パッケージ構造の問題

### 現状の問題点
- `main.go` がルートディレクトリに直接配置されている（通常は `cmd/blog4/main.go` に配置）
- `internal/` ディレクトリがない（プライベートパッケージの明確な分離がない）
- パッケージ名がディレクトリ構造と一致していない（例: `server/admin` ディレクトリ内のパッケージ名が `admin` のみ）

### 影響
- Go の標準的な慣習に従っていないため、他の開発者が理解しにくい
- パッケージの公開/非公開の境界が不明確

## 2. フロントエンドの配置

### 現状の問題点
- Svelte フロントエンドが `server/admin/frontend/` という深い階層に配置されている
- `package.json` と `package-lock.json` がプロジェクトルートにある（本来は frontend ディレクトリ内にあるべき）
- `node_modules` が Go のソースツリー内に存在する可能性

### 影響
- Go のツールチェーンがフロントエンドファイルを不要に処理する可能性
- ビルド時間の増加
- プロジェクト構造の複雑化

## 3. 生成コードの配置

### 現状の問題点
- SQLC で生成されたコードが `db/admin/admindb/` と `db/public/publicdb/` に分散
- データベースごとにディレクトリを分けるのは珍しい構造

### 影響
- インポートパスが長くなる
- 生成コードの管理が複雑

## 4. テストファイルの不足

### 現状の問題点
- テストファイルが2つしか見つからない（`server/date_test.go`、`server/admin/amazon_test.go`）
- Go プロジェクトとしてはテストカバレッジが低すぎる

### 影響
- コードの品質保証が不十分
- リファクタリング時のリスク増大

## 5. 設定ファイルの配置

### 現状の問題点
- 複数の設定ファイル（`app.json`、`app.jsonnet`、`deploy.json`）がルートに散在
- 専用の `configs/` ディレクトリがない

### 影響
- プロジェクトルートが散らかる
- 設定ファイルの管理が困難

## 6. 静的アセットの配置

### 現状の問題点
- `server/static/` と `server/templates/` が server ディレクトリ内にある
- 通常は `web/` や `assets/` などのトップレベルディレクトリに配置

### 影響
- アセットとサーバーコードの分離が不明確
- デプロイメント時の複雑性

## 7. 推奨されるプロジェクト構造

```
blog4/
├── cmd/
│   └── blog4/              # メインアプリケーション
│       └── main.go
├── internal/               # プライベートパッケージ
│   ├── admin/
│   │   ├── handler.go
│   │   └── service.go
│   ├── public/
│   ├── markdown/
│   ├── middleware/
│   └── worker/
├── pkg/                    # 公開パッケージ（必要な場合）
├── web/                    # フロントエンドとアセット
│   ├── admin/             # Svelte フロントエンド
│   │   ├── package.json
│   │   ├── src/
│   │   └── dist/
│   ├── static/
│   └── templates/
├── db/
│   ├── migrations/
│   ├── queries/           # SQLC クエリ
│   └── generated/         # SQLC 生成コード
├── api/                   # API定義
│   └── typespec/
├── configs/               # 設定ファイル
│   ├── app.jsonnet
│   └── deploy.json
├── scripts/               # ビルド・デプロイスクリプト
├── docs/                  # ドキュメント
├── build/                 # ビルド成果物
├── go.mod
├── go.sum
├── tools.go
├── Taskfile.yml
├── Dockerfile
├── .gitignore
└── README.md
```

## 8. その他の改善点

### コード品質
- エラーハンドリングで `os.Exit()` を使っている箇所がある（main.go 以外では避けるべき）
- グローバル変数や設定の管理が分散している
- ビルド成果物（`/blog4` バイナリ）のための明確な場所がない

### Go のベストプラクティス
- パッケージ名は短く、小文字で、アンダースコアを含まない
- インターフェースは利用側で定義する
- エラーは値として扱い、適切にラップする

## 実装の優先順位

1. **高優先度**
   - テストの追加
   - エラーハンドリングの改善（os.Exit の除去）

2. **中優先度**
   - `cmd/` ディレクトリの作成と main.go の移動
   - フロントエンドの再配置

3. **低優先度**
   - 全体的なディレクトリ構造の再編成
   - 生成コードの配置見直し

## 参考資料

- [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
- [Effective Go](https://golang.org/doc/effective_go.html)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)