# blog4 Terraform

`sacloud/sakura` provider (v3) で blog4 のインフラを管理する。

## 設計方針

- Provider: `sacloud/sakura` v3 に統一 (TiDB / AppRun 共用 / Container Registry / WebAccel / Object Storage が全部1プロバイダで揃う)
- 認証: 環境変数 (`SAKURACLOUD_ACCESS_TOKEN` / `SAKURACLOUD_ACCESS_TOKEN_SECRET`)。アクセストークンはレポジトリにストアしない
- パスワード等の機微情報: `password_wo` (write-only argument, Terraform 1.11+) と `TF_VAR_*` 環境変数で渡す。state に平文では入らない
- tfstate: 現状はローカル。後続 PR で Object Storage backend に切り替える予定
- 移行方針の全体像: [../architecture/migration-plan.md](../architecture/migration-plan.md)

## このディレクトリで定義しているリソース

| ファイル | リソース |
|---|---|
| `versions.tf` | terraform / provider 版固定 |
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
export SAKURACLOUD_ACCESS_TOKEN='...'
export SAKURACLOUD_ACCESS_TOKEN_SECRET='...'
export TF_VAR_tidb_password='...'   # 16文字以上推奨
```

### 3. init / plan / apply

```bash
cd terraform
terraform init
terraform plan
terraform apply
```

`terraform output tidb_hostname` で接続先 (FQDN) が取れる。ポートは TiDB 固定の **4000**。

### パスワードのローテート

`var.tidb_password_version` を `1` → `2` に上げて、`TF_VAR_tidb_password` に新パスワードを入れて `apply`。リソース再作成なしで反映される。
