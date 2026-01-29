export function VisibilityControl({ visibility, onVisibilityChange }) {
    const handleChange = (newVisibility) => {
        if (newVisibility === visibility) return;
        const msg = `Change visibility to ${newVisibility}?`;
        if (!confirm(msg)) return;
        onVisibilityChange(newVisibility);
    };

    return (
        <div class="control-panel">
            <h3>Visibility</h3>
            <div class="visibility-control">
                <label class="radio-label">
                    <input
                        type="radio"
                        name="visibility"
                        value="private"
                        checked={visibility === 'private'}
                        onChange={() => handleChange('private')}
                    />
                    <span>Private</span>
                </label>
                <label class="radio-label">
                    <input
                        type="radio"
                        name="visibility"
                        value="public"
                        checked={visibility === 'public'}
                        onChange={() => handleChange('public')}
                    />
                    <span>Public</span>
                </label>
            </div>
        </div>
    );
}
