import { Box, Divider, Grid, Paper, Typography } from "@mui/material";
import { useNavigate } from "react-router-dom";
import { createAdminApiClient } from "../admin_api";
import type {
	EntryWithImage,
	GetLatestEntriesRow,
	LinkPalletData,
} from "../generated-client/model";
import AdminEntryCardItem from "./AdminEntryCardItem";
import CardItem from "./CardItem";

interface LinkPalletProps {
	linkPallet: LinkPalletData;
}

const api = createAdminApiClient();

// Convert EntryWithImage to GetLatestEntriesRow format
function convertToGetLatestEntriesRow(
	entry: EntryWithImage,
): GetLatestEntriesRow {
	return {
		Path: entry.path,
		Title: entry.title,
		Body: entry.body,
		Visibility: entry.visibility,
		Format: entry.format,
		ImageUrl: entry.imageUrl,
		PublishedAt: null,
		LastEditedAt: null,
		CreatedAt: null,
		UpdatedAt: null,
	};
}

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
		<Paper sx={{ p: 2 }}>
			<Typography variant="h6" gutterBottom>
				Related Links
			</Typography>

			{linkPallet.links.length > 0 && (
				<Box sx={{ mb: 3 }}>
					<Typography variant="subtitle2" color="text.secondary" gutterBottom>
						Direct Links
					</Typography>
					<Grid container spacing={1} sx={{ maxWidth: "540px" }}>
						{linkPallet.links.map((link) => (
							<Grid item key={link.path} xs={12} sm={6} md={4}>
								<AdminEntryCardItem
									entry={convertToGetLatestEntriesRow(link)}
								/>
							</Grid>
						))}
					</Grid>
				</Box>
			)}

			{linkPallet.twohops.length > 0 && (
				<>
					<Divider sx={{ my: 2 }} />
					<Box sx={{ mb: 3 }}>
						<Typography variant="subtitle2" color="text.secondary" gutterBottom>
							Two-hop Links
						</Typography>
						{linkPallet.twohops.map((twohops) => (
							<Box
								key={`${twohops.src.path || twohops.src.dstTitle}-twohop`}
								sx={{ mb: 2 }}
							>
								<Grid container spacing={1} sx={{ maxWidth: "540px" }}>
									<Grid item xs={12} sm={6} md={4}>
										{twohops.src.title ? (
											<AdminEntryCardItem
												entry={convertToGetLatestEntriesRow(twohops.src)}
												backgroundColor="#c8e6c9"
											/>
										) : (
											<CardItem
												onClick={() => createNewEntry(twohops.src.dstTitle)}
												title={twohops.src.dstTitle}
												content=""
												backgroundColor="#e0f2f1"
												color="gray"
											/>
										)}
									</Grid>
									{twohops.links.map((link) => (
										<Grid item key={link.path} xs={12} sm={6} md={4}>
											<AdminEntryCardItem
												entry={convertToGetLatestEntriesRow(link)}
											/>
										</Grid>
									))}
								</Grid>
							</Box>
						))}
					</Box>
				</>
			)}

			{linkPallet.newLinks.length > 0 && (
				<>
					<Divider sx={{ my: 2 }} />
					<Box>
						<Typography variant="subtitle2" color="text.secondary" gutterBottom>
							New Links to Create
						</Typography>
						<Grid container spacing={1} sx={{ maxWidth: "540px" }}>
							{linkPallet.newLinks.map((title) => (
								<Grid item key={title} xs={12} sm={6} md={4}>
									<CardItem
										onClick={() => createNewEntry(title)}
										title={title}
										content=""
										backgroundColor="#fff3e0"
										color="text.secondary"
									/>
								</Grid>
							))}
						</Grid>
					</Box>
				</>
			)}
		</Paper>
	);
}
