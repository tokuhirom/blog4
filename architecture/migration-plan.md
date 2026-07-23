# blog4 インフラ移行計画 (EDB MariaDB → TiDB + Terraform IaC 化)

現状整理は [current-state.md](./current-state.md) を参照。本書は移行の
ターゲット構成と作業計画。

## ゴール

1. 廃止予定の EDB (Lab) MariaDB から、さくらの **オンデマンドデータベース
   TiDB (CR版)** へ移行する
2. インフラを **Terraform (`sacloud/sakura` v3) で IaC 化** する
3. 上記に伴い壊れる**全文検索をフロントエンド検索へ作り替える**
4. AppRun は **共用型のまま**、追加の固定費を増やさない

## ターゲット構成

```
                         Internet
                            │
                            ▼
                  ┌────────────────────┐
                  │ Route 53 (AWS)     │  blog.64p.org → WebAccel へ CNAME
                  └─────────┬──────────┘
                            ▼
                  ┌────────────────────┐
                  │ WebAccel (CDN)     │  sakura_webaccel (+ _certificate)
                  │ TLS 終端 / cache   │  blog.64p.org の TLS はここ
                  └─────────┬──────────┘
                            │ origin (X-WebAccel-Guard 付与)
                            ▼
                  ┌────────────────────┐
                  │ AppRun 共用型       │  sakura_apprun_shared
                  │ (Sakura Cloud)     │  image = container registry
                  └──┬─────────────┬───┘
                     │ TLS+3306     │ HTTPS (S3 API)
            ┌────────▼───────┐  ┌──▼──────────────┐
            │ オンデマンドDB  │  │ Object Storage  │
            │ TiDB (CR/無償) │  │ attachments     │
            │ public + TLS   │  │ backup          │
            └────────────────┘  └─────────────────┘
                     ▲
                     │ (pull)
            ┌────────┴────────┐
            │ Container        │  sakura_container_registry
            │ Registry         │  *.sakuracr.jp
            └──────────────────┘
```

## Terraform Provider 決定

**`sacloud/sakura` v3 に一本化**する。理由: 必要なリソースがすべて v3 に
そろっており、v2 (`sacloud/sakuracloud`) と混在させずに済む。

| 用途 | v3 リソース | 備考 |
|---|---|---|
| DB (TiDB) | `sakura_ondemand_db` | `enhanced_db` は deprecated。`database_type = "tidb"` 固定 |
| アプリ実行 | `sakura_apprun_shared` | `components[].deploy_source.container_registry.image` |
| CDN / TLS | `sakura_webaccel` + `sakura_webaccel_certificate` (+ `_acl` / `_activation`) | `blog.64p.org` の入口・TLS 終端をここで表現 |
| コンテナレジストリ | `sakura_container_registry` | 既存を import |
| (任意) Object Storage | ※後述 | バケットは provider 対応状況を要確認 |

## 主要リソースのスキーマ要点

### `sakura_ondemand_db` (TiDB)
- Required: `name`, `database_name`, `database_type = "tidb"`, `region` (`is1`/`tk1`), `password_wo` (write-only, TF 1.11+)
- Optional: `allowed_networks` (CIDR の IP 制限), `password_wo_version`
- Read-Only: `hostname` (`<database_name>.tidb-is1.db.sakurausercontent.com`), `max_connections`
- **`port` 属性なし** → アプリの `DATABASE_PORT` には **3306** を明示する。
  TiDB 標準の 4000 ではない (2026-07-23 に疎通確認。4000 は接続できない)

### `sakura_apprun_shared`
- Required: `name`, `components`, `min_scale`, `max_scale`, `port`, `timeout_seconds`
- `components[]`: `name`, `max_cpu` (`0.5`/`1`/`2`), `max_memory` (`1Gi`/`2Gi`/`4Gi`), `deploy_source`, `env`, `probe`
- `deploy_source.container_registry`: `image`, `server`, `username`, `password_wo`
- `env[]`: `{ key, value }` (value は常に Sensitive。secret フラグは無い)
- `traffics[]`: `{ version_index, percent }`
- Read-Only: `public_url`

## 移行で「変わる点 / 変わらない点」

| | 内容 |
|---|---|
| 変わる (DB) | 接続先が TiDB の public ホスト名 + TLS に。`DATABASE_*` を更新 (ポートは 3306 のまま) |
| 変わる (検索) | `MATCH/AGAINST` 撤廃。全文検索はフロントエンドで実施 |
| 変わる (スキーマ) | `FULLTEXT KEY idx_bigram` 削除、`ENGINE=InnoDB` 句削除 |
| 変わる (運用) | 構成が Terraform 管理に。手動コンパネ操作を廃する |
| 変わらない | AppRun 共用型、**WebAccel (前段 CDN/TLS) + `X-WebAccel-Guard` ガード**、Object Storage、Container Registry、Route 53、Dockerfile |

## 作業計画

### フェーズ 1: アプリ改修 (TiDB 互換 + フロント検索)
DB 移行の前提。複数 PR に分ける。

1. **全文検索のフロント化**
   - public: 公開エントリの検索用データ (title / 抜粋 or 本文 / path) を
     JSON で配信する API を追加 → クライアント JS で絞り込み
   - admin: 認証済み前提で全エントリ (private 含む) を配信 → Preact 側で検索
   - `AdminFullTextSearchEntries` / public の `MATCH...AGAINST` クエリを撤去
2. **スキーマの TiDB 互換化**
   - `FULLTEXT KEY idx_bigram (title, body)` を削除
   - `ENGINE=InnoDB` 句を削除 (TiDB は無視するが明示削除でクリーンに)
   - `FOREIGN KEY ... ON DELETE CASCADE` の挙動を TiDB で検証
     (EDB TiDB の FK は v8.5 で GA。バージョン次第。最悪アプリ側で CASCADE 相当)
   - `SELECT ... FOR UPDATE` (visibility.sql) の悲観ロック挙動を検証
3. **接続設定**: `DATABASE_PORT=3306`、TLS 必須化、接続文字列の TLS パラメータ

### フェーズ 2: Terraform 化
4. `terraform/` を作り直す (前回 stash した雛形が流用可。AppRun 専有系は捨てる)
   - `sakura_ondemand_db` (TiDB)
   - `sakura_apprun_shared`
   - `sakura_webaccel` + `sakura_webaccel_certificate` (+ `_acl` / `_activation`)
   - `sakura_container_registry` (`terraform import` で既存取り込み)
   - secrets は sops + age (`WEBACCEL_GUARD` トークンもここで管理)
   - tfstate 置き場は別途決定 (前回は専用バケット。コスト要相談)
5. `terraform plan` で既存リソースとの差分を確認しながら import を進める

### フェーズ 3: データ移行 & カットオーバー — **2026-07-24 完了**

実作業の手順とコマンドは [migration-runbook.md](./migration-runbook.md) に切り出した
(`scripts/db-count.sh` / `db-dump.sh` / `db-restore.sh` を使う)。

6. MariaDB から論理ダンプ (`mysqldump --no-tablespaces` 等) を取得 — 済
7. TiDB へリストア (FULLTEXT 行は schema から除外済みのものを使う) — 済
8. 1Password の `blog4-app-db` を TiDB の値に更新し `terraform apply` で
   AppRun の env を差し替える — 済 (in-place update。`traffics` を使った
   段階切替はしていない。共用型は min/max scale = 1 で、旧 EDB を残しておけば
   1Password を戻して apply するだけで切り戻せるため)
9. 旧 EDB を廃止 — **未実施**。数日運用して問題がなければ廃止する

## 既知の注意点 / 未確認事項

解決済み (2026-07-23〜24 の移行作業で確認):

1. **TiDB のバージョンと FK の GA 状況** — さくらの TiDB CR は **v8.5.0**。
   FK は GA で、`ON DELETE CASCADE` 付きの FK をそのまま作成できた。
2. **AppRun 共用型の送信元 IP が不定** — 問題にならなかった。TiDB CR の
   `allowed_networks` は未設定 = 制限なしで、送信元を絞る必要がない。
   `max_connections` は 50 で、アプリのプール + 手元の作業クライアントでも余裕がある。
3. **WebAccel の TLS 証明書** — Let's Encrypt 自動運用。証明書リソースは管理しない
   (`terraform/webaccel.tf` 参照)。

残っているもの:

4. **Object Storage バケットの IaC** — provider にバケット作成リソースが
   あるか要確認。無ければ S3 互換 API or コンパネで管理し Terraform 外に。
   なお WebAccel サイト (バケット origin 含む) は v3.12.2 以降 import できるようになり、
   `blog-attachments.64p.org` も Terraform 管理下に入れた。
5. **tfstate 置き場** — Sakura Object Storage 専用バケット (`blog4-tfstate`) で運用中。
   月 500 円のコスト懸念は残るが、当面はこのまま。

## 未決事項 (この計画の先で決める)

- フロント検索の実装方式 (全文 JSON のサイズ・キャッシュ戦略、初回ロード)
- tfstate 置き場の最終決定
