import { Outlet } from "react-router-dom";
import { Box, Container } from "@mui/material";
import AdminHeader from "./components/AdminHeader";
import Footer from "./components/Footer";

export default function Layout() {
	return (
		<Box sx={{ display: "flex", flexDirection: "column", minHeight: "100vh" }}>
			<AdminHeader />
			<Container component="main" sx={{ mt: 8, mb: 4, flex: 1 }} maxWidth="xl">
				<Outlet />
			</Container>
			<Footer />
		</Box>
	);
}
