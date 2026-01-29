import { SaveFeedback } from './SaveFeedback.jsx';
import { VisibilityControl } from './VisibilityControl.jsx';
import { ActionButtons } from './ActionButtons.jsx';

export function Sidebar({ feedback, visibility, path, onVisibilityChange, onDelete, onRegenerateImage }) {
    return (
        <div class="edit-sidebar">
            <SaveFeedback feedback={feedback} />
            {visibility === 'public' && (
                <div class="control-panel">
                    <a href={`/entry/${path}`} target="_blank" class="btn btn-secondary" style={{ textAlign: 'center', textDecoration: 'none', display: 'block' }}>
                        View Public Page
                    </a>
                </div>
            )}
            <VisibilityControl visibility={visibility} onVisibilityChange={onVisibilityChange} />
            <ActionButtons onDelete={onDelete} onRegenerateImage={onRegenerateImage} />
        </div>
    );
}
