import { useCallback, useReducer, useRef } from 'preact/hooks';
import { fetchEntries } from './api.js';
import { EntryGrid } from './components/EntryGrid.jsx';
import { NewEntryButton } from './components/NewEntryButton.jsx';
import { SearchBox } from './components/SearchBox.jsx';

function reducer(state, action) {
    switch (action.type) {
        case 'SET_QUERY':
            return { ...state, query: action.query };
        case 'SEARCH_START':
            return { ...state, loading: true };
        case 'SEARCH_RESULT':
            return {
                ...state,
                entries: action.entries,
                hasMore: action.hasMore,
                lastCursor: action.lastCursor,
                loading: false,
            };
        case 'LOAD_MORE_RESULT':
            return {
                ...state,
                entries: [...state.entries, ...action.entries],
                hasMore: action.hasMore,
                lastCursor: action.lastCursor,
                loading: false,
            };
        case 'SEARCH_ERROR':
            return { ...state, loading: false };
        default:
            return state;
    }
}

export function App({ initData }) {
    const [state, dispatch] = useReducer(reducer, {
        entries: initData.entries || [],
        hasMore: initData.has_more,
        lastCursor: initData.last_cursor,
        query: '',
        loading: false,
    });

    const debounceRef = useRef(null);

    const handleSearch = useCallback((query) => {
        dispatch({ type: 'SET_QUERY', query });

        if (debounceRef.current) clearTimeout(debounceRef.current);

        debounceRef.current = setTimeout(async () => {
            dispatch({ type: 'SEARCH_START' });
            try {
                const data = await fetchEntries({ query });
                dispatch({
                    type: 'SEARCH_RESULT',
                    entries: data.entries || [],
                    hasMore: data.has_more,
                    lastCursor: data.last_cursor,
                });
            } catch {
                dispatch({ type: 'SEARCH_ERROR' });
            }
        }, 500);
    }, []);

    const handleLoadMore = useCallback(async () => {
        if (state.loading || !state.hasMore) return;
        dispatch({ type: 'SEARCH_START' });
        try {
            const data = await fetchEntries({
                query: state.query,
                lastCursor: state.lastCursor,
            });
            dispatch({
                type: 'LOAD_MORE_RESULT',
                entries: data.entries || [],
                hasMore: data.has_more,
                lastCursor: data.last_cursor,
            });
        } catch {
            dispatch({ type: 'SEARCH_ERROR' });
        }
    }, [state.loading, state.hasMore, state.query, state.lastCursor]);

    return (
        <>
            <div class="search-and-actions">
                <SearchBox onSearch={handleSearch} />
                <NewEntryButton />
            </div>
            <EntryGrid entries={state.entries} />
            {state.hasMore && (
                <div class="load-more">
                    <button type="button" onClick={handleLoadMore} disabled={state.loading}>
                        {state.loading ? 'Loading...' : 'Load More'}
                    </button>
                </div>
            )}
        </>
    );
}
