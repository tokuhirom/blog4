import React, { useState } from "react";
import { debounce } from "../utils";
import styles from "./SearchBox.module.css";

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

	return (
		<div className={styles.container}>
			<input
				type="text"
				placeholder="Search entries..."
				className={styles.input}
				value={keyword}
				onChange={handleInput}
			/>
		</div>
	);
}
