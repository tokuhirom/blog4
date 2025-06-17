import React, { useEffect, useRef, useState } from "react";
import { createAdminApiClient } from "../admin_api";
import AdminEntryCardItem from "../components/AdminEntryCardItem";
import SearchBox from "../components/SearchBox";
import type { GetLatestEntriesRow } from "../generated-client/model";

const api = createAdminApiClient();

export default function TopPage() {
	const [searchKeyword, setSearchKeyword] = useState("");
	const [allEntries, setAllEntries] = useState<GetLatestEntriesRow[]>([]);
	const [isLoading, setIsLoading] = useState(false);
	const [hasMore, setHasMore] = useState(true);
	const loadIntervalRef = useRef<NodeJS.Timeout | null>(null);

	const filteredEntries = React.useMemo(() => {
		if (searchKeyword === "") {
			return allEntries;
		}

		const lowerKeyword = searchKeyword.toLowerCase();
		return allEntries.filter(
			(entry) =>
				entry.Title?.toLowerCase()?.includes(lowerKeyword) ||
				entry.Body?.toLowerCase()?.includes(lowerKeyword),
		);
	}, [searchKeyword, allEntries]);

	function handleSearch(keyword: string) {
		setSearchKeyword(keyword);
	}

	async function loadMoreEntries() {
		if (!allEntries) {
			console.log("allEntries is not initialized yet");
			return;
		}
		console.log(`loadMoreEntries ${isLoading} ${hasMore} ${allEntries.length}`);
		if (isLoading || !hasMore) return;

		setIsLoading(true);

		const last_last_edited_at = allEntries[allEntries.length - 1]?.LastEditedAt;
		if (allEntries.length > 0 && !last_last_edited_at) {
			setIsLoading(false);
			setHasMore(false);
			return;
		}

		try {
			console.log(`loadMoreEntries ${last_last_edited_at}`);
			console.log(allEntries);
			const rawEntries = await api.getLatestEntries(
				last_last_edited_at
					? {
							lastLastEditedAt: last_last_edited_at,
						}
					: {},
			);
			// Filter out any entries without a Path (note: PascalCase from API)
			const newEntries = (rawEntries || []).filter((entry) => entry?.Path);

			if (newEntries.length === 0) {
				console.log(
					`No more entries to load... stopping loading more entries. last_last_edited_at=${last_last_edited_at}`,
				);
				setHasMore(false);
			} else {
				const existingPaths = allEntries.map((entry) => entry.Path);
				const addingNewEntries = newEntries.filter(
					(entry) => !existingPaths.includes(entry.Path),
				);
				if (addingNewEntries.length === 0) {
					console.log(
						`All entries are duplicated... stopping loading more entries. last_last_edited_at=${last_last_edited_at}, newEntries=${newEntries.map((entry) => entry.Title)}`,
					);
					setHasMore(false);
				} else {
					console.log(
						`Adding new entries... last_last_edited_at=${last_last_edited_at}, newEntries=${newEntries.map((entry) => entry.Title)}`,
					);
					setAllEntries((prev) => [...prev, ...addingNewEntries]);
				}
			}
		} catch (err) {
			setHasMore(false);
			console.error(err);
		} finally {
			setIsLoading(false);
		}
	}

	async function handleKeydown(event: KeyboardEvent) {
		if (
			event.key === "c" &&
			!event.ctrlKey &&
			!event.altKey &&
			!event.metaKey &&
			!event.shiftKey
		) {
			event.preventDefault();
			event.stopPropagation();
			try {
				const response = await fetch("/admin/api/entry", {
					method: "POST",
					headers: {
						"Content-Type": "application/json",
					},
					body: JSON.stringify({}),
				});
				if (response.ok) {
					const data = await response.json();
					location.href = `/admin/entry/${data.path}`;
				} else {
					alert(
						`Failed to create new entry: ${response.status} ${response.statusText}`,
					);
				}
			} catch (err) {
				console.error(err);
				alert(`Failed to create new entry: ${err}`);
			}
		}
	}

	useEffect(() => {
		window.addEventListener("keydown", handleKeydown);

		const loadInitialEntries = async () => {
			try {
				console.log("Loading entries...");
				setIsLoading(true);
				const entries = await api.getLatestEntries();
				console.log("Loaded entries", entries);
				// Filter out any entries without a Path (note: PascalCase from API)
				setAllEntries((entries || []).filter((entry) => entry?.Path));
				setIsLoading(false);

				console.log("Start loading more entries...");
				// start loading more entries
				loadIntervalRef.current = setInterval(() => {
					if (!isLoading && hasMore) {
						loadMoreEntries();
					}
				}, 10);
			} catch (err) {
				console.error(err);
				alert(`Failed to load entries: ${err}`);
			}
		};

		loadInitialEntries();

		return () => {
			if (loadIntervalRef.current) {
				clearInterval(loadIntervalRef.current);
			}
			window.removeEventListener("keydown", handleKeydown);
		};
	}, []);

	const styles = {
		container: {
			padding: "1rem",
			margin: "0 auto",
			maxWidth: "1200px",
		},
		loadingMessage: {
			marginTop: "1rem",
			textAlign: "center" as const,
			color: "#6b7280",
		},
		entryList: {
			display: "flex",
			flexWrap: "wrap" as const,
			margin: "auto",
			gap: "1rem",
			justifyContent: "flex-start",
			maxWidth: "1200px",
		},
	};

	return (
		<div style={styles.container}>
			<SearchBox onSearch={handleSearch} />

			<div style={styles.entryList}>
				{filteredEntries.map((entry, index) => (
					<AdminEntryCardItem
						key={entry.Path || `index-${index}`}
						entry={entry}
					/>
				))}
			</div>
			{(isLoading || hasMore) && (
				<p style={styles.loadingMessage}>Loading more entries...</p>
			)}
			{!hasMore && allEntries.length > 0 && (
				<p style={styles.loadingMessage}>No more entries to load</p>
			)}
		</div>
	);
}
