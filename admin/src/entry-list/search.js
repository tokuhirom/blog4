// 軽いスコアリング検索: スペース区切りで複数語 AND、title ヒットを加点、大小無視。
// query が空なら元の並び (last_edited_at DESC) のまま返す。
export function searchEntries(entries, query) {
    const q = query.trim();
    if (!q) {
        return entries;
    }

    const terms = q.toLowerCase().split(/\s+/).filter(Boolean);
    const scored = [];
    for (const entry of entries) {
        const title = entry.title.toLowerCase();
        const body = (entry.body || '').toLowerCase();
        let score = 0;
        let matchedAll = true;
        for (const term of terms) {
            const inTitle = title.includes(term);
            const inBody = body.includes(term);
            if (!inTitle && !inBody) {
                matchedAll = false;
                break;
            }
            if (inTitle) {
                score += 10;
            }
            if (inBody) {
                score += 1;
            }
        }
        if (matchedAll) {
            scored.push({ entry, score });
        }
    }

    // 高スコア順。Array.sort は stable なので同点は元の順 (last_edited_at DESC) を保つ。
    scored.sort((a, b) => b.score - a.score);
    return scored.map((x) => x.entry);
}
