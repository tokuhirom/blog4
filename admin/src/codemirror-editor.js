// CodeMirror 6 Editor for Blog4 Admin

import {
    autocompletion,
    closeBrackets,
    closeBracketsKeymap,
    completionKeymap,
} from '@codemirror/autocomplete';
import { defaultKeymap, history, historyKeymap, indentWithTab } from '@codemirror/commands';
import { markdown, markdownLanguage } from '@codemirror/lang-markdown';
import {
    bracketMatching,
    defaultHighlightStyle,
    foldGutter,
    foldKeymap,
    indentOnInput,
    syntaxHighlighting,
} from '@codemirror/language';
import { languages } from '@codemirror/language-data';
import { highlightSelectionMatches, searchKeymap } from '@codemirror/search';
import { EditorState } from '@codemirror/state';
import {
    EditorView,
    crosshairCursor,
    drawSelection,
    dropCursor,
    highlightActiveLine,
    highlightSpecialChars,
    keymap,
    rectangularSelection,
} from '@codemirror/view';

// Setup without line numbers (based on basicSetup)
const setupWithoutLineNumbers = [
    highlightSpecialChars(),
    history(),
    foldGutter(),
    drawSelection(),
    dropCursor(),
    EditorState.allowMultipleSelections.of(true),
    indentOnInput(),
    syntaxHighlighting(defaultHighlightStyle, { fallback: true }),
    bracketMatching(),
    closeBrackets(),
    autocompletion(),
    rectangularSelection(),
    crosshairCursor(),
    highlightActiveLine(),
    highlightSelectionMatches(),
    keymap.of([
        ...closeBracketsKeymap,
        ...defaultKeymap,
        ...searchKeymap,
        ...historyKeymap,
        ...foldKeymap,
        ...completionKeymap,
    ]),
];

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
                setupWithoutLineNumbers,
                markdown({
                    base: markdownLanguage,
                    codeLanguages: languages,
                }),
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
