import React from "react";
import type { GetLatestEntriesRow } from "../generated-client/model";
import CardItem from "./CardItem";

interface EntryCardItemProps {
	entry: GetLatestEntriesRow;
	backgroundColor?: string;
	color?: string;
	onClick: (event: React.MouseEvent) => void;
}

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

export default function EntryCardItem({
	entry,
	backgroundColor = entry.Visibility === "private" ? "#cccccc" : "#f6f6f6",
	color = "#0f0f0f",
	onClick,
}: EntryCardItemProps) {
	const title = entry.Title;
	const content = entry.Body
		? `${simplifyMarkdown(entry.Body).slice(0, 100)}...`
		: "";
	const imgSrc = entry.ImageUrl;

	return (
		<CardItem
			onClick={onClick}
			backgroundColor={backgroundColor}
			color={color}
			title={title}
			content={content}
			imgSrc={imgSrc}
		/>
	);
}
