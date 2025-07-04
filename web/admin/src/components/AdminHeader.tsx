import AddIcon from "@mui/icons-material/Add";
import LogoutIcon from "@mui/icons-material/Logout";
import { AppBar, Button, Container, Toolbar, Typography } from "@mui/material";
import { format } from "date-fns";
import { useNavigate } from "react-router-dom";
import { createAdminApiClient } from "../admin_api";
import { useAuth } from "../hooks/useAuth";

const api = createAdminApiClient();

export default function AdminHeader() {
	const navigate = useNavigate();
	const { logout } = useAuth();

	async function handleNewEntry() {
		try {
			// Generate a placeholder title with current date/time
			const now = new Date();
			const placeholderTitle = format(now, "yyyy-MM-ddTHH-mm-ss");

			const data = await api.createEntry({ title: placeholderTitle });
			console.log(`New entry created: ${data.path}`);
			navigate(`/entry/${data.path}`);
		} catch (e) {
			console.error("Error creating new entry:", e);
			alert("Failed to create new entry");
		}
	}

	return (
		<AppBar position="fixed">
			<Container maxWidth="xl">
				<Toolbar disableGutters>
					<Typography
						variant="h6"
						component="a"
						href="/admin/"
						sx={{
							flexGrow: 1,
							textDecoration: "none",
							color: "inherit",
							fontWeight: 700,
						}}
					>
						Blog Admin
					</Typography>
					<Button
						color="inherit"
						onClick={handleNewEntry}
						startIcon={<AddIcon />}
						variant="outlined"
						sx={{ borderColor: "rgba(255, 255, 255, 0.5)", mr: 2 }}
					>
						New Entry
					</Button>
					<Button
						color="inherit"
						onClick={logout}
						startIcon={<LogoutIcon />}
						variant="outlined"
						sx={{ borderColor: "rgba(255, 255, 255, 0.5)" }}
					>
						Logout
					</Button>
				</Toolbar>
			</Container>
		</AppBar>
	);
}
