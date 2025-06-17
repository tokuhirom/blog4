<script lang="ts">
import type { GetLatestEntriesRow } from "../generated-client/model";
import EntryCardItem from "./EntryCardItem.svelte";

const {
	entry,
	backgroundColor = entry.Visibility === "private" ? "#cccccc" : "#f6f6f6",
	color = "#0f0f0f",
	onClick = (event: MouseEvent) => {
		if (event.metaKey || event.ctrlKey) {
			// Commandキー (Mac) または Ctrlキー (Windows/Linux) が押されている場合、別タブで開く
			window.open(`/admin/entry/${entry.Path}`, "_blank");
		} else {
			// 通常クリック時は同じタブで開く
			location.href = `/admin/entry/${entry.Path}`;
		}
	},
}: {
	entry: GetLatestEntriesRow;
	backgroundColor?: string;
	color?: string;
	onClick?: (event: MouseEvent) => void;
} = $props();
</script>

<EntryCardItem {onClick} {backgroundColor} {color} {entry} />
