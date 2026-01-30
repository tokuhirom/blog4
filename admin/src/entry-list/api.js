export async function fetchEntries({ query, lastCursor }) {
    const params = new URLSearchParams();
    if (query) params.set('q', query);
    if (lastCursor) params.set('last_last_edited_at', lastCursor);

    const res = await fetch(`/admin/api/entries?${params.toString()}`);
    if (!res.ok) throw new Error('Failed to fetch entries');
    return res.json();
}

export async function createEntry(title) {
    const res = await fetch('/admin/api/entries/create', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ title }),
    });
    if (!res.ok) throw new Error('Failed to create entry');
    return res.json();
}
