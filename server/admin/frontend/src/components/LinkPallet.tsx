import React from "react";
import type { LinkPalletData } from "../generated-client/model";
import AdminEntryCardItem from "./AdminEntryCardItem";
import CardItem from "./CardItem";
import { createAdminApiClient } from "../admin_api";

interface LinkPalletProps {
	linkPallet: LinkPalletData;
}

const api = createAdminApiClient();

export default function LinkPallet({ linkPallet }: LinkPalletProps) {
	function createNewEntry(title: string) {
		api
			.createEntry({
				createEntryRequest: { title },
			})
			.then((data) => {
				location.href = `/admin/entry/${data.Path}`;
			})
			.catch((err) => {
				console.error("Error creating new entry:", err);
				alert("Failed to create new entry");
			});
	}

	const styles = {
		linkPallet: {},
		oneHopLink: {
			display: "flex",
			flexWrap: "wrap" as const,
			gap: "1rem",
			clear: "both" as const,
		},
		twoHopLink: {
			display: "flex",
			flexWrap: "wrap" as const,
			gap: "1rem",
			clear: "both" as const,
			marginTop: "1rem",
		},
	};

	return (
		<div style={styles.linkPallet}>
			<div style={styles.oneHopLink}>
				{linkPallet.links.map((link) => (
					<AdminEntryCardItem key={link.Path} entry={link} />
				))}
			</div>
			{linkPallet.twohops.map((twohops, index) => (
				<div key={index} style={styles.twoHopLink}>
					{twohops.src.Title ? (
						<AdminEntryCardItem
							entry={twohops.src}
							backgroundColor="yellowgreen"
						/>
					) : (
						<CardItem
							onClick={() => createNewEntry(twohops.src.dstTitle)}
							title={twohops.src.dstTitle}
							content=""
							backgroundColor="#c0f6f6"
							color="gray"
						/>
					)}
					{twohops.links.map((link) => (
						<AdminEntryCardItem key={link.Path} entry={link} />
					))}
				</div>
			))}
			{linkPallet.newLinks.length > 0 && (
				<div style={styles.oneHopLink}>
					<CardItem
						onClick={() => false}
						title="New Item"
						content=""
						backgroundColor="darkgoldenrod"
					/>
					{linkPallet.newLinks.map((title) => (
						<CardItem
							key={title}
							onClick={() => createNewEntry(title)}
							title={title}
							content=""
							color="gray"
						/>
					))}
				</div>
			)}
		</div>
	);
}
