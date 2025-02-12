<script lang="ts">
    import CardItem from './CardItem.svelte';
    import type {GetLatestEntriesRow} from "../generated-client";

	let {
		entry,
		backgroundColor = entry.visibility == 'private' ? '#cccccc' : '#f6f6f6',
		color = '#0f0f0f',
		onClick
	}: {
		entry: GetLatestEntriesRow;
		backgroundColor: string;
		color: string;
		onClick: (event: MouseEvent) => void;
	} = $props();

	function simplifyMarkdown(text: string): string {
		return text
			.replaceAll(/\n/g, ' ')
			.replaceAll(/\[(.*?)\]\(.*?\)/g, '$1')
			.replace(/\[\[(.*?)\]\]/g, '$1')
			.replace(/`.*?`/g, '')
			.replace(/#+/g, '')
			.replace(/\s+/g, ' ')
			.replace(/https?:\/\/\S+/g, ' ')
			.trim();
	}

	let title = entry.title;
	let content = entry.body ? simplifyMarkdown(entry.body).slice(0, 100) + '...' : '';
	let imgSrc = entry.imageUrl;
</script>

<CardItem {onClick} {backgroundColor} {color} {title} {content} {imgSrc} />
