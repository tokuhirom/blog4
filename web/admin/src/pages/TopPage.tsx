import { Alert, Box, CircularProgress, Grid, Typography } from "@mui/material";
import { format } from "date-fns";
import React, { useEffect, useRef, useState } from "react";
import { useNavigate } from "react-router-dom";
import { createAdminApiClient } from "../admin_api";
import AdminEntryCardItem from "../components/AdminEntryCardItem";
import SearchBox from "../components/SearchBox";
import type { GetLatestEntriesRow } from "../generated-client/model";

const api = createAdminApiClient();

export default function TopPage() {
	const navigate = useNavigate();
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
			const lastEditedAt = lastEntry.LastEditedAt;

			console.log(
				`Loading more entries... last_last_edited_at=${lastEditedAt} title=${lastEntry.Title}`,
			);

			const rawEntries = await api.getLatestEntries({
				last_last_edited_at: lastEditedAt,
			});

			// Filter out any entries without a Path
			const newEntries = (rawEntries || []).filter((entry) => entry?.Path);

			if (newEntries.length === 0) {
				console.log("No more entries to load");
				setHasMore(false);
			} else {
				// Filter out duplicates
				const existingPaths = new Set(allEntries.map((entry) => entry.Path));
				const uniqueNewEntries = newEntries.filter(
					(entry) => !existingPaths.has(entry.Path),
				);

				if (uniqueNewEntries.length === 0) {
					console.log(
						`All entries are duplicates, no more to load: ${newEntries.length}, ${newEntries[0].Title}`,
					);
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

	const handleKeydown = React.useCallback(
		async (event: KeyboardEvent) => {
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
					// Generate a placeholder title with current date/time
					const now = new Date();
					const placeholderTitle = format(now, "yyyy-MM-ddTHH-mm-ss");

					const data = await api.createEntry({ title: placeholderTitle });
					console.log(`New entry created: ${data.path}`);
					navigate(`/entry/${data.path}`);
				} catch (err) {
					console.error("Error creating new entry:", err);
					alert(`Failed to create new entry: ${err}`);
				}
			}
		},
		[navigate],
	);

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
				console.log(
					`number of valid entries: ${validEntries?.length || 0} initial entries`,
				);
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
	}, [isInitialized, isLoading]);

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
	}, [isInitialized, hasMore, loadMoreEntries]);

	// Keyboard shortcuts
	useEffect(() => {
		window.addEventListener("keydown", handleKeydown);
		return () => {
			window.removeEventListener("keydown", handleKeydown);
		};
	}, [handleKeydown]);

	return (
		<Box sx={{ width: "100%" }}>
			<SearchBox onSearch={handleSearch} />

			<Grid container spacing={2}>
				{filteredEntries.map((entry) => (
					<Grid item xs={12} sm={6} md={4} lg={3} xl={2} key={entry.Path}>
						<AdminEntryCardItem entry={entry} />
					</Grid>
				))}
			</Grid>

			{isLoading && (
				<Box
					sx={{
						display: "flex",
						justifyContent: "center",
						alignItems: "center",
						mt: 4,
					}}
				>
					<CircularProgress />
					<Typography sx={{ ml: 2 }}>Loading more entries...</Typography>
				</Box>
			)}

			{!hasMore && allEntries.length > 0 && (
				<Box sx={{ display: "flex", justifyContent: "center", mt: 4 }}>
					<Typography color="text.secondary">
						No more entries to load
					</Typography>
				</Box>
			)}

			{!isLoading && allEntries.length === 0 && isInitialized && (
				<Box sx={{ display: "flex", justifyContent: "center", mt: 4 }}>
					<Alert severity="info">No entries found</Alert>
				</Box>
			)}
		</Box>
	);
}
