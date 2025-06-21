import { useNavigate } from "react-router-dom";
import { Box, Paper, Typography, Divider, Grid } from "@mui/material";
import type { LinkPalletData } from "../generated-client/model";
import AdminEntryCardItem from "./AdminEntryCardItem";
import CardItem from "./CardItem";
import { createAdminApiClient } from "../admin_api";

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
							<Grid item key={link.Path} xs={12} sm={6} md={4}>
								<AdminEntryCardItem entry={link} />
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
								key={`${twohops.src.Path || twohops.src.dstTitle}-twohop`}
								sx={{ mb: 2 }}
							>
								<Grid container spacing={1} sx={{ maxWidth: "540px" }}>
									<Grid item xs={12} sm={6} md={4}>
										{twohops.src.Title ? (
											<AdminEntryCardItem
												entry={twohops.src}
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
										<Grid item key={link.Path} xs={12} sm={6} md={4}>
											<AdminEntryCardItem entry={link} />
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
