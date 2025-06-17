import type React from "react";
import type { GetLatestEntriesRow } from "../generated-client/model";
import EntryCardItem from "./EntryCardItem";

interface AdminEntryCardItemProps {
	entry: GetLatestEntriesRow;
	backgroundColor?: string;
	color?: string;
	onClick?: (event: React.MouseEvent) => void;
}

export default function AdminEntryCardItem({
	entry,
	backgroundColor = entry.Visibility === "private" ? "#cccccc" : "#f6f6f6",
	color = "#0f0f0f",
	onClick = (event: React.MouseEvent) => {
		if (event.metaKey || event.ctrlKey) {
			// Commandキー (Mac) または Ctrlキー (Windows/Linux) が押されている場合、別タブで開く
			window.open(`/admin/entry/${entry.Path}`, "_blank");
		} else {
			// 通常クリック時は同じタブで開く
			location.href = `/admin/entry/${entry.Path}`;
		}
	},
}: AdminEntryCardItemProps) {
	return (
		<EntryCardItem
			entry={entry}
			backgroundColor={backgroundColor}
			color={color}
			onClick={onClick}
		/>
	);
}
