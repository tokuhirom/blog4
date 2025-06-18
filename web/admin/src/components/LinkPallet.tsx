import { useNavigate } from "react-router-dom";
import type { LinkPalletData } from "../generated-client/model";
import AdminEntryCardItem from "./AdminEntryCardItem";
import CardItem from "./CardItem";
import { createAdminApiClient } from "../admin_api";
import styles from "./LinkPallet.module.css";

interface LinkPalletProps {
	linkPallet: LinkPalletData;
}

const api = createAdminApiClient();

export default function LinkPallet({ linkPallet }: LinkPalletProps) {
	const navigate = useNavigate();

	function createNewEntry(title: string) {
		api
			.createEntry({ title })
			.then((data) => {
				navigate(`/admin/entry/${data.Path}`);
			})
			.catch((err) => {
				console.error("Error creating new entry:", err);
				alert("Failed to create new entry");
			});
	}

	return (
		<div>
			<div className={styles.oneHopLink}>
				{linkPallet.links.map((link) => (
					<AdminEntryCardItem key={link.Path} entry={link} />
				))}
			</div>
			{linkPallet.twohops.map((twohops) => (
				<div
					key={`${twohops.src.Path || twohops.src.dstTitle}-twohop`}
					className={styles.twoHopLink}
				>
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
				<div className={styles.oneHopLink}>
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
