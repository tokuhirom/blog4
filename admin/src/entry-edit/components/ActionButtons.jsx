export function ActionButtons({ onDelete, onRegenerateImage }) {
    const handleDelete = () => {
        if (!confirm('Are you sure you want to delete this entry? This action cannot be undone.')) return;
        onDelete();
    };

    return (
        <div class="control-panel">
            <h3>Actions</h3>
            <button class="btn btn-danger" onClick={handleDelete}>
                Delete Entry
            </button>
            <button class="btn btn-secondary" onClick={onRegenerateImage}>
                Regenerate Image
            </button>
        </div>
    );
}
