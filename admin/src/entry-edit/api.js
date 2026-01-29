export async function updateTitle(path, title, updatedAt) {
    const res = await fetch(`/admin/api/entries/title?path=${encodeURIComponent(path)}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ title, updated_at: updatedAt }),
    });
    return res.json();
}

export async function updateBody(path, body, updatedAt) {
    const res = await fetch(`/admin/api/entries/body?path=${encodeURIComponent(path)}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ body, updated_at: updatedAt }),
    });
    return res.json();
}

export async function updateVisibility(path, visibility) {
    const res = await fetch(`/admin/api/entries/visibility?path=${encodeURIComponent(path)}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ visibility }),
    });
    return res.json();
}

export async function deleteEntry(path) {
    const res = await fetch(`/admin/api/entries/delete?path=${encodeURIComponent(path)}`, {
        method: 'DELETE',
    });
    return res.json();
}

export async function regenerateImage(path) {
    const res = await fetch(`/admin/api/entries/image/regenerate?path=${encodeURIComponent(path)}`, {
        method: 'POST',
    });
    return res.json();
}

export async function uploadImage(file) {
    const formData = new FormData();
    formData.append('file', file);
    const res = await fetch('/admin/entries/upload', {
        method: 'POST',
        body: formData,
    });
    return res.json();
}
