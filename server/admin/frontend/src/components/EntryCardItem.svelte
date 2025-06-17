<script lang="ts">
import CardItem from "./CardItem.svelte";
import type { GetLatestEntriesRow } from "../generated-client/model";

const {
	entry,
	backgroundColor = entry.Visibility === "private" ? "#cccccc" : "#f6f6f6",
	color = "#0f0f0f",
	onClick,
}: {
	entry: GetLatestEntriesRow;
	backgroundColor: string;
	color: string;
	onClick: (event: MouseEvent) => void;
} = $props();

function simplifyMarkdown(text: string): string {
	return text
		.replaceAll(/\n/g, " ")
		.replaceAll(/\[(.*?)]\(.*?\)/g, "$1")
		.replace(/\[\[(.*?)]]/g, "$1")
		.replace(/`.*?`/g, "")
		.replace(/#+/g, "")
		.replace(/\s+/g, " ")
		.replace(/https?:\/\/\S+/g, " ")
		.trim();
}

const title = entry.Title;
const content = entry.Body
	? `${simplifyMarkdown(entry.Body).slice(0, 100)}...`
	: "";
const imgSrc = entry.ImageUrl;
</script>

<CardItem {onClick} {backgroundColor} {color} {title} {content} {imgSrc} />
