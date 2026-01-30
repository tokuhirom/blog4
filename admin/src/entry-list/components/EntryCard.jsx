export function EntryCard({ entry }) {
    const className = `entry-card${entry.visibility === 'private' ? ' private' : ''}`;

    return (
        <a href={`/admin/entries/edit?path=${entry.path}`} class={className}>
            <div class="entry-card-content">
                <h3 class="entry-title">{entry.title}</h3>
                {entry.image_url && <img src={entry.image_url} class="entry-image" alt="" />}
                <p class="entry-body">{entry.body_preview}</p>
            </div>
        </a>
    );
}
