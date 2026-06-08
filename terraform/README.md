# blog4 Terraform

`sacloud/sakura` provider (v3) で blog4 のインフラを管理する。

## 設計方針

- Provider: `sacloud/sakura` v3 に統一 (TiDB / AppRun 共用 / Container Registry / WebAccel / Object Storage が全部1プロバイダで揃う)
- 認証: 環境変数 (`SAKURACLOUD_ACCESS_TOKEN` / `SAKURACLOUD_ACCESS_TOKEN_SECRET`)。アクセストークンはレポジトリにストアしない
- パスワード等の機微情報: `password_wo` (write-only argument, Terraform 1.11+) と `TF_VAR_*` 環境変数で渡す。state に平文では入らない
- tfstate: さくらのオブジェクトストレージ (S3 互換) backend。`backend.tf` 参照。state ロックは S3 ネイティブの `use_lockfile`
- 移行方針の全体像: [../architecture/migration-plan.md](../architecture/migration-plan.md)

## このディレクトリで定義しているリソース

| ファイル | リソース |
|---|---|
| `versions.tf` | terraform / provider 版固定 |
| `backend.tf` | tfstate の Object Storage (S3 互換) backend |
| `provider.tf` | `sakura` provider 設定 |
| `variables.tf` | 入力変数 |
| `tidb.tf` | `sakura_ondemand_db` — TiDB CR インスタンス |
| `outputs.tf` | hostname など |

WebAccel / AppRun 共用 / Container Registry / Object Storage は後続 PR で追加 (既存リソースを import)。

## 使い方

### 1. mise で Terraform を入れる

```bash
mise install   # .mise.toml の terraform = "1.14.3" を入れる
```

### 2. 認証情報を環境変数で渡す

```bash
# さくらのクラウド API (provider 用)
export SAKURACLOUD_ACCESS_TOKEN='...'
export SAKURACLOUD_ACCESS_TOKEN_SECRET='...'
export TF_VAR_tidb_password='...'   # 16文字以上推奨

# オブジェクトストレージの S3 アクセスキー (backend 用)
export AWS_ACCESS_KEY_ID='...'
export AWS_SECRET_ACCESS_KEY='...'
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

```bash
cd terraform
terraform init                 # 新規。backend は最初から Object Storage
terraform plan
terraform apply
```

ローカル tfstate から移行する場合は init 時に state コピーを促される:

```bash
terraform init -migrate-state   # 既存 local state を Object Storage へコピー
```

移行が済んだらローカルの `terraform.tfstate*` は削除してよい。

`terraform output tidb_hostname` で接続先 (FQDN) が取れる。ポートは TiDB 固定の **4000**。

### パスワードのローテート

`var.tidb_password_version` を `1` → `2` に上げて、`TF_VAR_tidb_password` に新パスワードを入れて `apply`。リソース再作成なしで反映される。
