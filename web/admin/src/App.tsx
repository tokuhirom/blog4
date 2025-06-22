import { BrowserRouter, Route, Routes } from "react-router-dom";
import { ThemeProvider, createTheme } from "@mui/material/styles";
import CssBaseline from "@mui/material/CssBaseline";
import Layout from "./Layout";
import NotFound from "./NotFound";
import AdminEntryPage from "./pages/AdminEntryPage";
import TopPage from "./pages/TopPage";
import LoginPage from "./pages/LoginPage";
import { AuthProvider } from "./hooks/useAuth";
import ProtectedRoute from "./components/ProtectedRoute";

const theme = createTheme({
	palette: {
		mode: "light",
		primary: {
			main: "#dc004e", // Red as primary
		},
		secondary: {
			main: "#1976d2",
		},
	},
});

export default function App() {
	return (
		<ThemeProvider theme={theme}>
			<CssBaseline />
			<BrowserRouter basename="/admin">
				<AuthProvider>
					<Routes>
						<Route path="/login" element={<LoginPage />} />
						<Route
							path="/"
							element={
								<ProtectedRoute>
									<Layout />
								</ProtectedRoute>
							}
						>
							<Route index element={<TopPage />} />
							<Route path="/entry/*" element={<AdminEntryPage />} />
							<Route path="*" element={<NotFound />} />
						</Route>
					</Routes>
				</AuthProvider>
			</BrowserRouter>
		</ThemeProvider>
	);
}
