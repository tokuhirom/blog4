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

	useEffect(() => {
		if (!containerRef.current) return;

		const state = EditorState.create({
			doc: initialContent,
			extensions: [
				markdown(),
				keymap.of([...defaultKeymap, indentWithTab]),
				syntaxHighlighting(oneDarkHighlightStyle),
				EditorView.updateListener.of((update) => {
					if (update.docChanged && onUpdateText) {
						onUpdateText(update.state.doc.toString());
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
	}, [initialContent, onUpdateText]);

	return (
		<div
			ref={containerRef}
			style={{ height: "400px", border: "1px solid #ccc" }}
		/>
	);
}
