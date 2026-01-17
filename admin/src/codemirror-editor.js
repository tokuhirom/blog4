// CodeMirror 6 Editor for Blog4 Admin

import { indentWithTab } from '@codemirror/commands';
import { markdown, markdownLanguage } from '@codemirror/lang-markdown';
import { HighlightStyle, syntaxHighlighting } from '@codemirror/language';
import { languages } from '@codemirror/language-data';
import { EditorState } from '@codemirror/state';
import { EditorView, keymap } from '@codemirror/view';
import { tags } from '@lezer/highlight';
import { basicSetup } from 'codemirror';

// Custom Markdown highlight style
const markdownHighlightStyle = HighlightStyle.define([
    { tag: tags.heading1, fontWeight: 'bold', fontSize: '1.4em', color: '#1a1a1a' },
    { tag: tags.heading2, fontWeight: 'bold', fontSize: '1.3em', color: '#2a2a2a' },
    { tag: tags.heading3, fontWeight: 'bold', fontSize: '1.2em', color: '#3a3a3a' },
    { tag: tags.heading4, fontWeight: 'bold', fontSize: '1.1em', color: '#4a4a4a' },
    { tag: tags.heading5, fontWeight: 'bold', color: '#5a5a5a' },
    { tag: tags.heading6, fontWeight: 'bold', color: '#6a6a6a' },
    { tag: tags.strong, fontWeight: 'bold' },
    { tag: tags.emphasis, fontStyle: 'italic' },
    { tag: tags.strikethrough, textDecoration: 'line-through' },
    { tag: tags.link, color: '#1976d2', textDecoration: 'underline' },
    { tag: tags.url, color: '#1976d2' },
    { tag: tags.monospace, fontFamily: 'Monaco, Menlo, monospace', backgroundColor: '#f5f5f5' },
    { tag: tags.quote, color: '#666', fontStyle: 'italic' },
    { tag: tags.list, color: '#e65100' },
    { tag: tags.processingInstruction, color: '#9c27b0' }, // code fence markers
]);

/**
 * Initialize CodeMirror editor
 * @param {HTMLElement} container - The container element for the editor
 * @param {string} initialValue - Initial content
 * @param {Function} onUpdate - Callback when content changes (debounced)
 * @returns {EditorView} The editor instance
 */
function createEditor(container, initialValue, onUpdate) {
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
                markdown({
                    base: markdownLanguage,
                    codeLanguages: languages,
                }),
                syntaxHighlighting(markdownHighlightStyle),
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
function getContent(editor) {
    return editor.state.doc.toString();
}

/**
 * Insert text at cursor position
 * @param {EditorView} editor
 * @param {string} text
 */
function insertAtCursor(editor, text) {
    const cursor = editor.state.selection.main.head;
    editor.dispatch({
        changes: { from: cursor, insert: text },
        selection: { anchor: cursor + text.length },
    });
    editor.focus();
}

// Export functions
export { createEditor, getContent, insertAtCursor };
