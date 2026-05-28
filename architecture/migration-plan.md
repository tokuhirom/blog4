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
                     │ TLS+4000     │ HTTPS (S3 API)
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
- **`port` 属性なし** → TiDB の **4000 番**をアプリの `DATABASE_PORT` に明示

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
| 変わる (DB) | 接続先が TiDB の public ホスト名 + TLS + 4000 番に。`DATABASE_*` を更新 |
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
3. **接続設定**: `DATABASE_PORT=4000`、TLS 必須化、接続文字列の TLS パラメータ

### フェーズ 2: Terraform 化
4. `terraform/` を作り直す (前回 stash した雛形が流用可。AppRun 専有系は捨てる)
   - `sakura_ondemand_db` (TiDB)
   - `sakura_apprun_shared`
   - `sakura_webaccel` + `sakura_webaccel_certificate` (+ `_acl` / `_activation`)
   - `sakura_container_registry` (`terraform import` で既存取り込み)
   - secrets は sops + age (`WEBACCEL_GUARD` トークンもここで管理)
   - tfstate 置き場は別途決定 (前回は専用バケット。コスト要相談)
5. `terraform plan` で既存リソースとの差分を確認しながら import を進める

### フェーズ 3: データ移行 & カットオーバー
6. MariaDB から論理ダンプ (`mysqldump --no-tablespaces` 等) を取得
7. TiDB へリストア (FULLTEXT 行は schema から除外済みのものを使う)
8. ステージング的に新 DB を指す AppRun バージョンを用意し検証
9. `traffics` を新バージョンへ切替 (AppRun 共用型のトラフィック比率)
10. 旧 EDB を廃止

## 既知の注意点 / 未確認事項

1. ~~AppRun 共用型のカスタムドメイン~~ **(解決済み)**
   v3 `sakura_apprun_shared` にカスタムドメイン / TLS 属性は無いが、
   **`blog.64p.org` の入口と TLS 終端は前段の WebAccel が担う**ため問題なし。
   WebAccel は `sakura_webaccel` + `sakura_webaccel_certificate` で IaC 化でき、
   AppRun はオリジン (`public_url`) のままでよい。オリジン直アクセスは
   `X-WebAccel-Guard` (`WEBACCEL_GUARD`) で防止する既存の仕組みを維持する。
   - 残課題: WebAccel の TLS 証明書が Let's Encrypt 自動か持込か、現状の設定を確認
2. **TiDB のバージョンと FK の GA 状況** — EDB TiDB が v8.5 以上かで FK の
   挙動が変わる。要確認。
3. **AppRun 共用型の送信元 IP が不定** → `allowed_networks` での IP 制限が
   使いづらい。`0.0.0.0/0` + TLS + 強パスワードになる可能性。コネクション数
   上限 (`max_connections`) と合わせて要確認。
4. **Object Storage バケットの IaC** — provider にバケット作成リソースが
   あるか要確認。無ければ S3 互換 API or コンパネで管理し Terraform 外に。
5. **tfstate 置き場** — Sakura Object Storage 専用バケットは月 500 円。
   コスト懸念あり。ローカル + バックアップ / 既存バケット同居 等も含め再検討。

## 未決事項 (この計画の先で決める)

- フロント検索の実装方式 (全文 JSON のサイズ・キャッシュ戦略、初回ロード)
- tfstate 置き場の最終決定
- カスタムドメイン紐付けの扱い (上記注意点 1)
