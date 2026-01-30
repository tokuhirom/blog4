import { useCallback } from 'preact/hooks';

export function SearchBox({ onSearch }) {
    const handleInput = useCallback(
        (e) => {
            onSearch(e.target.value);
        },
        [onSearch],
    );

    return (
        <div class="search-box">
            <input
                type="text"
                placeholder="Search entries..."
                onInput={handleInput}
                autocomplete="off"
            />
        </div>
    );
}
