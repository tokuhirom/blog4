import { EntryCard } from './EntryCard.jsx';

export function EntryGrid({ entries }) {
    return (
        <div class="entry-grid">
            {entries.map((entry) => (
                <EntryCard key={entry.path} entry={entry} />
            ))}
        </div>
    );
}
