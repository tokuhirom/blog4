import React, { useState } from "react";
import { debounce } from "../utils";

interface SearchBoxProps {
	onSearch: (keyword: string) => void;
}

export default function SearchBox({ onSearch }: SearchBoxProps) {
	const [keyword, setKeyword] = useState("");

	const debouncedSearch = React.useMemo(
		() =>
			debounce((newKeyword: string) => {
				onSearch(newKeyword);
			}, 1000),
		[onSearch],
	);

	function handleInput(event: React.ChangeEvent<HTMLInputElement>) {
		const newKeyword = event.target.value;
		setKeyword(newKeyword);
		debouncedSearch(newKeyword);
	}

	const styles = {
		searchInput: {
			marginBottom: "1rem",
			width: "100%",
			borderRadius: "0.375rem",
			border: "1px solid #d1d5db",
			padding: "0.5rem",
		},
	};

	return (
		<input
			type="text"
			placeholder="Search entries..."
			style={styles.searchInput}
			value={keyword}
			onChange={handleInput}
		/>
	);
}
