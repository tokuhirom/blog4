import { format } from "date-fns";
import { useNavigate } from "react-router-dom";
import { createAdminApiClient } from "../admin_api";
import styles from "./AdminHeader.module.css";

const api = createAdminApiClient();

export default function AdminHeader() {
	const navigate = useNavigate();

	async function handleNewEntry() {
		try {
			// Generate a placeholder title with current date/time
			const now = new Date();
			const placeholderTitle = format(now, "'New Entry' yyyy-MM-dd HH-mm-ss");

			const data = await api.createEntry({ title: placeholderTitle });
			console.log(`New entry created: ${data.path}`);
			navigate(`/admin/entry/${data.path}`);
		} catch (e) {
			console.error("Error creating new entry:", e);
			alert("Failed to create new entry");
		}
	}

	return (
		<header className={styles.header}>
			<div className={styles.container}>
				<a href="/admin/" className={`${styles.link} ${styles.textXl}`}>
					Blog Admin
				</a>
				<nav className={styles.nav}>
					<button
						type="button"
						onClick={handleNewEntry}
						className={styles.button}
					>
						New Entry
					</button>
				</nav>
			</div>
		</header>
	);
}
