import { useRef, useEffect } from 'preact/hooks';
import { createEditor, getContent, insertAtCursor } from '../../codemirror-editor.js';
import { uploadImage } from '../api.js';

export function BodyEditor({ initialBody, onBodyChange, onFeedback }) {
    const containerRef = useRef(null);
    const editorRef = useRef(null);

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

    return (
        <div class="body-section">
            <div ref={containerRef} class="editor-container" />
        </div>
    );
}
