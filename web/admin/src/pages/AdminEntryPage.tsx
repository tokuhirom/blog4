import React, { useEffect, useRef, useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import {
	Box,
	TextField,
	Button,
	Paper,
	Typography,
	Alert,
	Snackbar,
	FormControl,
	FormLabel,
	RadioGroup,
	FormControlLabel,
	Radio,
	Link,
} from "@mui/material";
import Grid from "@mui/material/Grid2";
import DeleteIcon from "@mui/icons-material/Delete";
import RefreshIcon from "@mui/icons-material/Refresh";
import { createAdminApiClient } from "../admin_api";
import LinkPallet from "../components/LinkPallet";
import MarkdownEditor from "../components/MarkdownEditor";
import { extractLinks } from "../extractLinks";
import type {
	GetLatestEntriesRow,
	LinkPalletData,
} from "../generated-client/model";
import { debounce } from "../utils";

const api = createAdminApiClient();

export default function AdminEntryPage() {
	const location = useLocation();
	const navigate = useNavigate();
	const path = location.pathname.replace("/admin/entry/", "");

	const [entry, setEntry] = useState<GetLatestEntriesRow>(
		{} as GetLatestEntriesRow,
	);
	const [, setLinks] = useState<{ [key: string]: string | null }>({});
	const [title, setTitle] = useState("");
	const [body, setBody] = useState("");
	const [visibility, setVisibility] = useState("private");
	const [currentLinks, setCurrentLinks] = useState<string[]>([]);
	const [linkPallet, setLinkPallet] = useState<LinkPalletData>({
		links: [],
		twohops: [],
		newLinks: [],
	});
	const [, setMessage] = useState("");
	const [, setMessageType] = useState<"success" | "error" | "">("");
	const [updatedMessage, setUpdatedMessage] = useState("");
	const [, setIsDirty] = useState(false);
	const titleRef = useRef(title);
	const bodyRef = useRef(body);

	// Update refs when state changes
	useEffect(() => {
		titleRef.current = title;
	}, [title]);

	useEffect(() => {
		bodyRef.current = body;
	}, [body]);

	const showUpdatedMessage = React.useCallback((text: string) => {
		console.log("Showing updated message:", text);
		setUpdatedMessage(text);
		setTimeout(() => {
			setUpdatedMessage("");
		}, 1000);
	}, []);

	const clearMessage = React.useCallback(() => {
		setMessage("");
		setMessageType("");
	}, []);

	const showMessage = React.useCallback(
		(type: "success" | "error", text: string) => {
			setMessageType(type);
			setMessage(text);
			setTimeout(() => {
				setMessage("");
				setMessageType("");
			}, 5000);
		},
		[],
	);

	const loadLinks = React.useCallback(() => {
		api
			.getLinkPallet({ path: encodeURIComponent(path) })
			.then((data) => {
				console.log("Got link pallet data", data);
				setLinkPallet(data);
			})
			.catch((error) => {
				console.error("Failed to get links:", error);
			});
	}, [path]);

	const handleDelete = React.useCallback(
		async (event: React.MouseEvent) => {
			event.preventDefault();

			const confirmed = confirm(
				`Are you sure you want to delete the entry "${title}"?`,
			);
			if (confirmed) {
				clearMessage();

				try {
					console.log("Deleting entry:", entry.Path);
					await api.deleteEntry({
						path: encodeURIComponent(entry.Path),
					});
					console.log("Entry deleted:", entry.Path);
					showMessage("success", "Entry deleted successfully");
					// Small delay to ensure the message is shown
					setTimeout(() => {
						navigate("/admin/");
					}, 500);
				} catch (e) {
					console.log(e);
					showMessage("error", "Failed to delete entry");
				}
			}
		},
		[entry.Path, title, clearMessage, showMessage, navigate],
	);

	const handleRegenerateEntryImage = React.useCallback(
		async (event: React.MouseEvent) => {
			event.preventDefault();
			console.log("Regenerating entry image for:", entry.Path);

			if (!entry.Path) {
				showMessage("error", "No entry path available");
				return;
			}

			clearMessage();

			try {
				await api.regenerateEntryImage({
					path: encodeURIComponent(entry.Path),
				});
				showUpdatedMessage("Entry image regenerated successfully");
				// Don't redirect, just show the message
				// location.href = "/admin/";
			} catch (e) {
				console.error("Failed to regenerate entry image:", e);
				showMessage("error", "Failed to regenerate entry image");
			}
		},
		[entry.Path, clearMessage, showMessage, showUpdatedMessage],
	);

	const handleUpdateBody = React.useCallback(async () => {
		clearMessage();

		const currentBody = bodyRef.current;
		if (currentBody === "") {
			showMessage("error", "Body cannot be empty");
			return;
		}

		try {
			await api.updateEntryBody(
				{ path: encodeURIComponent(path) },
				{
					body: currentBody,
				},
			);

			showUpdatedMessage("Updated");
			setIsDirty(false);
		} catch (e) {
			showMessage("error", "Failed to update entry body");
			console.error("Failed to update entry body:", e);
		}
	}, [path, clearMessage, showMessage, showUpdatedMessage]);

	const handleUpdateTitle = React.useCallback(async () => {
		clearMessage();
		const currentTitle = titleRef.current;
		if (currentTitle === "") {
			showMessage("error", "Title cannot be empty");
			return;
		}

		try {
			await api.updateEntryTitle(
				{ path: encodeURIComponent(path) },
				{
					title: currentTitle,
				},
			);

			showUpdatedMessage("Updated");
			setIsDirty(false);
		} catch (e) {
			showMessage("error", "Failed to update entry title");
			console.error("Failed to update entry title:", e);
		}
	}, [path, clearMessage, showMessage, showUpdatedMessage]);

	const debouncedUpdateBody = React.useMemo(
		() => debounce(handleUpdateBody, 800),
		[handleUpdateBody],
	);

	const debouncedTitleUpdate = React.useMemo(
		() => debounce(handleUpdateTitle, 500),
		[handleUpdateTitle],
	);

	const handleInputBody = React.useCallback(() => {
		setIsDirty(true);
		debouncedUpdateBody();

		const newLinks = extractLinks(body);
		if (JSON.stringify(currentLinks) !== JSON.stringify(newLinks)) {
			setCurrentLinks(newLinks);
			loadLinks();
		}
	}, [body, currentLinks, debouncedUpdateBody, loadLinks]);

	const handleDropFiles = React.useCallback(
		async (files: File[]): Promise<string[]> => {
			console.log("Uploading files:", files);
			const urls: string[] = [];

			for (const file of files) {
				try {
					const formData = new FormData();
					formData.append("file", file);

					const response = await api.uploadFile(formData);
					console.log("Upload response:", response);

					if (response.url) {
						urls.push(response.url);
					}
				} catch (err) {
					console.error("Failed to upload file:", file.name, err);
					showMessage("error", `Failed to upload ${file.name}`);
				}
			}

			return urls;
		},
		[showMessage],
	);

	useEffect(() => {
		const loadEntry = async () => {
			try {
				const loadedEntry = await api.getEntryByDynamicPath({
					path: encodeURIComponent(path),
				});
				setEntry(loadedEntry);
				setTitle(loadedEntry.Title);
				setBody(loadedEntry.Body);
				setVisibility(loadedEntry.Visibility);
				setCurrentLinks(extractLinks(loadedEntry.Body));
			} catch (e) {
				console.error("Failed to get entry:", e);
				if (e instanceof Error && e.message.includes("404")) {
					navigate("/admin/");
				}
			}
		};

		loadEntry();
		api.getLinkedEntryPaths({ path: encodeURIComponent(path) }).then(setLinks);
		loadLinks();
	}, [path, loadLinks, navigate]);

	return (
		<Box>
			<Grid container spacing={3}>
				<Grid size={8}>
					<Paper
						sx={{
							p: 3,
							backgroundColor:
								entry.Visibility === "private" ? "#f5f5f5" : "background.paper",
						}}
					>
						<TextField
							fullWidth
							label="Title"
							value={title}
							onChange={(e) => {
								setTitle(e.target.value);
								setIsDirty(true);
								debouncedTitleUpdate();
							}}
							placeholder="Entry Title"
							variant="outlined"
							sx={{ mb: 3 }}
						/>

						<Box sx={{ mb: 3 }}>
							<Typography variant="h6" gutterBottom>
								Body
							</Typography>
							<Paper
								variant="outlined"
								sx={{
									p: 0,
									height: "700px",
									maxHeight: "70vh",
									overflow: "hidden",
									position: "relative",
								}}
							>
								<MarkdownEditor
									key={path}
									initialContent={entry.Body}
									onUpdateText={(text) => {
										setBody(text);
										handleInputBody();
									}}
									onDropFiles={handleDropFiles}
								/>
							</Paper>
						</Box>
					</Paper>
				</Grid>

				<Grid size={4}>
					<Box sx={{ display: "flex", flexDirection: "column", gap: 2 }}>
						{/* Entry Controls */}
						<Paper sx={{ p: 2 }}>
							<FormControl component="fieldset" sx={{ width: "100%", mb: 2 }}>
								<FormLabel component="legend">Visibility</FormLabel>
								<RadioGroup
									value={visibility}
									onChange={(e) => {
										if (
											confirm(
												"Are you sure you want to change the visibility of this entry?",
											)
										) {
											const newVisibility = e.target.value;
											api
												.updateEntryVisibility(
													{ path: encodeURIComponent(entry.Path) },
													{ visibility: newVisibility },
												)
												.then((data) => {
													setVisibility(data.Visibility);
												})
												.catch((error) => {
													console.error("Failed to update visibility:", error);
													showMessage(
														"error",
														`Failed to update visibility: ${error.message}`,
													);
												});
										}
									}}
								>
									<FormControlLabel
										value="private"
										control={<Radio />}
										label="Private"
									/>
									<FormControlLabel
										value="public"
										control={<Radio />}
										label="Public"
									/>
								</RadioGroup>
							</FormControl>

							<Box sx={{ display: "flex", flexDirection: "column", gap: 1 }}>
								<Button
									variant="outlined"
									color="error"
									startIcon={<DeleteIcon />}
									onClick={handleDelete}
									fullWidth
								>
									Delete
								</Button>
								<Button
									variant="outlined"
									startIcon={<RefreshIcon />}
									onClick={handleRegenerateEntryImage}
									fullWidth
								>
									Regenerate entry_image
								</Button>
							</Box>

							{visibility === "public" && (
								<Box sx={{ mt: 2 }}>
									<Link
										href={`/entry/${entry.Path}`}
										target="_blank"
										rel="noopener"
									>
										Go to User Side Page
									</Link>
								</Box>
							)}
						</Paper>

						{/* Link Pallet */}
						<LinkPallet linkPallet={linkPallet} />
					</Box>
				</Grid>
			</Grid>

			<Snackbar
				open={updatedMessage !== ""}
				autoHideDuration={1000}
				onClose={() => setUpdatedMessage("")}
				anchorOrigin={{ vertical: "bottom", horizontal: "right" }}
			>
				<Alert severity="success" sx={{ width: "100%" }}>
					{updatedMessage}
				</Alert>
			</Snackbar>
		</Box>
	);
}
