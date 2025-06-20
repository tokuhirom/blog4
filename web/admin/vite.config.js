import react from "@vitejs/plugin-react";
import { defineConfig } from "vite";

// https://vite.dev/config/
export default defineConfig({
	base: "/admin/",
	plugins: [react()],
	server: {
		port: 6173,
	},
});
