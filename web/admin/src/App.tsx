import { BrowserRouter, Route, Routes } from "react-router-dom";
import { ThemeProvider, createTheme } from "@mui/material/styles";
import CssBaseline from "@mui/material/CssBaseline";
import Layout from "./Layout";
import NotFound from "./NotFound";
import AdminEntryPage from "./pages/AdminEntryPage";
import TopPage from "./pages/TopPage";

const theme = createTheme({
	palette: {
		mode: "light",
		primary: {
			main: "#1976d2",
		},
		secondary: {
			main: "#dc004e",
		},
	},
});

export default function App() {
	return (
		<ThemeProvider theme={theme}>
			<CssBaseline />
			<BrowserRouter>
				<Routes>
					<Route path="/" element={<Layout />}>
						<Route path="/admin/" element={<TopPage />} />
						<Route path="/admin/entry/*" element={<AdminEntryPage />} />
						<Route path="*" element={<NotFound />} />
					</Route>
				</Routes>
			</BrowserRouter>
		</ThemeProvider>
	);
}
