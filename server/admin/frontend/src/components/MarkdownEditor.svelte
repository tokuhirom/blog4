<script lang="ts">
import { onMount } from "svelte";
import { EditorView, keymap } from "@codemirror/view";
import { EditorState, Transaction } from "@codemirror/state";
import { markdown, markdownLanguage } from "@codemirror/lang-markdown";
import { history, historyKeymap, indentWithTab } from "@codemirror/commands";
import { languages } from "@codemirror/language-data";
import { oneDarkHighlightStyle } from "@codemirror/theme-one-dark";
import { syntaxHighlighting } from "@codemirror/language";
import { defaultKeymap } from "@codemirror/commands";
import { internalLinkPlugin } from "./markdown/InternalLink";
import {
	autocompletion,
	type CompletionContext,
} from "@codemirror/autocomplete";

let container: HTMLDivElement;

const {
	initialContent = "",
	onUpdateText,
	onSave = () => {},
	existsEntryByTitle,
	onClickEntry,
	content = initialContent,
	pageTitles = [],
}: {
	initialContent: string;
	onUpdateText: (content: string) => void;
	onSave: (content: string) => void;
	existsEntryByTitle: (title: string) => boolean;
	onClickEntry: (title: string) => void;
	content: string;
	pageTitles: string[];
} = $props();

let editor: EditorView;

let isUploading: boolean = $state(false); // アップロード中の状態
let errorMessage: string = $state(""); // エラー通知メッセージ

const myCompletion = (context: CompletionContext) => {
	{
		// `[[foobar]]` style notation
		const word = context.matchBefore(/\[\[(?:(?!\]\].*).)*/);
		if (word) {
			console.log("Return links");
			const options = pageTitles.map((title) => {
				return {
					label: `[[${title}]]`,
					type: "keyword",
				};
			});
			return {
				from: word.from,
				options: options,
			};
		}
	}
	return null;
};

onMount(() => {
	const startState = EditorState.create({
		doc: initialContent,
		extensions: [
			internalLinkPlugin(existsEntryByTitle, onClickEntry),
			markdown({
				base: markdownLanguage,
				codeLanguages: languages,
			}),
			history(),
			syntaxHighlighting(oneDarkHighlightStyle),
			autocompletion({ override: [myCompletion], closeOnBlur: false }),
			EditorView.lineWrapping,
			keymap.of([...historyKeymap, ...defaultKeymap, indentWithTab]),
			EditorView.theme({
				".cm-editor": { height: "600px" },
				".cm-content": { overflowY: "auto" },
			}),
			EditorView.domEventHandlers({
				paste: handlePaste,
			}),
			EditorView.updateListener.of((update) => {
				if (update.changes) {
					const isUserInput = update.transactions.some(
						(tr) =>
							tr.annotation(Transaction.userEvent) !== "program" &&
							tr.docChanged,
					);
					if (isUserInput && onUpdateText) {
						onUpdateText(editor.state.doc.toString());
					}
				}
			}),
		],
	});

	editor = new EditorView({
		state: startState,
		parent: container,
	});

	$effect(() => {
		console.log("content:", content);
		if (editor && content !== editor.state.doc.toString()) {
			const transaction = editor.state.update({
				changes: {
					from: 0,
					to: editor.state.doc.length,
					insert: content,
				},
			});
			editor.dispatch(transaction);
		}
	});

	return () => {
		editor.destroy();
	};
});

function handlePaste(event: ClipboardEvent): void {
	const items = event.clipboardData?.items;
	if (!items) return;

	for (const item of items) {
		if (item.type.startsWith("image/")) {
			event.preventDefault();

			const file = item.getAsFile();
			if (file) {
				isUploading = true; // アップロード中フラグを立てる
				uploadImage(file)
					.then((url) => {
						insertMarkdownImage(url);
					})
					.catch((error) => {
						console.error("Image upload failed:", error);
						showError("Image upload failed. Please try again.");
					})
					.finally(() => {
						isUploading = false; // アップロード完了でフラグを下ろす
					});
			}
		}
	}
}

async function uploadImage(file: File): Promise<string> {
	const formData = new FormData();
	formData.append("file", file);

	const response = await fetch("/admin/api/upload", {
		method: "POST",
		body: formData,
	});

	if (!response.ok) {
		throw new Error("Failed to upload image");
	}

	const data = await response.json();
	return data.url;
}

function insertMarkdownImage(url: string) {
	const markdownImage = `![Image](${url})`;
	const transaction = editor.state.update({
		changes: {
			from: editor.state.selection.main.from,
			insert: markdownImage,
		},
	});
	editor.dispatch(transaction);
}

function showError(message: string) {
	errorMessage = message;
	setTimeout(() => {
		errorMessage = ""; // 7秒後にエラーメッセージを非表示
	}, 7000);
}

const handleKeyDown = (event: KeyboardEvent) => {
	if ((event.ctrlKey || event.metaKey) && event.key === "s") {
		event.preventDefault();
		onSave(editor.state.doc.toString());
	}
};
</script>

<svelte:window on:keydown={handleKeyDown} />

<div class="wrapper">
	<div bind:this={container}></div>
	{#if isUploading}
		<div class="upload-indicator">Uploading image...</div>
	{/if}
	{#if errorMessage}
		<div class="error-indicator">{errorMessage}</div>
	{/if}
</div>

<style>
	.wrapper {
		width: 100%;
		height: 100%;
	}

	.upload-indicator {
		position: fixed;
		bottom: 20px;
		right: 20px;
		background: rgba(0, 0, 0, 0.7);
		color: white;
		padding: 10px 15px;
		border-radius: 5px;
		font-size: 14px;
	}

	.error-indicator {
		position: fixed;
		bottom: 60px;
		right: 20px;
		background: rgba(255, 0, 0, 0.8);
		color: white;
		padding: 10px 15px;
		border-radius: 5px;
		font-size: 14px;
	}
</style>
