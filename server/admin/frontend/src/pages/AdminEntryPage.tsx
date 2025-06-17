import React, { useEffect, useState } from "react";
import { useLocation } from "react-router-dom";
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

	const showUpdatedMessage = React.useCallback((text: string) => {
		setUpdatedMessage(text);
		setTimeout(() => {
			setUpdatedMessage("");
		}, 1000);
	}, []);

	const clearMessage = React.useCallback(() => {
		setMessage("");
		setMessageType("");
	}, []);

	const showMessage = React.useCallback((type: "success" | "error", text: string) => {
		setMessageType(type);
		setMessage(text);
		setTimeout(() => {
			setMessage("");
			setMessageType("");
		}, 5000);
	}, []);

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

	async function handleDelete(event: React.MouseEvent) {
		event.preventDefault();

		const confirmed = confirm(
			`Are you sure you want to delete the entry "${title}"?`,
		);
		if (confirmed) {
			clearMessage();

			try {
				await api.deleteEntry({
					path: encodeURIComponent(entry.Path),
				});
				showMessage("success", "Entry deleted successfully");
				location.href = "/admin/";
			} catch (e) {
				console.log(e);
				showMessage("error", "Failed to delete entry");
			}
		}
	}

	async function handleRegenerateEntryImage(event: React.MouseEvent) {
		event.preventDefault();

		clearMessage();

		try {
			await api.regenerateEntryImage({
				path: encodeURIComponent(entry.Path),
			});
			showMessage("success", "Entry image regenerated successfully");
			location.href = "/admin/";
		} catch (e) {
			console.log(e);
			showMessage("error", "Failed to regenerate entry image");
		}
	}

	const handleUpdateBody = React.useCallback(async () => {
		clearMessage();

		if (body === "") {
			showMessage("error", "Body cannot be empty");
			return;
		}

		try {
			await api.updateEntryBody(
				{ path: encodeURIComponent(path) },
				{
					body: body,
				},
			);

			showUpdatedMessage("Updated");
			setIsDirty(false);
		} catch (e) {
			showMessage("error", "Failed to update entry body");
			console.error("Failed to update entry body:", e);
		}
	}, [body, path, clearMessage, showMessage, showUpdatedMessage]);

	const handleUpdateTitle = React.useCallback(async () => {
		clearMessage();
		if (title === "") {
			showMessage("error", "Title cannot be empty");
			return;
		}

		try {
			await api.updateEntryTitle(
				{ path: encodeURIComponent(path) },
				{
					title,
				},
			);

			showMessage("success", "Entry updated successfully");
			setIsDirty(false);
		} catch (e) {
			showMessage("error", "Failed to update entry title");
			console.error("Failed to update entry title:", e);
		}
	}, [title, path, clearMessage, showMessage]);

	const debouncedUpdateBody = React.useMemo(
		() => debounce(() => handleUpdateBody(), 800),
		[handleUpdateBody],
	);

	const debouncedTitleUpdate = React.useMemo(
		() => debounce(() => handleUpdateTitle(), 500),
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

	function toggleVisibility(event: React.MouseEvent) {
		event.preventDefault();
		event.stopPropagation();

		const newVisibility = visibility === "private" ? "public" : "private";

		if (
			!confirm("Are you sure you want to change the visibility of this entry?")
		) {
			return;
		}

		console.log("Updating visibility to", newVisibility);

		api
			.updateEntryVisibility(
				{ path: encodeURIComponent(entry.Path) },
				{
					visibility: newVisibility,
				},
			)
			.then((data) => {
				setVisibility(data.Visibility);
			})
			.catch((error) => {
				console.error("Failed to update visibility:", error);
				showMessage("error", `Failed to update visibility: ${error.message}`);
			});
	}

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
					location.href = "/admin/";
				}
			}
		};

		loadEntry();
		api.getLinkedEntryPaths({ path: encodeURIComponent(path) }).then(setLinks);
		loadLinks();
	}, [path, loadLinks, location]);

	const styles = {
		parent: {},
		container: {
			backgroundColor: entry.Visibility === "private" ? "#f3f4f6" : "white",
			minHeight: "100vh",
		},
		leftPane: {
			float: "left" as const,
			width: "49%",
			padding: "1rem",
		},
		rightPane: {
			float: "right" as const,
			width: "49%",
			padding: "1rem",
		},
		form: {},
		titleContainer: {
			marginBottom: "1rem",
		},
		input: {
			width: "100%",
			padding: "0.5rem",
			border: "1px solid #d1d5db",
			borderRadius: "0.25rem",
		},
		label: {
			display: "block",
			marginBottom: "0.5rem",
			fontWeight: "bold",
		},
		bodyContainer: {
			marginBottom: "1rem",
		},
		editor: {
			minHeight: "400px",
			border: "1px solid #d1d5db",
			borderRadius: "0.25rem",
		},
		visibilityContainer: {
			marginBottom: "1rem",
		},
		buttonContainer: {
			marginTop: "1rem",
			display: "flex",
			gap: "1rem",
		},
		deleteButton: {
			backgroundColor: "#ef4444",
			color: "white",
			padding: "0.5rem 1rem",
			border: "none",
			borderRadius: "0.25rem",
			cursor: "pointer",
		},
		regenerateButton: {
			backgroundColor: "#3b82f6",
			color: "white",
			padding: "0.5rem 1rem",
			border: "none",
			borderRadius: "0.25rem",
			cursor: "pointer",
		},
		linkContainer: {
			marginTop: "1rem",
		},
		link: {
			color: "#3b82f6",
			textDecoration: "underline",
		},
		updatedMessage: {
			position: "fixed" as const,
			top: "20px",
			right: "20px",
			backgroundColor: "#10b981",
			color: "white",
			padding: "0.5rem 1rem",
			borderRadius: "0.25rem",
		},
	};

	return (
		<div style={styles.parent}>
			<div style={styles.container}>
				<div style={styles.leftPane}>
					<form style={styles.form}>
						<div style={styles.titleContainer}>
							<input
								name="title"
								type="text"
								style={styles.input}
								value={title}
								onChange={(e) => {
									setTitle(e.target.value);
									setIsDirty(true);
									debouncedTitleUpdate();
								}}
								placeholder="Entry Title"
							/>
						</div>

						<div style={styles.bodyContainer}>
							<label htmlFor="body" style={styles.label}>
								Body
							</label>
							<div style={styles.editor}>
								<MarkdownEditor
									initialContent={body}
									onUpdateText={(text) => {
										setBody(text);
										handleInputBody();
									}}
								/>
							</div>
						</div>

						<div style={styles.visibilityContainer}>
							<label style={styles.label}>
								Visibility: {visibility}
								<button
									type="button"
									onClick={toggleVisibility}
									style={{ marginLeft: "1rem" }}
								>
									Toggle
								</button>
							</label>
						</div>
					</form>

					<div style={styles.buttonContainer}>
						<button
							type="submit"
							style={styles.deleteButton}
							onClick={handleDelete}
						>
							Delete
						</button>
						<button
							type="submit"
							style={styles.regenerateButton}
							onClick={handleRegenerateEntryImage}
						>
							Regenerate entry_image
						</button>
					</div>

					{visibility === "public" && (
						<div style={styles.linkContainer}>
							<a href={`/entry/${entry.Path}`} style={styles.link}>
								Go to User Side Page
							</a>
						</div>
					)}

					{updatedMessage !== "" && (
						<div style={styles.updatedMessage}>{updatedMessage}</div>
					)}
				</div>

				<div style={styles.rightPane}>
					<LinkPallet linkPallet={linkPallet} />
				</div>
			</div>
		</div>
	);
}
