# blog4 DB 移行 runbook (EDB MariaDB → TiDB CR)

[migration-plan.md](./migration-plan.md) のフェーズ3 (データ移行 & カットオーバー) の実作業手順。
想定する流れは **両方カウント → 両方バックアップ → MariaDB から TiDB へ移す → 両方カウント**。

## 使うスクリプト

| スクリプト | 役割 |
|---|---|
| `scripts/db-count.sh [tidb\|mariadb\|both]` | テーブルごとの件数 (記事数、公開/非公開の内訳) を表示。`both` は差分に `*` を付けて並べる |
| `scripts/db-dump.sh [options] [tidb\|mariadb\|both]` | mysqldump して `/tmp` に置く |
| `scripts/db-restore.sh [--target tidb\|mariadb] <file>` | SQL ファイルを流し込む |

共通事項:

- **必ず `op run --env-file=terraform/.env --` 経由で実行する。** 接続情報は 1Password が正本
  (TiDB = `blog4-tidb`、MariaDB = `blog4-app-db`)
- mysql クライアントはローカルに入れず docker で動かす (TiDB へは `mysql:8.4`、
  旧 EDB へは `mariadb:10.11.18` のクライアント)
- パスワードは 600 の一時ファイル経由でクライアントに渡す。`ps` や `docker inspect` には出ない
- TiDB のホスト名は `terraform output tidb_hostname` から取る。terraform を通したくないときは
  `TIDB_HOST=... TIDB_USER=...` を環境変数で明示すれば terraform なしで動く

TiDB の接続先は 2026-07-23 に疎通確認済み:

- ユーザ名は **データベース名と同じ** (`blog5`)
- ポートは **3306**。TiDB 標準の 4000 は閉じている (さくら側が 3306 で待ち受けている)
- TLS 必須 (`ssl-mode=REQUIRED` で接続できる)
- 送信元ネットワーク制限 (`allowed_networks`) は未設定 = 制限なし。手元から素で繋がる

## 移行前に片付ける必要があるアプリ側の課題

対応済み:

- **TLS 接続** — `cmd/blog4/main.go` で `LOCAL_DEV` 以外は `mysql.Config.TLSConfig` を
  有効にした。TiDB CR は TLS 必須、ローカルの docker-compose MariaDB は TLS 非対応なので分岐する
- **DB 名の環境変数名** — アプリ (`internal/config.go`) / `docker-compose.yml` /
  `terraform/apprun.tf` を `DATABASE_NAME` に統一した。旧名は `DATABASE_DB`

未対応:

1. **アプリ内の日次バックアップが `mariadb-dump`** — `internal/backup.go:43`。TiDB 相手に
   動くか確認が要る (認証プラグインと TLS)。動かないなら `mysqldump` に替える。
2. **FOREIGN KEY の対応状況** — `db/init/01-schema.sql` は `entry_image` / `entry_link` に
   FK を張っている。TiDB の FK は v8.5 で GA。手順4 (スキーマ作成) が通れば対応済みと判断でき、
   落ちるならスキーマから FK を外し、削除時の CASCADE をアプリ側で持つ。

## 手順

以降のコマンドはリポジトリルートで実行する。

### 0. 疎通確認

```bash
op signin
op run --env-file=terraform/.env -- ./scripts/db-count.sh mariadb
op run --env-file=terraform/.env -- ./scripts/db-count.sh tidb
```

TiDB 側は移行前なので「テーブルが1つもない」と出るのが正常。ここで認証エラーが出るなら
ユーザ名/パスワードの仮定を見直す (上記の `TIDB_USER` 注記を参照)。

### 1. 移行前のカウントを記録する

```bash
op run --env-file=terraform/.env -- ./scripts/db-count.sh both | tee /tmp/blog4-count-before.txt
```

この出力が移行の正解値になる。手順6 で突き合わせる。

### 2. MariaDB のバックアップを取る

```bash
op run --env-file=terraform/.env -- ./scripts/db-dump.sh --gzip mariadb
# -> /tmp/blog4-mariadb-<timestamp>.sql.gz (スキーマ + データのフルダンプ)
```

これは**切り戻し用の保険**。移行に使うダンプ (手順5) とは別物なので、消さずに取っておく。

### 3. TiDB のバックアップを取る

```bash
op run --env-file=terraform/.env -- ./scripts/db-dump.sh --gzip tidb
```

初回は空だが、やり直し (2回目以降の移行) では TiDB 側に中途半端なデータが入った状態で
再実行することになるので、毎回取る。

### 4. TiDB にスキーマを作る

```bash
op run --env-file=terraform/.env -- ./scripts/db-restore.sh db/init/01-schema.sql
```

**本番 MariaDB のスキーマをそのままコピーしない。** 本番のテーブルには `FULLTEXT KEY idx_bigram`
が残っている可能性があり、TiDB では作れない。`db/init/01-schema.sql` (FULLTEXT 削除済み) を正とする。

やり直すときは先に既存テーブルを消す (FK があるので子テーブルから):

```bash
cat >/tmp/blog4-drop.sql <<'SQL'
DROP TABLE IF EXISTS entry_link, entry_image, admin_session, amazon_cache, entry;
SQL
op run --env-file=terraform/.env -- ./scripts/db-restore.sh --yes /tmp/blog4-drop.sql
```

### 5. MariaDB から TiDB へデータを移す

```bash
# データのみダンプ (CREATE TABLE を含めない。セッションは移さない)
DUMP=$(op run --env-file=terraform/.env -- ./scripts/db-dump.sh --data-only --exclude-sessions mariadb | tail -1)
echo "$DUMP"

# TiDB へ流し込む
op run --env-file=terraform/.env -- ./scripts/db-restore.sh "$DUMP"
```

`admin_session` を除外しているので、カットオーバー後は管理画面に再ログインが必要になる。

### 6. カウントを照合する

```bash
op run --env-file=terraform/.env -- ./scripts/db-count.sh both
```

`diff` 列の `*` が `admin_session` だけになっていれば成功
(セッションは意図的に移していないため差が出る)。`entry` / `entry:public` / `entry_image` /
`entry_link` / `amazon_cache` が一致していることを手順1の出力と突き合わせて確認する。

### 7. カットオーバー

1. 移行中に記事を書かない (書いたら手順5をやり直す)。確実にやるなら AppRun を止めるか、
   管理画面を触らない時間帯を選ぶ
2. 1Password の `blog4-app-db` を TiDB の値に更新する:

   | field | 新しい値 |
   |---|---|
   | `host` | `terraform output tidb_hostname` の値 (`blog5.tidb-is1.db.sakurausercontent.com`) |
   | `port` | `3306` (TiDB 標準の 4000 ではない) |
   | `name` | `blog5` |
   | `user` | `blog5` (= データベース名) |
   | `password` | `blog4-tidb` の `root-password` と同じ値 |

3. AppRun に反映する:

   ```bash
   cd terraform
   op run --env-file=.env -- terraform plan    # env 以外に差分が出ていないことを確認
   op run --env-file=.env -- terraform apply
   ```

   アプリが読む DB 名の env は `DATABASE_NAME`。AppRun 側に旧名の `DATABASE_DB` しか
   入っていない状態でデプロイすると、env が無視されてデフォルトの `blog3` に繋ぎにいく。
   apply 後に AppRun の env に `DATABASE_NAME` が入っていることを確認する。

4. 動作確認: `/healthz`、トップページ、管理画面のログインと記事の保存
5. 直前にもう一度 `db-count.sh both` を回し、切り替え後に記事数が減っていないか見る

### 8. 切り戻し

TiDB 側で問題が出たら、1Password の `blog4-app-db` を旧 MariaDB の値に戻して
`terraform apply` する。DB のデータは手順2のダンプから戻せる:

```bash
op run --env-file=terraform/.env -- ./scripts/db-restore.sh --target mariadb /tmp/blog4-mariadb-<timestamp>.sql.gz
```

旧 EDB を廃止するまでは、この経路を残しておく。

### 9. 後片付け

- `/tmp` のダンプは**平文の記事データ**。作業が終わったら消す
  (`shred -u /tmp/blog4-*.sql*`)。長期保管するなら `blog4-backup` バケットへ
- 数日運用して問題なければ旧 EDB MariaDB を廃止する
- 廃止後、この runbook の MariaDB 側の手順は不要になる

## 補足

- TiDB CR の `max_connections` は 50 (`terraform output tidb_max_connections`)。
  ダンプ/リストア中にアプリが動いていても枯渇しない範囲だが、並列実行はしない
- `db-dump.sh` は `--single-transaction --no-tablespaces --skip-add-locks --skip-disable-keys`
  を付けている。EDB の制限された権限でも通り、`LOCK TABLES` / `DISABLE KEYS` が
  TiDB へのリストアで問題にならないようにするため
- ローカルの docker-compose の MariaDB を相手にスクリプトを試すときは
  `DB_DOCKER_NETWORK=blog4_default MARIADB_TLS_CNF=ssl=0` を付ける
