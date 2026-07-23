#!/usr/bin/env bash
# blog4 の DB に入っている件数 (記事数など) を表示する。
#
# 使い方:
#   op run --env-file=terraform/.env -- ./scripts/db-count.sh            # TiDB (デフォルト)
#   op run --env-file=terraform/.env -- ./scripts/db-count.sh mariadb    # 旧 EDB MariaDB
#   op run --env-file=terraform/.env -- ./scripts/db-count.sh both       # 両方を並べて比較
#
# 移行後の件数一致チェックには both を使う。差があるテーブルには * が付く。
set -euo pipefail

cd "$(dirname "$0")/.."
# shellcheck source=scripts/lib/db-common.sh
. scripts/lib/db-common.sh

target="${1:-tidb}"
case "$target" in
tidb | mariadb | both) ;;
*) die "usage: $0 [tidb|mariadb|both]" ;;
esac

# 指定 DB の "名前<TAB>件数" を outfile に書く。
collect_counts() {
	local t="$1" outfile="$2"

	db_config "$t"
	info "counting: $(db_describe)"

	local tables
	tables="$(db_query_tsv "SELECT table_name FROM information_schema.tables
	                        WHERE table_schema = DATABASE() AND table_type = 'BASE TABLE'
	                        ORDER BY table_name")"

	if [ -z "$tables" ]; then
		info "テーブルが1つもない (移行前の TiDB ならこれが正常)"
		: >"$outfile"
		return
	fi

	local sql="" t_name
	while IFS= read -r t_name; do
		[ -n "$t_name" ] || continue
		if [ -n "$sql" ]; then
			sql+=" UNION ALL "
		fi
		sql+="SELECT '$t_name' AS n, COUNT(*) AS c FROM \`$t_name\`"
		# entry は visibility の内訳も出す (公開記事が何件かは移行検証で見たい)。
		if [ "$t_name" = "entry" ]; then
			sql+=" UNION ALL SELECT CONCAT('entry:', visibility), COUNT(*) FROM \`entry\` GROUP BY visibility"
		fi
	done <<<"$tables"

	db_query_tsv "$sql" >"$outfile"
}

# 表示順: テーブル名のアルファベット順。entry の直後に entry:public / entry:private。
sort_keys() {
	LC_ALL=C sort -u -t: -k1,1 -k2,2
}

# counts ファイル群に現れる行名を表示順に列挙する。
list_keys() {
	cut -f1 "$@" | sort_keys
}

# counts ファイルから行名に対応する件数を取る (無ければ空)。
lookup_count() {
	awk -F'\t' -v n="$2" '$1 == n { print $2 }' "$1"
}

tmpdir="$(mktemp -d "${TMPDIR:-/tmp}/blog4-count.XXXXXX")"
trap 'rm -rf "$tmpdir"; db_cleanup' EXIT

case "$target" in
tidb | mariadb)
	collect_counts "$target" "$tmpdir/counts"
	echo
	printf '%-20s %12s\n' "table" "$target"
	printf '%-20s %12s\n' "--------------------" "------------"
	list_keys "$tmpdir/counts" | while IFS= read -r name; do
		[ -n "$name" ] || continue
		printf '%-20s %12s\n' "$name" "$(lookup_count "$tmpdir/counts" "$name")"
	done
	;;
both)
	collect_counts mariadb "$tmpdir/mariadb"
	collect_counts tidb "$tmpdir/tidb"
	echo
	printf '%-20s %12s %12s   %s\n' "table" "mariadb" "tidb" "diff"
	printf '%-20s %12s %12s   %s\n' "--------------------" "------------" "------------" "----"
	# 両方に現れる行名を突き合わせる。片方にしかないものは - を表示し、差があれば * を付ける。
	list_keys "$tmpdir/mariadb" "$tmpdir/tidb" | while IFS= read -r name; do
		[ -n "$name" ] || continue
		m_count="$(lookup_count "$tmpdir/mariadb" "$name")"
		t_count="$(lookup_count "$tmpdir/tidb" "$name")"
		mark=""
		if [ "$m_count" != "$t_count" ]; then
			mark="*"
		fi
		printf '%-20s %12s %12s   %s\n' "$name" "${m_count:--}" "${t_count:--}" "$mark"
	done
	;;
esac
