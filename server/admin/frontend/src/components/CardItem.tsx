import type React from "react";
import styles from "./CardItem.module.css";

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
	const cardStyle = {
		backgroundColor,
		color,
	};

	return (
		<button
			type="button"
			className={styles.card}
			style={cardStyle}
			onClick={onClick}
		>
			{title && <span className={styles.title}>{title}</span>}
			{imgSrc && <img src={imgSrc} alt="great one" className={styles.img} />}
			{title && <span className={styles.content}>{content}</span>}
		</button>
	);
}
