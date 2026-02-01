import { useReducer, useCallback, useRef } from 'preact/hooks';
import { TitleInput } from './components/TitleInput.jsx';
import { BodyEditor } from './components/BodyEditor.jsx';
import { Sidebar } from './components/Sidebar.jsx';
import { useAutoSave } from './hooks/useAutoSave.js';
import * as api from './api.js';

function reducer(state, action) {
    switch (action.type) {
        case 'SET_TITLE':
            return { ...state, title: action.value };
        case 'SET_BODY':
            return { ...state, body: action.value };
        case 'SET_VISIBILITY':
            return { ...state, visibility: action.value };
        case 'SET_UPDATED_AT':
            return { ...state, updatedAt: action.value };
        case 'SET_FEEDBACK':
            return { ...state, feedback: action.value };
        case 'CLEAR_FEEDBACK':
            return { ...state, feedback: null };
        default:
            return state;
    }
}

export function App({ initData }) {
    const [state, dispatch] = useReducer(reducer, {
        title: initData.title,
        body: initData.body,
        visibility: initData.visibility,
        updatedAt: initData.updated_at,
        feedback: null,
    });

    const updatedAtRef = useRef(state.updatedAt);
    const feedbackTimerRef = useRef(null);

    const showFeedback = useCallback((feedback) => {
        dispatch({ type: 'SET_FEEDBACK', value: feedback });
        if (feedbackTimerRef.current) clearTimeout(feedbackTimerRef.current);
        const delay = feedback.type === 'error' ? 5000 : 3000;
        feedbackTimerRef.current = setTimeout(() => {
            dispatch({ type: 'CLEAR_FEEDBACK' });
        }, delay);
    }, []);

    const handleApiResponse = useCallback((data) => {
        if (data.error) {
            showFeedback({ type: 'error', message: data.error });
            return false;
        }
        if (data.updated_at) {
            updatedAtRef.current = data.updated_at;
            dispatch({ type: 'SET_UPDATED_AT', value: data.updated_at });
        }
        if (data.message) {
            showFeedback({ type: 'success', message: data.message });
        }
        return true;
    }, [showFeedback]);

    const saveTitle = useCallback(async (title) => {
        try {
            const data = await api.updateTitle(initData.path, title, updatedAtRef.current);
            handleApiResponse(data);
        } catch (err) {
            showFeedback({ type: 'error', message: `Failed to save title: ${err.message}` });
        }
    }, [initData.path, handleApiResponse, showFeedback]);

    const saveBody = useCallback(async (body) => {
        try {
            const data = await api.updateBody(initData.path, body, updatedAtRef.current);
            handleApiResponse(data);
        } catch (err) {
            showFeedback({ type: 'error', message: `Failed to save body: ${err.message}` });
        }
    }, [initData.path, handleApiResponse, showFeedback]);

    const debouncedSaveTitle = useAutoSave(saveTitle, 500);
    const debouncedSaveBody = useAutoSave(saveBody, 800);

    const handleTitleChange = useCallback((title) => {
        dispatch({ type: 'SET_TITLE', value: title });
        debouncedSaveTitle(title);
    }, [debouncedSaveTitle]);

    const handleBodyChange = useCallback((body) => {
        dispatch({ type: 'SET_BODY', value: body });
        debouncedSaveBody(body);
    }, [debouncedSaveBody]);

    const handleVisibilityChange = useCallback(async (visibility) => {
        try {
            const data = await api.updateVisibility(initData.path, visibility);
            if (handleApiResponse(data)) {
                dispatch({ type: 'SET_VISIBILITY', value: visibility });
            }
        } catch (err) {
            showFeedback({ type: 'error', message: `Failed to update visibility: ${err.message}` });
        }
    }, [initData.path, handleApiResponse, showFeedback]);

    const handleDelete = useCallback(async () => {
        try {
            const data = await api.deleteEntry(initData.path);
            if (data.redirect) {
                window.location.href = data.redirect;
            } else if (data.error) {
                showFeedback({ type: 'error', message: data.error });
            }
        } catch (err) {
            showFeedback({ type: 'error', message: `Failed to delete entry: ${err.message}` });
        }
    }, [initData.path, showFeedback]);

    const handleRegenerateImage = useCallback(async () => {
        try {
            const data = await api.regenerateImage(initData.path);
            handleApiResponse(data);
        } catch (err) {
            showFeedback({ type: 'error', message: `Failed to regenerate image: ${err.message}` });
        }
    }, [initData.path, handleApiResponse, showFeedback]);

    return (
        <div class="edit-container">
            <div class="edit-main">
                <TitleInput value={state.title} onChange={handleTitleChange} />
                <BodyEditor
                    initialBody={initData.body}
                    currentBody={state.body}
                    onBodyChange={handleBodyChange}
                    onFeedback={showFeedback}
                />
            </div>
            <Sidebar
                feedback={state.feedback}
                visibility={state.visibility}
                path={initData.path}
                onVisibilityChange={handleVisibilityChange}
                onDelete={handleDelete}
                onRegenerateImage={handleRegenerateImage}
            />
        </div>
    );
}
