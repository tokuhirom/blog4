// Import styles that are needed for components
import "../src/app.css";

// Set up Material-UI for component testing
import { beforeMount } from "@playwright/experimental-ct-react/hooks";
import { ThemeProvider, createTheme } from "@mui/material/styles";
import CssBaseline from "@mui/material/CssBaseline";
import type { ReactNode } from "react";

const theme = createTheme();

beforeMount(async ({ App }) => {
	return (
		<ThemeProvider theme={theme}>
			<CssBaseline />
			<App />
		</ThemeProvider>
	);
});
