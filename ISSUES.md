# Go プロジェクト構造の改善点

このドキュメントは、blog4 プロジェクトの Go プロジェクトとしての構造的な問題点と改善案をまとめたものです。

## 最近の改善 (2025-01)

以下の改善が実施されました：
- ✅ `main.go` を `cmd/blog4/main.go` に移動（[PR #265](https://github.com/tokuhirom/blog4/pull/265)）
- ✅ `os.Exit()` の使用を main.go のみに限定（[PR #264](https://github.com/tokuhirom/blog4/pull/264)）
- ✅ フロントエンドビルドツールを OpenAPI Generator から Orval に移行（Java 依存を削除）
- ✅ パッケージ名とディレクトリ構造の一致を確認（すべてのパッケージが正しく命名されていることを確認済み）
- ✅ `internal/` ディレクトリを作成し、プライベートパッケージを整理（admin、markdown、middleware を internal/ へ移動）

## 1. パッケージ構造の問題

### 現状の問題点
- ~~`main.go` がルートディレクトリに直接配置されている~~ **✅ 解決済み**
- ~~`internal/` ディレクトリがない（プライベートパッケージの明確な分離がない）~~ **✅ 解決済み** - internal/ ディレクトリを作成し、admin、markdown、middleware パッケージを移動
- ~~パッケージ名がディレクトリ構造と一致していない（例: `server/admin` ディレクトリ内のパッケージ名が `admin` のみ）~~ **✅ 解決済み** - すべてのパッケージが Go の慣習に従って正しく命名されている

### 影響
- Go の標準的な慣習に従っていないため、他の開発者が理解しにくい
- パッケージの公開/非公開の境界が不明確

## 2. フロントエンドの配置

### 現状の問題点
- Svelte フロントエンドが `internal/admin/frontend/` という深い階層に配置されている
- ~~`package.json` と `package-lock.json` がプロジェクトルートにある~~ **✅ 解決済み** - frontend ディレクトリ内に正しく配置されている
- `node_modules` が Go のソースツリー内に存在する（ただし .gitignore で除外済み）

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
- ~~テストファイルが2つしか見つからない（`server/date_test.go`、`server/admin/amazon_test.go`）~~ **一部改善** - 基本的なテストを追加
- Go プロジェクトとしてはテストカバレッジが低すぎる
- データベース依存のコードが多く、モックなしではテストが困難

### 影響
- コードの品質保証が不十分
- リファクタリング時のリスク増大

### テストが困難な部分（リファクタリング後に対応予定）

#### 1. データベース依存のコード
- `server/entry_image_service.go` - `getImageFromEntry` メソッドがデータベースクエリに直接依存
- `internal/admin/admin_api_service.go` - すべてのメソッドがデータベースに直接アクセス
- `internal/admin/twohop.go` - `getLinkPalletData` と `getEntriesByLinkedTitles` がデータベースに依存
- `internal/markdown/asin.go` と `wiki_link.go` のレンダラー - データベースクエリを実行

**解決策**: インターフェースを導入してモック可能にする、またはリポジトリパターンを採用

#### 2. 外部サービス依存
- `internal/admin/amazon_paapi5.go` - Amazon API への直接依存
- `server/backup.go` - S3 サービスへの直接依存
- `server/keep_alive.go` - HTTP クライアントの直接使用
- `internal/admin/pubsubhubbub.go` - 外部 Hub への HTTP POST

**解決策**: HTTPクライアントをインターフェース化、外部サービスクライアントの抽象化

#### 3. ファイルシステム依存
- `internal/admin/admin.go` - 静的ファイルの直接読み込み
- `server/public.go` - テンプレートと静的ファイルの embed

**解決策**: ファイルシステムアクセスの抽象化

#### 4. 時間依存
- `server/backup.go` - 現在時刻に依存したバックアップファイル名生成

**解決策**: 時刻を注入可能にする

#### 5. サードパーティライブラリのインターフェース
- `internal/markdown/wiki_link.go` と `asin.go` - Goldmark の `util.BufWriter` インターフェースを実装するモックの作成が複雑

**解決策**: 統合テストの活用、またはレンダリングロジックの抽象化

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

## 7. 現在のプロジェクト構造

```
blog4/
├── cmd/
│   └── blog4/              # メインアプリケーション ✅
│       └── main.go
├── server/                 # サーバーコード
│   ├── admin/
│   ├── router/
│   └── ...
├── db/                     # データベース関連
├── typespec/               # API定義
└── ...
```

## 8. 推奨される最終的なプロジェクト構造

```
blog4/
├── cmd/
│   └── blog4/              # メインアプリケーション ✅
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

## 9. その他の改善点

### コード品質
- ~~エラーハンドリングで `os.Exit()` を使っている箇所がある~~ **✅ 解決済み** - main.go のみに限定
- グローバル変数や設定の管理が分散している
- ビルド成果物（`/blog4` バイナリ）のための明確な場所がない

### Go のベストプラクティス
- パッケージ名は短く、小文字で、アンダースコアを含まない
- インターフェースは利用側で定義する
- エラーは値として扱い、適切にラップする

## 実装の優先順位

1. **高優先度**
   - テストの追加
   - ~~エラーハンドリングの改善（os.Exit の除去）~~ **✅ 解決済み**

2. **中優先度**
   - ~~`cmd/` ディレクトリの作成と main.go の移動~~ **✅ 解決済み**
   - フロントエンドの再配置（`web/` ディレクトリへの移動）
   - ~~`internal/` ディレクトリの作成とプライベートパッケージの整理~~ **✅ 解決済み**

3. **低優先度**
   - 全体的なディレクトリ構造の再編成
   - 生成コードの配置見直し
   - 設定ファイルの `configs/` ディレクトリへの集約

## 参考資料

- [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
- [Effective Go](https://golang.org/doc/effective_go.html)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)