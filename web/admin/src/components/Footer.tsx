import { Box, Link, Typography } from "@mui/material";
import { useEffect, useState } from "react";
import { getBuildInfo } from "../generated-client";

interface BuildInfo {
	buildTime: string;
	gitCommit: string;
	gitShortCommit: string;
	gitBranch: string;
	gitTag?: string;
	githubUrl: string;
}

export default function Footer() {
	const [buildInfo, setBuildInfo] = useState<BuildInfo | null>(null);

	useEffect(() => {
		getBuildInfo()
			.then((response) => {
				if (response.status === 200 && response.data) {
					setBuildInfo(response.data);
				}
			})
			.catch((error) => {
				console.error("Failed to fetch build info:", error);
			});
	}, []);

	if (!buildInfo) {
		return null;
	}

	return (
		<Box
			component="footer"
			sx={{
				py: 2,
				px: 2,
				mt: "auto",
				backgroundColor: (theme) =>
					theme.palette.mode === "light"
						? theme.palette.grey[200]
						: theme.palette.grey[800],
			}}
		>
			<Typography variant="body2" color="text.secondary" align="center">
				Build: {new Date(buildInfo.buildTime).toLocaleString()} |{" "}
				{buildInfo.gitBranch}
				{buildInfo.gitTag ? `@${buildInfo.gitTag}` : ""} |{" "}
				<Link href={buildInfo.githubUrl} target="_blank" rel="noopener">
					{buildInfo.gitShortCommit}
				</Link>
			</Typography>
		</Box>
	);
}
