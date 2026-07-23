# blog4 Terraform

`sacloud/sakura` provider (v3) で blog4 のインフラを管理する。

## 設計方針

- Provider: `sacloud/sakura` v3 に統一 (TiDB / AppRun 共用 / Container Registry / WebAccel / Object Storage が全部1プロバイダで揃う)
- 認証: 環境変数。実値は **1Password (`blog4` vault) が正本**で、リポジトリには `op://` 参照 (`.env`) だけを置く。`op run` で実行時に注入し、平文をディスクに置かない
- パスワード等の機微情報: `password_wo` (write-only argument, Terraform 1.11+) と `TF_VAR_*` 環境変数で渡す。state に平文では入らない
- tfstate: さくらのオブジェクトストレージ (S3 互換) backend。`backend.tf` 参照。state ロックは S3 ネイティブの `use_lockfile`
- 移行方針の全体像: [../architecture/migration-plan.md](../architecture/migration-plan.md)
- DB 移行の実作業手順: [../architecture/migration-runbook.md](../architecture/migration-runbook.md)
  (件数確認・ダンプ・リストアは `scripts/db-count.sh` / `db-dump.sh` / `db-restore.sh`)

## このディレクトリで定義しているリソース

| ファイル | リソース |
|---|---|
| `versions.tf` | terraform / provider 版固定 |
| `backend.tf` | tfstate の Object Storage (S3 互換) backend |
| `.env` | 認証情報の `op://` 参照テンプレ (1Password CLI 用、実値は入れない) |
| `provider.tf` | `sakura` provider 設定 |
| `variables.tf` | 入力変数 |
| `tidb.tf` | `sakura_ondemand_db` — TiDB CR インスタンス |
| `outputs.tf` | hostname など |

WebAccel / AppRun 共用 / Container Registry / Object Storage は後続 PR で追加 (既存リソースを import)。

## 使い方

### 1. mise でツールを入れる

```bash
mise install   # terraform / awscli / 1password-cli (op) が入る
```

### 2. 認証情報 (1Password CLI で注入)

実値は 1Password の `blog4` vault が正本。リポジトリの `terraform/.env` は
`op://` 参照だけを持つ (commit 済み)。Terraform は `op run` 経由で実行し、
実値は実行時にメモリへ注入される (ディスクに平文を書かない)。

初回は 1Password の `blog4` vault に以下の item/field を作る:

| item | field | 中身 |
|---|---|---|
| `blog4-sakura-api` | `token` / `secret` | さくらのクラウド API キー |
| `blog4-object-storage` | `access-key-id` / `secret-access-key` | Object Storage S3 キー |
| `blog4-tidb` | `root-password` | TiDB CR root パスワード (16文字以上推奨) |

```bash
op signin
```

### 3. state バケットの用意 (初回のみ・済)

`backend.tf` が使う state バケット `blog4-tfstate` は鶏卵問題のため
Terraform 管理外。**さくらのクラウド コントロールパネルから private バケットとして作成済み**。
作り直す場合もコンパネから一度だけ手動で作る。

state 破損時の巻き戻し用に **versioning を有効化済み**。これはコンパネでは
設定できないため S3 互換 CLI で行う:

```bash
aws --endpoint-url=https://s3.isk01.sakurastorage.jp --region jp-north-1 \
    s3api put-bucket-versioning --bucket blog4-tfstate \
    --versioning-configuration Status=Enabled
```

### 4. init / plan / apply

`op run --env-file=.env` で各コマンドをラップして認証情報を注入する:

```bash
cd terraform
op run --env-file=.env -- terraform init    # 新規。backend は最初から Object Storage
op run --env-file=.env -- terraform plan
op run --env-file=.env -- terraform apply
```

ローカル tfstate から移行する場合は init 時に state コピーを促される:

```bash
op run --env-file=.env -- terraform init -migrate-state   # 既存 local state を Object Storage へコピー
```

移行が済んだらローカルの `terraform.tfstate*` は削除してよい。

`terraform output tidb_hostname` で接続先 (FQDN) が取れる。ポートは **3306**
(TiDB 標準の 4000 ではない。さくら側が 3306 で待ち受けており、4000 は閉じている)。

### パスワードのローテート

1Password の `blog4-tidb/root-password` を新パスワードに更新し、`.env` の
`TF_VAR_tidb_password_version` を `1` → `2` に上げて `op run --env-file=.env -- terraform apply`。
リソース再作成なしで反映される。
