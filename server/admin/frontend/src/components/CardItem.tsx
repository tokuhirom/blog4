import React from "react";

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
	const styles = {
		card: {
			display: "flex",
			flexDirection: "column" as const,
			alignItems: "flex-start",
			justifyContent: "flex-start",
			float: "left" as const,
			width: "150px",
			height: "150px",
			margin: "2px",
			padding: "4px",
			overflowY: "hidden" as const,
			borderRadius: "2px",
			textAlign: "left" as const,
			verticalAlign: "top" as const,
			overflowX: "hidden" as const,
			cursor: "pointer",
			backgroundColor,
			color,
		},
		title: {
			fontSize: "0.875rem",
			fontWeight: "bold",
			wordBreak: "break-all" as const,
		},
		content: {
			marginTop: "4px",
			color: "#4b5563",
			fontSize: "0.675rem",
		},
		img: {
			maxWidth: "100%",
			maxHeight: "100%",
			border: "1px solid #000",
		},
	};

	return (
		<button className="card" style={styles.card} onClick={onClick}>
			{title && <span style={styles.title}>{title}</span>}
			{imgSrc && <img src={imgSrc} alt="great one" style={styles.img} />}
			{title && <span style={styles.content}>{content}</span>}
		</button>
	);
}
