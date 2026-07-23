#!/usr/bin/env bash
# SQL ファイルを DB に流し込む。移行のリストア (MariaDB のダンプ → TiDB) 用。
#
# 使い方:
#   # スキーマを作る (db/init/01-schema.sql が TiDB 互換スキーマの正)
#   op run --env-file=terraform/.env -- ./scripts/db-restore.sh db/init/01-schema.sql
#
#   # MariaDB から取ったデータダンプを流し込む
#   op run --env-file=terraform/.env -- ./scripts/db-restore.sh /tmp/blog4-mariadb-20260723-101500-data.sql
#
# デフォルトの流し込み先は TiDB。MariaDB へ戻す (切り戻し) 場合のみ --target mariadb。
set -euo pipefail

cd "$(dirname "$0")/.."
# shellcheck source=scripts/lib/db-common.sh
. scripts/lib/db-common.sh

usage() {
	cat >&2 <<'EOF'
usage: db-restore.sh [options] <file.sql|file.sql.gz>

options:
  --target tidb|mariadb   流し込み先 (default: tidb)
  --yes                   確認プロンプトを出さない
  -h, --help              このヘルプ
EOF
	exit 1
}

target="tidb"
assume_yes=0
file=""

while [ $# -gt 0 ]; do
	case "$1" in
	--target)
		target="${2:-}"
		[ -n "$target" ] || usage
		shift
		;;
	--yes) assume_yes=1 ;;
	-h | --help) usage ;;
	-*) usage ;;
	*) file="$1" ;;
	esac
	shift
done

[ -n "$file" ] || usage
[ -f "$file" ] || die "ファイルがない: $file"

trap db_cleanup EXIT

db_config "$target"

# 流し込む前に現状を見せる。空でない DB に上書きするときに気づけるようにする。
existing="$(db_query_tsv "SELECT COUNT(*) FROM information_schema.tables
                          WHERE table_schema = DATABASE() AND table_type = 'BASE TABLE'")"

echo "restore 先 : $(db_describe)" >&2
echo "既存テーブル: ${existing} 個" >&2
echo "流し込む file: $file" >&2

if [ "$assume_yes" -ne 1 ]; then
	printf 'この内容で実行する? [y/N] ' >&2
	read -r answer
	case "$answer" in
	y | Y | yes | YES) ;;
	*) die "中止した" ;;
	esac
fi

info "restoring..."
case "$file" in
*.gz) gzip -dc "$file" | db_client mysql --database="$DB_NAME" ;;
*) db_client mysql --database="$DB_NAME" <"$file" ;;
esac

info "done. ./scripts/db-count.sh $target で件数を確認する"
