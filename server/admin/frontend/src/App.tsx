import { BrowserRouter, Route, Routes } from "react-router-dom";
import Layout from "./Layout";
import NotFound from "./NotFound";
import AdminEntryPage from "./pages/AdminEntryPage";
import TopPage from "./pages/TopPage";

export default function App() {
	return (
		<BrowserRouter>
			<Routes>
				<Route path="/" element={<Layout />}>
					<Route path="/admin/" element={<TopPage />} />
					<Route path="/admin/entry/*" element={<AdminEntryPage />} />
					<Route path="*" element={<NotFound />} />
				</Route>
			</Routes>
		</BrowserRouter>
	);
}
