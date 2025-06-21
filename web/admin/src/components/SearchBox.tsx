import React, { useState } from "react";
import { TextField, InputAdornment } from "@mui/material";
import SearchIcon from "@mui/icons-material/Search";
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

	return (
		<TextField
			fullWidth
			placeholder="Search entries..."
			value={keyword}
			onChange={handleInput}
			variant="outlined"
			size="small"
			InputProps={{
				startAdornment: (
					<InputAdornment position="start">
						<SearchIcon />
					</InputAdornment>
				),
			}}
			sx={{ mb: 3 }}
		/>
	);
}
