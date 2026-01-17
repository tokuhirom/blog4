// CodeMirror 6 Editor for Blog4 Admin
// Uses Import Maps defined in layout.html for module resolution

import { indentWithTab } from '@codemirror/commands';
import { markdown } from '@codemirror/lang-markdown';
import { EditorState } from '@codemirror/state';
import { keymap } from '@codemirror/view';
import { basicSetup, EditorView } from 'codemirror';

/**
 * Initialize CodeMirror editor
 * @param {HTMLElement} container - The container element for the editor
 * @param {string} initialValue - Initial content
 * @param {Function} onUpdate - Callback when content changes (debounced)
 * @returns {EditorView} The editor instance
 */
export function createEditor(container, initialValue, onUpdate) {
    let debounceTimer = null;
    const debounceDelay = 800; // Match the original textarea delay

    const updateListener = EditorView.updateListener.of((update) => {
        if (update.docChanged) {
            clearTimeout(debounceTimer);
            debounceTimer = setTimeout(() => {
                onUpdate(update.state.doc.toString());
            }, debounceDelay);
        }
    });

    const editor = new EditorView({
        state: EditorState.create({
            doc: initialValue,
            extensions: [
                basicSetup,
                markdown(),
                keymap.of([indentWithTab]),
                updateListener,
                EditorView.lineWrapping,
                EditorView.theme({
                    '&': {
                        height: '100%',
                        minHeight: '600px',
                    },
                    '.cm-scroller': {
                        overflow: 'auto',
                        fontFamily: 'Monaco, Menlo, "Ubuntu Mono", monospace',
                        fontSize: '14px',
                    },
                    '.cm-content': {
                        minHeight: '600px',
                    },
                }),
            ],
        }),
        parent: container,
    });

    return editor;
}

/**
 * Get the current content from the editor
 * @param {EditorView} editor
 * @returns {string}
 */
export function getContent(editor) {
    return editor.state.doc.toString();
}

/**
 * Set content in the editor
 * @param {EditorView} editor
 * @param {string} content
 */
export function setContent(editor, content) {
    editor.dispatch({
        changes: {
            from: 0,
            to: editor.state.doc.length,
            insert: content,
        },
    });
}

/**
 * Insert text at cursor position
 * @param {EditorView} editor
 * @param {string} text
 */
export function insertAtCursor(editor, text) {
    const cursor = editor.state.selection.main.head;
    editor.dispatch({
        changes: { from: cursor, insert: text },
        selection: { anchor: cursor + text.length },
    });
    editor.focus();
}
