import { useNavigate } from "react-router-dom";
import { createAdminApiClient } from "../admin_api";

const api = createAdminApiClient();

export default function AdminHeader() {
	const navigate = useNavigate();
	
	async function handleNewEntry() {
		try {
			// Generate a placeholder title with current date/time
			const now = new Date();
			const dateStr = now.toISOString().split('T')[0]; // YYYY-MM-DD
			const timeStr = now.toTimeString().split(' ')[0].replace(/:/g, '-'); // HH-MM-SS
			const placeholderTitle = `New Entry ${dateStr} ${timeStr}`;
			
			const data = await api.createEntry({ title: placeholderTitle });
			console.log(`New entry created: ${data.Path}`);
			navigate(`/admin/entry/${data.Path}`);
		} catch (e) {
			console.error("Error creating new entry:", e);
			alert("Failed to create new entry");
		}
	}

	const styles = {
		header: {
			position: "fixed" as const,
			left: 0,
			top: 0,
			zIndex: 10,
			width: "100%",
			backgroundColor: "#d92706",
			color: "white",
			height: "62px",
			verticalAlign: "middle",
			fontFamily: "'Hiragino Kaku Gothic ProN', 'Meiryo', sans-serif",
		},
		container: {
			maxWidth: "1200px",
			margin: "0 auto",
		},
		link: {
			display: "block",
			float: "left" as const,
			textDecoration: "none",
			color: "white",
			padding: "8px",
		},
		linkHover: {
			textDecoration: "underline",
			cursor: "pointer",
		},
		textXl: {
			fontSize: "1.25rem",
			fontWeight: "bold",
		},
		nav: {
			float: "right" as const,
		},
		button: {
			background: "none",
			border: "none",
			color: "white",
			textDecoration: "none",
			cursor: "pointer",
			font: "inherit",
		},
		buttonHover: {
			textDecoration: "underline",
		},
	};

	return (
		<header style={styles.header}>
			<div style={styles.container}>
				<a
					href="/admin/"
					style={{ ...styles.link, ...styles.textXl }}
					onMouseEnter={(e) => {
						e.currentTarget.style.textDecoration = "underline";
					}}
					onMouseLeave={(e) => {
						e.currentTarget.style.textDecoration = "none";
					}}
				>
					Blog Admin
				</a>
				<nav style={styles.nav}>
					<button
						type="button"
						onClick={handleNewEntry}
						style={styles.button}
						onMouseEnter={(e) => {
							e.currentTarget.style.textDecoration = "underline";
						}}
						onMouseLeave={(e) => {
							e.currentTarget.style.textDecoration = "none";
						}}
					>
						New Entry
					</button>
				</nav>
			</div>
		</header>
	);
}
