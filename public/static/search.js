(() => {
    const resultsEl = document.getElementById('search-results');
    const inputEl = document.querySelector('.search-input');
    const formEl = document.querySelector('.search-form');

    let entries = [];

    // サーバの summarizeEntry 相当の軽量版: URL と [[...]] を整理して先頭 length 文字。
    function summarize(body, length) {
        const s = body.replace(/https?:\/\/\S+/g, '').replace(/\[\[(.*?)\]\]/g, '$1');
        return [...s].slice(0, length).join('');
    }

    // 軽いスコアリング: 全語必須 (AND)、title ヒットを加点。
    function scoreEntry(entry, terms) {
        const title = entry.title.toLowerCase();
        const body = entry.body.toLowerCase();
        let score = 0;
        for (const term of terms) {
            const inTitle = title.includes(term);
            const inBody = body.includes(term);
            if (!inTitle && !inBody) {
                return -1;
            }
            if (inTitle) {
                score += 10;
            }
            if (inBody) {
                score += 1;
            }
        }
        return score;
    }

    function renderInfo(message) {
        const div = document.createElement('div');
        div.className = 'search-result-info';
        div.textContent = message;
        resultsEl.replaceChildren(div);
    }

    function renderResults(list) {
        if (list.length === 0) {
            renderInfo('No results found. Try different keywords.');
            return;
        }

        const ul = document.createElement('ul');
        ul.className = 'card-container';
        for (const entry of list) {
            const li = document.createElement('li');
            li.className = 'card';

            const a = document.createElement('a');
            a.href = `/entry/${entry.path}`;
            a.className = 'card-link';

            const head = document.createElement('div');
            head.className = 'entry-head';
            const title = document.createElement('span');
            title.className = 'entry-title';
            title.textContent = entry.title;
            head.appendChild(title);
            a.appendChild(head);

            if (entry.image_url) {
                const wrap = document.createElement('div');
                const img = document.createElement('img');
                img.src = entry.image_url;
                img.className = 'entry-image';
                img.alt = entry.title;
                wrap.appendChild(img);
                a.appendChild(wrap);
            } else {
                const preview = document.createElement('div');
                preview.className = 'entry-text-preview';
                preview.textContent = summarize(entry.body, 100);
                a.appendChild(preview);
            }

            const date = document.createElement('span');
            date.className = 'published-date';
            date.textContent = entry.published_at;
            a.appendChild(date);

            li.appendChild(a);
            ul.appendChild(li);
        }
        resultsEl.replaceChildren(ul);
    }

    function doSearch() {
        const q = inputEl.value.trim();
        if (!q) {
            renderInfo('Enter keywords to search entries.');
            return;
        }

        const terms = q.toLowerCase().split(/\s+/).filter(Boolean);
        const scored = [];
        for (const entry of entries) {
            const score = scoreEntry(entry, terms);
            if (score >= 0) {
                scored.push({ entry, score });
            }
        }
        // 高スコア順。Array.sort は stable なので同点は元の published_at DESC を保つ。
        scored.sort((a, b) => b.score - a.score);
        renderResults(scored.map((x) => x.entry));
    }

    let timer = null;
    function scheduleSearch() {
        clearTimeout(timer);
        timer = setTimeout(doSearch, 200);
    }

    async function init() {
        renderInfo('Loading...');
        try {
            const res = await fetch('/search-index.json');
            if (!res.ok) {
                throw new Error(`HTTP ${res.status}`);
            }
            entries = await res.json();
        } catch (err) {
            renderInfo('Failed to load search index.');
            console.error(err);
            return;
        }

        inputEl.addEventListener('input', scheduleSearch);
        formEl.addEventListener('submit', (e) => {
            e.preventDefault();
            doSearch();
        });

        // ?q= に初期値があれば即検索 (input の value はサーバが埋めている)。
        doSearch();
    }

    init();
})();
