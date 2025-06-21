import type React from "react";
import {
	Card,
	CardActionArea,
	CardContent,
	CardMedia,
	Typography,
} from "@mui/material";

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
				height: "100%",
				display: "flex",
				flexDirection: "column",
				backgroundColor,
				color,
			}}
		>
			<CardActionArea onClick={onClick} sx={{ height: "100%" }}>
				{imgSrc && (
					<CardMedia
						component="img"
						height="140"
						image={imgSrc}
						alt={title || "Entry image"}
						sx={{ objectFit: "cover" }}
					/>
				)}
				<CardContent>
					{title && (
						<Typography
							gutterBottom
							variant="h6"
							component="h2"
							sx={{
								overflow: "hidden",
								textOverflow: "ellipsis",
								display: "-webkit-box",
								WebkitLineClamp: 2,
								WebkitBoxOrient: "vertical",
							}}
						>
							{title}
						</Typography>
					)}
					{content && (
						<Typography
							variant="body2"
							color="text.secondary"
							sx={{
								overflow: "hidden",
								textOverflow: "ellipsis",
								display: "-webkit-box",
								WebkitLineClamp: 3,
								WebkitBoxOrient: "vertical",
							}}
						>
							{content}
						</Typography>
					)}
				</CardContent>
			</CardActionArea>
		</Card>
	);
}
