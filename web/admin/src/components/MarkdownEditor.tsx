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
	onDropFiles,
}: MarkdownEditorProps) {
	const containerRef = useRef<HTMLDivElement>(null);
	const editorRef = useRef<EditorView | null>(null);
	const onUpdateTextRef = useRef(onUpdateText);
	const onDropFilesRef = useRef(onDropFiles);

	// Update the refs when callbacks change
	useEffect(() => {
		onUpdateTextRef.current = onUpdateText;
	}, [onUpdateText]);

	useEffect(() => {
		onDropFilesRef.current = onDropFiles;
	}, [onDropFiles]);

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
				EditorView.domEventHandlers({
					paste: (event, view) => {
						console.log("Paste event triggered", event);
						const items = event.clipboardData?.items;
						if (!items) return false;

						const files: File[] = [];
						for (let i = 0; i < items.length; i++) {
							const item = items[i];
							console.log("Clipboard item:", item.type);
							if (item.type.indexOf("image") !== -1) {
								const file = item.getAsFile();
								if (file) {
									files.push(file);
								}
							}
						}

						if (files.length > 0 && onDropFilesRef.current) {
							event.preventDefault();
							const pos = view.state.selection.main.head;
							onDropFilesRef
								.current(files)
								.then((urls) => {
									const text = urls.map((url) => `![image](${url})`).join("\n");
									view.dispatch({
										changes: { from: pos, insert: text },
										selection: { anchor: pos + text.length },
									});
								})
								.catch((err) => {
									console.error("Failed to upload files:", err);
								});
							return true;
						}
						return false;
					},
					drop: (event, view) => {
						event.preventDefault();
						const files = Array.from(event.dataTransfer?.files || []);
						if (files.length > 0 && onDropFilesRef.current) {
							const pos = view.posAtCoords({
								x: event.clientX,
								y: event.clientY,
							});
							if (pos !== null) {
								onDropFilesRef
									.current(files)
									.then((urls) => {
										const text = urls
											.map((url) => `![image](${url})`)
											.join("\n");
										view.dispatch({
											changes: { from: pos, insert: text },
											selection: { anchor: pos + text.length },
										});
									})
									.catch((err) => {
										console.error("Failed to upload files:", err);
									});
							}
						}
						return true;
					},
					dragover: (event) => {
						event.preventDefault();
						return true;
					},
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
	}, [initialContent]);

	return (
		<div
			ref={containerRef}
			style={{
				height: "400px",
				border: "1px solid #e0e0e0",
				borderRadius: "4px",
			}}
		/>
	);
}
