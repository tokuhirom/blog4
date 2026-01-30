import { useCallback, useState } from 'preact/hooks';
import { createEntry } from '../api.js';

export function NewEntryButton() {
    const [creating, setCreating] = useState(false);

    const handleClick = useCallback(
        async (e) => {
            e.preventDefault();
            if (creating) return;
            setCreating(true);
            try {
                const now = new Date();
                const title = now.toISOString().slice(0, 19).replace(/[-:]/g, '').replace('T', '');
                const data = await createEntry(title);
                if (data.ok && data.redirect) {
                    window.location.href = data.redirect;
                }
            } catch {
                setCreating(false);
            }
        },
        [creating],
    );

    return (
        <button type="button" class="btn-new-entry" onClick={handleClick} disabled={creating}>
            <svg
                width="20"
                height="20"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                aria-hidden="true"
            >
                <line x1="12" y1="5" x2="12" y2="19" />
                <line x1="5" y1="12" x2="19" y2="12" />
            </svg>
            {creating ? 'Creating...' : 'New Entry'}
        </button>
    );
}
