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
	const [isInitialized, setIsInitialized] = useState(false);
	const loadingRef = useRef(false);
	
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

	const loadMoreEntries = React.useCallback(async () => {
		// Prevent concurrent loads
		if (loadingRef.current || !hasMore) {
			return;
		}

		loadingRef.current = true;
		setIsLoading(true);

		try {
			const lastEntry = allEntries[allEntries.length - 1];
			if (!lastEntry) {
				console.error("No entries loaded yet, loading initial entries");
				return;
			}
			const lastEditedAt = lastEntry?.LastEditedAt;

			console.log(`Loading more entries... last_last_edited_at=${lastEditedAt}`);

			const rawEntries = await api.getLatestEntries(
				lastEditedAt ? { lastLastEditedAt: lastEditedAt } : {}
			);

			// Filter out any entries without a Path
			const newEntries = (rawEntries || []).filter((entry) => entry?.Path);

			if (newEntries.length === 0) {
				console.log("No more entries to load");
				setHasMore(false);
			} else {
				// Filter out duplicates
				const existingPaths = new Set(allEntries.map((entry) => entry.Path));
				const uniqueNewEntries = newEntries.filter(
					(entry) => !existingPaths.has(entry.Path)
				);

				if (uniqueNewEntries.length === 0) {
					console.log("All entries are duplicates, no more to load");
					setHasMore(false);
				} else {
					console.log(`Adding ${uniqueNewEntries.length} new entries`);
					setAllEntries((prev) => [...prev, ...uniqueNewEntries]);
				}
			}
		} catch (err) {
			console.error("Failed to load more entries:", err);
			setHasMore(false);
		} finally {
			loadingRef.current = false;
			setIsLoading(false);
		}
	}, [allEntries, hasMore]);

	const handleKeydown = React.useCallback(async (event: KeyboardEvent) => {
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
	}, []);

	// Initial load effect
	useEffect(() => {
		if (isInitialized) return;
		setIsInitialized(true);

		const loadInitial = async () => {
			if (isLoading) return;

			console.log("Loading initial entries...");
			setIsLoading(true);

			try {
				const entries = await api.getLatestEntries();

				console.log(`Loaded ${entries?.length || 0} initial entries`);
				const validEntries = (entries || []).filter((entry) => entry?.Path);
				console.log(`number of valid entries: ${validEntries?.length || 0} initial entries`);
				setAllEntries(validEntries);
				setIsInitialized(true);
				
				// Only set hasMore if we got a full page of results
				// Assuming the API returns 20-50 items per page
				setHasMore(validEntries.length >= 20);
			} catch (err) {
				console.error("Failed to load initial entries:", err);
				setIsInitialized(true);
				setHasMore(false);
			} finally {
				setIsLoading(false);
			}
		};

		loadInitial();
	}, [isInitialized]);

	// Load more entries using timeout
	useEffect(() => {
		if (!isInitialized || !hasMore || loadingRef.current) return;

		const timeoutId = setTimeout(() => {
			if (hasMore && !loadingRef.current) {
				loadMoreEntries();
			}
		}, 100);

		return () => {
			clearTimeout(timeoutId);
		};
	}, [isInitialized, hasMore, allEntries, loadMoreEntries]);

	// Keyboard shortcuts
	useEffect(() => {
		window.addEventListener("keydown", handleKeydown);
		return () => {
			window.removeEventListener("keydown", handleKeydown);
		};
	}, [handleKeydown]);

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
				{filteredEntries.map((entry) => (
					<AdminEntryCardItem
						key={entry.Path}
						entry={entry}
					/>
				))}
			</div>
			
			{isLoading && (
				<p style={styles.loadingMessage}>Loading more entries...</p>
			)}
			
			{!hasMore && allEntries.length > 0 && (
				<p style={styles.loadingMessage}>No more entries to load</p>
			)}
			
			{!isLoading && allEntries.length === 0 && isInitialized && (
				<p style={styles.loadingMessage}>No entries found</p>
			)}
		</div>
	);
}