# blog4 インフラ 現状整理 (2026-05-28 時点)

このドキュメントは blog4 (`blog.64p.org`) の本番運用環境を把握するための
スナップショット。次の移行/再設計の議論の前提として使う。

## 全体図

```
                              Internet
                                 │
                                 ▼
                       ┌────────────────────┐
                       │ Route 53 (AWS)     │ ← blog.64p.org の権威 DNS
                       │ "64p.org" zone     │   (WebAccel エンドポイントへ CNAME)
                       └─────────┬──────────┘
                                 ▼
                       ┌────────────────────┐
                       │ WebAccel (CDN)     │ ← TLS 終端 / キャッシュ
                       │ (Sakura Cloud)     │   X-WebAccel-Guard を付与
                       └─────────┬──────────┘
                                 │ origin (X-WebAccel-Guard 必須)
                                 ▼
                       ┌────────────────────┐
                       │ AppRun 共用型       │ ← アプリ実行環境 (origin)
                       │ (Sakura Cloud)     │   /healthz は guard 免除
                       └──┬─────────────┬───┘
                          │             │
                  ┌───────▼──┐   ┌──────▼───────────────┐
                  │ EDB      │   │ Sakura Object        │
                  │ MariaDB  │   │ Storage              │
                  │ (Lab版)  │   │ blog4-attachments    │
                  │ ★廃止予告 │   │ blog4-backup         │
                  └──────────┘   │ (Ishikari 1)         │
                                 └──────────────────────┘

                            ── デプロイ ──
                                ▲
                                │ (PATCH apprun API)
                       ┌────────┴───────────┐
                       │ GitHub Actions     │
                       │ publish-image.yml  │
                       └────────┬───────────┘
                                │ docker push
                                ▼
                       ┌────────────────────┐
                       │ Sakura Container    │
                       │ Registry (.sakuracr │
                       │ .jp)                │
                       └─────────────────────┘
```

## コンポーネント一覧

| 役割 | サービス | 識別子 / 設定 | 備考 |
|---|---|---|---|
| CDN / TLS 終端 | さくらのウェブアクセラレータ (WebAccel) | オリジン = AppRun の公開 URL | **公開ドメイン `blog.64p.org` の入口**。TLS はここで終端 |
| アプリ実行 | さくらのクラウド AppRun (共用型) | `vars.APPRUN_APP_ID` | リッスンポート 8181 / `/healthz`。WebAccel の origin |
| データベース | さくらのクラウド エンハンスドDB (EDB) MariaDB (Lab版) | (環境変数で接続) | **★ 廃止予告あり、後述** |
| 添付ファイル ストレージ | さくらのオブジェクトストレージ | bucket `blog4-attachments` @ `s3.isk01.sakurastorage.jp` (jp-north-1) | 公開、CDN 配信用 |
| バックアップ ストレージ | さくらのオブジェクトストレージ | bucket `blog4-backup` @ `s3.isk01.sakurastorage.jp` (jp-north-1) | 非公開、アプリの `BACKUP_*` で利用 |
| コンテナレジストリ | さくらのクラウド コンテナレジストリ | `vars.SAKURA_REGISTRY_DOMAIN` (`*.sakuracr.jp`) | tag は git short hash |
| DNS | AWS Route 53 | zone `64p.org` | `blog` レコードは **WebAccel エンドポイント**へ |
| TLS | WebAccel で終端 | Let's Encrypt 自動 or 証明書持込 (現状どちらか要確認) | 証明書の出所は未確認 |

### オリジン保護 (WebAccel Guard)

公開トラフィックは必ず WebAccel を経由する。AppRun (origin) は
`internal/middleware/web_accel_guard.go` の `CheckWebAccelGuard` で
リクエストヘッダ `X-WebAccel-Guard` を検証し、環境変数 `WEBACCEL_GUARD` と
一致しないものを 400 で弾く = **オリジンへの直アクセスを防止**する。
ただし `/healthz` だけは guard を免除 (AppRun のヘルスチェックが WebAccel を
経由せず直接叩くため)。

## デプロイフロー

`.github/workflows/publish-image.yml` と `.github/actions/deploy-apprun/`:

1. `main` への push をトリガー
2. Docker buildx で Dockerfile (multi-stage、最終は `debian:trixie-slim`) をビルド
3. Sakura Container Registry へ `<registry>/blog4:<git-short-hash>` で push
4. `deploy.sh` が AppRun API を叩いて以下を実施:
   - 既存アプリ設定を GET (`apprun/.../applications/<APP_ID>`)
   - `jq` で `container_registry.image` を新タグに差し替え、
     `all_traffic_available = true` を追加、`container_registry.server` と
     `username` を削除 (PATCH 時の制約回避)
   - PATCH で新バージョンを反映
   - 反映後の `versions` を GET してログ表示

### デプロイで使う GitHub secrets / vars

| 種別 | 名前 | 用途 |
|---|---|---|
| vars | `SAKURA_REGISTRY_DOMAIN` | コンテナレジストリのホスト名 |
| vars | `SAKURA_REGISTRY_USERNAME` | レジストリ pull/push 用 |
| vars | `APPRUN_APP_ID` | デプロイ対象アプリの UUID |
| secret | `SAKURA_REGISTRY_PASSWORD` | レジストリ認証 |
| secret | `SACLOUD_API_TOKEN_ID` | AppRun API 用 (※環境変数名が不揃い: `SACLOUD_` と `SAKURA_API_*`) |
| secret | `SAKURA_API_TOKEN_SECRET` | 同上 |

## アプリの環境変数 (本番で要設定)

`internal/config.go` から抽出:

| 区分 | env 名 | デフォルト | メモ |
|---|---|---|---|
| アプリ | `BLOG_PORT` | `9191` (docker は 8181 override) | リッスンポート |
| DB | `DATABASE_USER` / `DATABASE_PASSWORD` / `DATABASE_HOST` / `DATABASE_PORT` / `DATABASE_NAME` | DB=`blog3` | エンハンスドDB を参照。`LOCAL_DEV` でなければ TLS で接続する |
| 管理 UI | `ADMIN_USER` / `ADMIN_PW` | user=`admin` | Basic 認証相当 |
| CORS | `ALLOWED_ORIGINS` | (empty) | カンマ区切り |
| 公開 URL | `SITE_BASE_URL` | `https://blog.64p.org` | |
| WebSub 通知 | `HUB_URLS` | 公的 hub × 2 | カンマ区切り |
| Amazon PA-API | `AMAZON_PAAPI5_ACCESS_KEY` / `_SECRET_KEY` | - | asin:... リンクで利用 |
| S3 添付 | `S3_ACCESS_KEY_ID` / `_SECRET_ACCESS_KEY` / `_REGION` / `_ATTACHMENTS_BUCKET_NAME` / `_ENDPOINT` / `_ATTACHMENTS_BASE_URL` | region=`jp-north-1` / endpoint=`s3.isk01.sakurastorage.jp` / bucket=`blog3-attachments` (default) | |
| S3 バックアップ | `S3_BACKUP_BUCKET_NAME` | `blog3-backup` | |
| バックアップ | `BACKUP_ENCRYPTION_KEY` | - | 暗号化キー |
| WebAccel | `WEBACCEL_GUARD` | - | キャッシュ無効化トークン |
| タイムゾーン | `TIMEZONE_OFFSET` | `32400` (JST) | |
| OG 画像 | `OG_IMAGE_ENABLED` / `OG_IMAGE_FONT_PATH` | true / `/usr/share/fonts/opentype/ipafont-gothic/ipagp.ttf` | コンテナ内で Puppeteer がフォント参照 |

ローカル開発: `docker-compose.yml` が MariaDB 10.11.17 + LocalStack で
S3/DB を再現。本番との差は DB エンドポイント / 認証情報のみ。

## 課題

### 1. ★ 緊急: エンハンスドDB (EDB) の MariaDB が提供終了予定

[【予告】エンハンスドデータベース(Lab)の提供形態の変更および一部提供終了のお知らせ](https://cloud.sakura.ad.jp/news/2026/04/02/enhanced-database-lab-end/)
(2026-04-02 公開) によれば:

- **2026-04-09**: MariaDB 新規申込終了
- **2026-04-23**: 提供形態変更実施
- **MariaDB は提供終了**。既存ユーザへの完全終了日は記事中明記なし
- TiDB は Lab → CR 版に継続。MariaDB ユーザの移行先指示は記載なし

blog4 はこの MariaDB に依存しているため、**期日内に別 DB へ移行する必要がある**。

候補:

| 候補 | プラス | マイナス |
|---|---|---|
| データベースアプライアンス (MariaDB 10.11) | マネージド、自動バックアップ、MariaDB 継続 | **vswitch 必須**で AppRun **共用型**からは直接接続不可 → 経路の橋渡しか AppRun 専有型への移行が必要 |
| 専有サーバ + 自前 MariaDB | 安い (月1.5k程度から)、構成自由 | バックアップ/アップデート/HA を自分で持つ |
| TiDB CR 版 | Sakura 内で継続提供 | MySQL 互換だが MariaDB 固有挙動の検証が必要、Lab 卒業直後で枯れていない可能性 |
| 外部 DBaaS (PlanetScale/Neon/Supabase 等) | 運用ほぼゼロ、無料枠あり | レイテンシ、撤退リスク、海外データ規約 |

### 2. インフラが IaC で管理されていない

現状、AppRun アプリの設定・コンテナレジストリ・オブジェクトストレージ
バケット・DB アプライアンスはすべて **コントロールパネル / 手動 API
コール** で作成・管理されている。

問題:
- 構成変更の履歴が残らない (誰が何をいつ変えたか追えない)
- 再現性がない (環境の作り直し / staging 環境の複製ができない)
- レビューが効かない (本番変更が口頭/手元判断で進む)
- 災害復旧時の手順が暗黙知になっている

対応方針 (案): `sacloud/sakura` (Terraform Provider) で `terraform/`
配下に集約する。1人運用でも履歴 + コードレビュー = 自己レビューの効用は大きい。
※ 過去に一度試作したが「共用型 AppRun と DB アプライアンス (= 別サービス
の DBA) の経路問題」「EDB の廃止予告」の整理が先と判断して一旦保留中
(`git stash`)。

### 3. シークレット/認証情報の置き場が散らばっている

- GitHub Actions secrets: `SAKURA_REGISTRY_PASSWORD`, `SACLOUD_API_TOKEN_ID`, `SAKURA_API_TOKEN_SECRET`
- AppRun の環境変数 (DB password, ADMIN_PW, S3 access key, AMAZON_PAAPI5_*, BACKUP_ENCRYPTION_KEY 等) はコンパネ直接入力と推測
- ローテーション手順が文書化されていない

IaC 化と合わせて sops + age 等で「正本はリポジトリ、復号鍵だけ別」に
まとめると一元管理しやすい。

### 4. デプロイの観測性

`deploy.sh` は `set -ex` で実行ログを出すが、失敗時の rollback 手順や
ヘルスチェック後の自動切り戻しはない。AppRun の `all_traffic_available
= true` で即時 100% 切替なので、起動失敗時にはダウンタイムが発生し得る。

### 5. 環境変数のデフォルトと実際の本番値の差

`config.go` のデフォルトには旧 `blog3-attachments` / `blog3-backup` などの
名前が残っている。本番では override されているはずだが、コードと運用の
ズレは事故のもと。IaC 化のタイミングで揃えたい。

## 次の検討トピック (本ドキュメントの先)

- DB 移行先の決定 (上記候補から)。これが決まらないと AppRun を専有型に
  するか共用型のままかの判断もブロックされる
- 決定後、IaC (Terraform) でカバーする範囲の確定
- シークレットローテーション運用の整備
