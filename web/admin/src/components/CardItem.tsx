import type React from "react";
import { Card, CardContent, Typography, Box } from "@mui/material";

interface CardItemProps {
	title?: string;
	content?: string;
	onClick?: (event: React.MouseEvent) => void;
	backgroundColor?: string;
	color?: string;
	imgSrc?: string | null;
}

export default function CardItem({
	title,
	content = "",
	onClick = () => {},
	backgroundColor = "#f6f6f6",
	color = "#0f0f0f",
	imgSrc,
}: CardItemProps) {
	return (
		<Card
			sx={{
				width: 170,
				height: 170,
				display: "flex",
				flexDirection: "column",
				backgroundColor,
				color,
				cursor: "pointer",
				"&:hover": {
					boxShadow: 3,
				},
			}}
			onClick={onClick}
		>
			<CardContent
				sx={{
					height: "100%",
					display: "flex",
					flexDirection: "column",
					p: 1.5,
				}}
			>
				{title && (
					<Typography
						variant="body2"
						component="h2"
						sx={{
							fontWeight: 600,
							fontSize: "0.875rem",
							overflow: "hidden",
							textOverflow: "ellipsis",
							display: "-webkit-box",
							WebkitLineClamp: 2,
							WebkitBoxOrient: "vertical",
							mb: 0.5,
						}}
					>
						{title}
					</Typography>
				)}
				{imgSrc && (
					<Box
						component="img"
						src={imgSrc}
						alt={title || "Entry image"}
						sx={{
							width: "100%",
							height: 60,
							objectFit: "cover",
							borderRadius: 1,
							mb: 0.5,
						}}
					/>
				)}
				{content && (
					<Typography
						variant="caption"
						color="text.secondary"
						sx={{
							fontSize: "0.75rem",
							overflow: "hidden",
							textOverflow: "ellipsis",
							display: "-webkit-box",
							WebkitLineClamp: imgSrc ? 3 : 6,
							WebkitBoxOrient: "vertical",
							flexGrow: 1,
						}}
					>
						{content}
					</Typography>
				)}
			</CardContent>
		</Card>
	);
}
