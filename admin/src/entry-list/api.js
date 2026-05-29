export async function fetchAllEntries() {
    const res = await fetch('/admin/api/entries');
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
