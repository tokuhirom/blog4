import { Outlet } from "react-router-dom";
import AdminHeader from "./components/AdminHeader";

export default function Layout() {
	return (
		<div>
			<AdminHeader />
			<main style={{ marginTop: "4rem", marginBottom: "4rem" }}>
				<Outlet />
			</main>
		</div>
	);
}
