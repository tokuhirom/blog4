import type React from "react";
import { useNavigate } from "react-router-dom";
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
	backgroundColor = entry.Visibility === "private" ? "#e0e0e0" : "#ffffff",
	color = "#0f0f0f",
	onClick,
}: AdminEntryCardItemProps) {
	const navigate = useNavigate();

	const handleClick = (event: React.MouseEvent) => {
		if (onClick) {
			onClick(event);
		} else {
			if (event.metaKey || event.ctrlKey) {
				// Commandキー (Mac) または Ctrlキー (Windows/Linux) が押されている場合、別タブで開く
				window.open(`/entry/${entry.Path}`, "_blank");
			} else {
				// 通常クリック時は同じタブで開く
				navigate(`/entry/${entry.Path}`);
			}
		}
	};

	return (
		<EntryCardItem
			entry={entry}
			backgroundColor={backgroundColor}
			color={color}
			onClick={handleClick}
		/>
	);
}
