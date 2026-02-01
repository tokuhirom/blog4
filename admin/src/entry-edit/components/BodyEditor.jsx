import { useRef, useEffect, useState, useCallback } from 'preact/hooks';
import { createEditor, getContent, insertAtCursor } from '../../codemirror-editor.js';
import { uploadImage, previewMarkdown } from '../api.js';

export function BodyEditor({ initialBody, currentBody, onBodyChange, onFeedback }) {
    const containerRef = useRef(null);
    const editorRef = useRef(null);
    const [activeTab, setActiveTab] = useState('edit');
    const [previewHtml, setPreviewHtml] = useState('');
    const [previewLoading, setPreviewLoading] = useState(false);

    useEffect(() => {
        if (!containerRef.current || editorRef.current) return;

        const editor = createEditor(containerRef.current, initialBody, (content) => {
            onBodyChange(content);
        });
        editorRef.current = editor;

        containerRef.current.addEventListener('paste', async (event) => {
            const items = event.clipboardData?.items || [];
            const imageFiles = [];
            for (const item of items) {
                if (item.kind === 'file' && item.type.startsWith('image/')) {
                    const file = item.getAsFile();
                    if (file) imageFiles.push(file);
                }
            }
            if (imageFiles.length === 0) return;
            event.preventDefault();

            for (const file of imageFiles) {
                try {
                    const data = await uploadImage(file);
                    if (data.error) {
                        onFeedback({ type: 'error', message: data.error });
                        return;
                    }
                    if (!data.url) {
                        onFeedback({ type: 'error', message: 'No URL in response' });
                        return;
                    }
                    const markdownImage = `![image](${data.url})`;
                    insertAtCursor(editor, markdownImage);
                    onBodyChange(getContent(editor));
                    onFeedback({ type: 'success', message: 'Image uploaded successfully' });
                } catch (err) {
                    onFeedback({ type: 'error', message: `Failed to upload image: ${err.message}` });
                }
            }
        });
    }, []);

    const handlePreviewClick = useCallback(async () => {
        setActiveTab('preview');
        setPreviewLoading(true);
        try {
            const body = editorRef.current ? getContent(editorRef.current) : (currentBody || '');
            const data = await previewMarkdown(body);
            if (data.error) {
                onFeedback({ type: 'error', message: data.error });
                setPreviewHtml('<p>Failed to load preview.</p>');
            } else {
                setPreviewHtml(data.html);
            }
        } catch (err) {
            onFeedback({ type: 'error', message: `Failed to load preview: ${err.message}` });
            setPreviewHtml('<p>Failed to load preview.</p>');
        } finally {
            setPreviewLoading(false);
        }
    }, [currentBody, onFeedback]);

    const handleEditClick = useCallback(() => {
        setActiveTab('edit');
    }, []);

    return (
        <div class="body-section">
            <div class="editor-tabs">
                <button
                    type="button"
                    class={`editor-tab ${activeTab === 'edit' ? 'editor-tab-active' : ''}`}
                    onClick={handleEditClick}
                >
                    Edit
                </button>
                <button
                    type="button"
                    class={`editor-tab ${activeTab === 'preview' ? 'editor-tab-active' : ''}`}
                    onClick={handlePreviewClick}
                >
                    Preview
                </button>
            </div>
            <div ref={containerRef} class="editor-container" style={{ display: activeTab === 'edit' ? '' : 'none' }} />
            {activeTab === 'preview' && (
                <div class="preview-container">
                    {previewLoading ? (
                        <div class="preview-loading">Loading preview...</div>
                    ) : (
                        <div class="preview-content" dangerouslySetInnerHTML={{ __html: previewHtml }} />
                    )}
                </div>
            )}
        </div>
    );
}
