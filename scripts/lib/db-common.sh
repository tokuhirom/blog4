#!/usr/bin/env bash
# blog4 の DB (TiDB CR / 旧 EDB MariaDB) に接続するための共通処理。
#
# 単体では実行しない。db-count.sh / db-dump.sh / db-restore.sh から source する。
#
# 認証情報の正本は 1Password (blog4 vault) なので、呼び出し側は op run 経由で実行する:
#   op run --env-file=terraform/.env -- ./scripts/db-count.sh
#
# ローカルに mysql クライアントを入れずに済むよう、クライアントは docker で動かす。
# TiDB は MySQL 8 互換なので mysql:8, 旧 EDB は mariadb: のクライアントを使う。

TIDB_CLIENT_IMAGE="${TIDB_CLIENT_IMAGE:-mysql:8.4}"
MARIADB_CLIENT_IMAGE="${MARIADB_CLIENT_IMAGE:-mariadb:10.11.18}"

DB_COMMON_REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"

die() {
	echo "error: $*" >&2
	exit 1
}

info() {
	echo "==> $*" >&2
}

# terraform output から TiDB のホスト名を取る (TIDB_HOST 未指定時のフォールバック)。
tidb_hostname_from_terraform() {
	local host
	if ! host="$(cd "$DB_COMMON_REPO_ROOT/terraform" && terraform output -raw tidb_hostname 2>/dev/null)"; then
		die "TiDB のホスト名を解決できない。TIDB_HOST を明示するか、terraform/ で init 済みの状態で実行する"
	fi
	[ -n "$host" ] || die "terraform output tidb_hostname が空"
	echo "$host"
}

# 接続情報を DB_* に展開する。
#   db_config tidb|mariadb
db_config() {
	local target="$1"

	case "$target" in
	tidb)
		# TiDB CR (さくらのオンデマンドデータベース)。TLS 必須。
		# ポートは TiDB 標準の 4000 ではなく 3306 (2026-07-23 に疎通確認。4000 は閉じている)。
		# ユーザ名はデータベース名と同じ (EDB/オンデマンドDB の仕様)。違う場合は TIDB_USER で上書き。
		DB_NAME="${TIDB_DATABASE:-${TF_VAR_tidb_database_name:-blog4}}"
		DB_HOST="${TIDB_HOST:-$(tidb_hostname_from_terraform)}"
		DB_PORT="${TIDB_PORT:-3306}"
		DB_USER="${TIDB_USER:-$DB_NAME}"
		DB_PASSWORD="${TIDB_PASSWORD:-${TF_VAR_tidb_password:-}}"
		DB_IMAGE="$TIDB_CLIENT_IMAGE"
		DB_TLS_CNF="${TIDB_TLS_CNF:-ssl-mode=REQUIRED}"
		;;
	mariadb)
		# 旧 EDB MariaDB = 現在アプリが繋いでいる DB。値は 1Password の blog4-app-db。
		DB_HOST="${TF_VAR_database_host:-}"
		DB_PORT="${TF_VAR_database_port:-3306}"
		DB_NAME="${TF_VAR_database_name:-}"
		DB_USER="${TF_VAR_database_user:-}"
		DB_PASSWORD="${TF_VAR_database_password:-}"
		DB_IMAGE="$MARIADB_CLIENT_IMAGE"
		# ローカルの docker-compose MariaDB を相手に検証するときは MARIADB_TLS_CNF="ssl=0"。
		DB_TLS_CNF="${MARIADB_TLS_CNF:-ssl=1}"
		[ -n "$DB_HOST" ] || die "TF_VAR_database_host が空。op run --env-file=terraform/.env 経由で実行しているか確認する"
		[ -n "$DB_NAME" ] || die "TF_VAR_database_name が空"
		[ -n "$DB_USER" ] || die "TF_VAR_database_user が空"
		;;
	*)
		die "unknown target: $target (tidb|mariadb)"
		;;
	esac

	[ -n "$DB_PASSWORD" ] || die "$target のパスワードが空。op run --env-file=terraform/.env 経由で実行しているか確認する"

	DB_TARGET="$target"
	db_write_defaults_file
}

# my.cnf の値をエスケープする (バックスラッシュとダブルクォート)。
cnf_escape() {
	printf '%s' "$1" | sed -e 's/\\/\\\\/g' -e 's/"/\\"/g'
}

# パスワードを引数に置くと ps や docker inspect から見えるので、
# 600 の一時ファイルに書いてコンテナへ read-only マウントする。
# 後始末は呼び出し側が `trap db_cleanup EXIT` で行う。
db_write_defaults_file() {
	# 同一プロセスで複数の接続先を切り替えるときに前の分を残さない。
	db_cleanup

	DB_DEFAULTS_FILE="$(mktemp "${TMPDIR:-/tmp}/blog4-client.cnf.XXXXXX")"
	chmod 600 "$DB_DEFAULTS_FILE"

	cat >"$DB_DEFAULTS_FILE" <<EOF
[client]
host="$(cnf_escape "$DB_HOST")"
port=$DB_PORT
user="$(cnf_escape "$DB_USER")"
password="$(cnf_escape "$DB_PASSWORD")"
$DB_TLS_CNF
default-character-set=utf8mb4
EOF
}

# 認証情報を書いた一時ファイルを消す。呼び出し側は trap db_cleanup EXIT すること。
db_cleanup() {
	if [ -n "${DB_DEFAULTS_FILE:-}" ]; then
		rm -f "$DB_DEFAULTS_FILE"
		DB_DEFAULTS_FILE=""
	fi
}

# docker 経由で mysql / mysqldump を実行する。
#   db_client mysql [args...]
db_client() {
	local prog="$1"
	shift

	local docker_args=(--rm -i -v "$DB_DEFAULTS_FILE:/etc/blog4-client.cnf:ro")
	# ローカルの docker-compose の DB を相手にするとき用 (例: DB_DOCKER_NETWORK=blog4_default)。
	if [ -n "${DB_DOCKER_NETWORK:-}" ]; then
		docker_args+=(--network "$DB_DOCKER_NETWORK")
	fi

	docker run "${docker_args[@]}" "$DB_IMAGE" \
		"$prog" --defaults-extra-file=/etc/blog4-client.cnf "$@"
}

# SQL を投げてタブ区切り (ヘッダなし) で受け取る。
db_query_tsv() {
	db_client mysql --batch --skip-column-names --database="$DB_NAME" -e "$1"
}

# 接続先の表示用文字列。
db_describe() {
	echo "$DB_TARGET ($DB_USER@$DB_HOST:$DB_PORT/$DB_NAME)"
}
