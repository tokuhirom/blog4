import {
	Decoration,
	type DecorationSet,
	EditorView,
	ViewPlugin,
	ViewUpdate
} from '@codemirror/view';
import { type Extension, RangeSetBuilder } from '@codemirror/state';

export function internalLinkPlugin(
	existsEntryByTitle: (pageName: string) => boolean,
	findOrCreateEntry: (pageName: string) => void
): Extension {
	return ViewPlugin.fromClass(
		class {
			decorations: DecorationSet;

			constructor(view: EditorView) {
				this.decorations = this.buildDecorations(view);
			}

			update(update: ViewUpdate) {
				if (update.docChanged || update.viewportChanged) {
					this.decorations = this.buildDecorations(update.view);
				}
			}

			buildDecorations(view: EditorView): DecorationSet {
				const builder = new RangeSetBuilder<Decoration>();
				const re = /\[\[([^|`]+?)]]/g; // 内部リンクの正規表現
				for (const { from, to } of view.visibleRanges) {
					const text = view.state.doc.sliceString(from, to);
					let match;
					while ((match = re.exec(text))) {
						const pos = from + match.index;
						const exists = !!existsEntryByTitle(match[1]);
						builder.add(
							pos,
							pos + match[0].length,
							Decoration.mark({
								class: `internal-link ${exists ? 'exists' : 'not-exists'}`
							})
						);
					}
				}
				return builder.finish();
			}
		},
		{
			decorations: (v) => v.decorations,
			eventHandlers: {
				click: (event) => {
					const linkElement = (event.target as HTMLElement).closest('.internal-link');
					if (linkElement) {
						const inlineLink = linkElement.textContent;
						// ここでキーワードに応じたアクションを実行
						console.log(`Keyword clicked: '${inlineLink}'`);
						if (inlineLink) {
							const m = inlineLink.match(/^\[\[(.+)]]$/);
							if (m) {
								event.preventDefault();
								event.stopPropagation();

								const keyword = m[1];

								console.log('Keyword clicked: ' + keyword);
								findOrCreateEntry(keyword);
							}
						}
					}
				}
			}
		}
	);
}
