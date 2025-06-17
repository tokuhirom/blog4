import { useEffect, useRef } from "react";
import { EditorState } from "@codemirror/state";
import { EditorView, keymap } from "@codemirror/view";
import { markdown } from "@codemirror/lang-markdown";
import { defaultKeymap, indentWithTab } from "@codemirror/commands";
import { oneDarkHighlightStyle } from "@codemirror/theme-one-dark";
import { syntaxHighlighting } from "@codemirror/language";

interface MarkdownEditorProps {
	initialContent?: string;
	onUpdateText?: (text: string) => void;
	onDropFiles?: (files: File[]) => Promise<string[]>;
	existsEntryByTitle?: (title: string) => boolean;
	findOrCreateEntry?: (title: string) => void;
}

export default function MarkdownEditor({
	initialContent = "",
	onUpdateText,
}: MarkdownEditorProps) {
	const containerRef = useRef<HTMLDivElement>(null);
	const editorRef = useRef<EditorView | null>(null);
	const onUpdateTextRef = useRef(onUpdateText);

	// Update the ref when onUpdateText changes
	useEffect(() => {
		onUpdateTextRef.current = onUpdateText;
	}, [onUpdateText]);

	// Initialize editor only once
	useEffect(() => {
		if (!containerRef.current) return;

		const state = EditorState.create({
			doc: initialContent,
			extensions: [
				markdown(),
				keymap.of([...defaultKeymap, indentWithTab]),
				syntaxHighlighting(oneDarkHighlightStyle),
				EditorView.updateListener.of((update) => {
					if (update.docChanged && onUpdateTextRef.current) {
						onUpdateTextRef.current(update.state.doc.toString());
					}
				}),
			],
		});

		const view = new EditorView({
			state,
			parent: containerRef.current,
		});

		editorRef.current = view;

		return () => {
			view.destroy();
		};
	}, []); // Empty dependency array - only run once on mount

	// Update editor content when initialContent changes (without recreating the editor)
	useEffect(() => {
		if (editorRef.current && initialContent !== undefined) {
			const currentContent = editorRef.current.state.doc.toString();
			if (currentContent !== initialContent) {
				editorRef.current.dispatch({
					changes: {
						from: 0,
						to: currentContent.length,
						insert: initialContent,
					},
				});
			}
		}
	}, [initialContent]);

	return (
		<div
			ref={containerRef}
			style={{ height: "400px", border: "1px solid #ccc" }}
		/>
	);
}