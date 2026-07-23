#!/usr/bin/env bash
# blog4 の DB を mysqldump して /tmp に置く。
#
# 使い方:
#   op run --env-file=terraform/.env -- ./scripts/db-dump.sh both        # 両方バックアップ
#   op run --env-file=terraform/.env -- ./scripts/db-dump.sh mariadb     # 旧 EDB だけ
#   op run --env-file=terraform/.env -- ./scripts/db-dump.sh --data-only --exclude-sessions mariadb
#
# 移行用のデータ投入ファイルを作るときは --data-only --exclude-sessions を使う。
# スキーマは本番 MariaDB のもの (FULLTEXT インデックスが残っている可能性がある) ではなく
# db/init/01-schema.sql を正とするため、移行では CREATE TABLE をダンプに含めない。
set -euo pipefail

cd "$(dirname "$0")/.."
# shellcheck source=scripts/lib/db-common.sh
. scripts/lib/db-common.sh

usage() {
	cat >&2 <<'EOF'
usage: db-dump.sh [options] [tidb|mariadb|both]

options:
  --data-only          CREATE TABLE を含めない (移行のデータ投入用)
  --schema-only        データを含めない
  --exclude-sessions   admin_session テーブルを除外する (移行では引き継がない)
  --out-dir DIR        出力先ディレクトリ (default: /tmp)
  --gzip               gzip で圧縮する
  -h, --help           このヘルプ
EOF
	exit 1
}

target=""
data_only=0
schema_only=0
exclude_sessions=0
out_dir="/tmp"
use_gzip=0

while [ $# -gt 0 ]; do
	case "$1" in
	--data-only) data_only=1 ;;
	--schema-only) schema_only=1 ;;
	--exclude-sessions) exclude_sessions=1 ;;
	--out-dir)
		out_dir="${2:-}"
		[ -n "$out_dir" ] || usage
		shift
		;;
	--gzip) use_gzip=1 ;;
	-h | --help) usage ;;
	tidb | mariadb | both) target="$1" ;;
	*) usage ;;
	esac
	shift
done

target="${target:-both}"
if [ "$data_only" -eq 1 ] && [ "$schema_only" -eq 1 ]; then
	die "--data-only と --schema-only は同時に指定できない"
fi
[ -d "$out_dir" ] || die "出力先ディレクトリがない: $out_dir"

timestamp="$(date +%Y%m%d-%H%M%S)"

trap db_cleanup EXIT

dump_one() {
	local t="$1"

	db_config "$t"

	local args=(
		--single-transaction # LOCK TABLES 権限なしでも一貫性のあるダンプを取る
		--no-tablespaces     # PROCESS 権限がなくても失敗しないようにする
		--skip-triggers      # トリガーは使っていない
		--skip-add-locks     # LOCK TABLES 文は TiDB へのリストアで邪魔になる
		--skip-disable-keys  # ALTER TABLE ... DISABLE KEYS も同様
		--complete-insert    # 列名付き INSERT にしてカラム順の差異に強くする
		--hex-blob
	)

	# mysql 8 系クライアント固有のオプション (mariadb クライアントには無い)。
	if [ "$t" = "tidb" ]; then
		args+=(--set-gtid-purged=OFF --column-statistics=0)
	fi

	local suffix=""
	if [ "$data_only" -eq 1 ]; then
		args+=(--no-create-info)
		suffix="-data"
	elif [ "$schema_only" -eq 1 ]; then
		args+=(--no-data)
		suffix="-schema"
	fi

	if [ "$exclude_sessions" -eq 1 ]; then
		args+=("--ignore-table=${DB_NAME}.admin_session")
	fi

	local out="$out_dir/blog4-${t}-${timestamp}${suffix}.sql"
	if [ "$use_gzip" -eq 1 ]; then
		out="${out}.gz"
	fi

	info "dumping: $(db_describe) -> $out"

	# ダンプ途中で失敗したファイルを残さない (中途半端なファイルをリストアしないため)。
	local tmp_out="${out}.partial"
	if [ "$use_gzip" -eq 1 ]; then
		db_client mysqldump "${args[@]}" "$DB_NAME" | gzip >"$tmp_out"
	else
		db_client mysqldump "${args[@]}" "$DB_NAME" >"$tmp_out"
	fi
	mv "$tmp_out" "$out"
	chmod 600 "$out"

	info "done: $out ($(du -h "$out" | cut -f1))"
	echo "$out"
}

case "$target" in
both)
	dump_one mariadb
	dump_one tidb
	;;
*)
	dump_one "$target"
	;;
esac
