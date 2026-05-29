import { useEffect, useMemo, useState } from 'preact/hooks';
import { fetchAllEntries } from './api.js';
import { EntryGrid } from './components/EntryGrid.jsx';
import { NewEntryButton } from './components/NewEntryButton.jsx';
import { SearchBox } from './components/SearchBox.jsx';
import { searchEntries } from './search.js';

export function App() {
    const [allEntries, setAllEntries] = useState([]);
    const [query, setQuery] = useState('');
    const [status, setStatus] = useState('loading'); // 'loading' | 'ready' | 'error'

    useEffect(() => {
        fetchAllEntries()
            .then((entries) => {
                setAllEntries(entries);
                setStatus('ready');
            })
            .catch(() => {
                setStatus('error');
            });
    }, []);

    const visibleEntries = useMemo(() => searchEntries(allEntries, query), [allEntries, query]);

    return (
        <>
            <div class="search-and-actions">
                <SearchBox onSearch={setQuery} />
                <NewEntryButton />
            </div>
            {status === 'loading' && <p class="entry-list-status">Loading...</p>}
            {status === 'error' && <p class="entry-list-status">Failed to load entries.</p>}
            {status === 'ready' && <EntryGrid entries={visibleEntries} />}
        </>
    );
}
