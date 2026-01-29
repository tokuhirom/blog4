export function TitleInput({ value, onChange }) {
    return (
        <div class="title-section">
            <input
                type="text"
                name="title"
                value={value}
                placeholder="Entry Title"
                class="title-input"
                onInput={(e) => onChange(e.target.value)}
            />
        </div>
    );
}
