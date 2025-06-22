import { Box, Typography, Button } from "@mui/material";
import { useNavigate } from "react-router-dom";

export default function NotFound() {
	const navigate = useNavigate();

	return (
		<Box
			sx={{
				display: "flex",
				flexDirection: "column",
				alignItems: "center",
				justifyContent: "center",
				minHeight: "50vh",
			}}
		>
			<Typography variant="h1" component="h1" gutterBottom>
				404
			</Typography>
			<Typography variant="h5" component="h2" gutterBottom>
				Page not found
			</Typography>
			<Typography variant="body1" color="text.secondary" paragraph>
				The page you are looking for doesn't exist.
			</Typography>
			<Button variant="contained" onClick={() => navigate("/")}>
				Go to Home
			</Button>
		</Box>
	);
}
