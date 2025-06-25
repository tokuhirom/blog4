import {
	Alert,
	Box,
	Button,
	Container,
	TextField,
	Typography,
} from "@mui/material";
import type React from "react";
import { useId, useState } from "react";
import { useNavigate } from "react-router-dom";
import { authLogin } from "../generated-client";
import { useAuth } from "../hooks/useAuth";

export default function LoginPage() {
	const [username, setUsername] = useState("");
	const [password, setPassword] = useState("");
	const [error, setError] = useState("");
	const [loading, setLoading] = useState(false);
	const navigate = useNavigate();
	const { checkAuth } = useAuth();
	const usernameId = useId();
	const passwordId = useId();

	const handleSubmit = async (e: React.FormEvent) => {
		e.preventDefault();
		setError("");
		setLoading(true);

		try {
			const response = await authLogin({ username, password });
			if (response.status === 200 && response.data.success) {
				// Update auth state before navigating
				await checkAuth();
				navigate("/");
			} else if (response.status === 401) {
				setError(response.data.message || "Invalid username or password");
			} else {
				setError("Login failed");
			}
		} catch (_err) {
			setError("Invalid username or password");
		} finally {
			setLoading(false);
		}
	};

	return (
		<Container component="main" maxWidth="xs">
			<Box
				sx={{
					marginTop: 8,
					display: "flex",
					flexDirection: "column",
					alignItems: "center",
				}}
			>
				<Typography component="h1" variant="h5">
					Admin Login
				</Typography>
				<Box component="form" onSubmit={handleSubmit} sx={{ mt: 1 }}>
					{error && (
						<Alert severity="error" sx={{ mb: 2 }}>
							{error}
						</Alert>
					)}
					<TextField
						margin="normal"
						required
						fullWidth
						id={usernameId}
						label="Username"
						name="username"
						autoComplete="username"
						autoFocus
						value={username}
						onChange={(e) => setUsername(e.target.value)}
					/>
					<TextField
						margin="normal"
						required
						fullWidth
						name="password"
						label="Password"
						type="password"
						id={passwordId}
						autoComplete="current-password"
						value={password}
						onChange={(e) => setPassword(e.target.value)}
					/>
					<Button
						type="submit"
						fullWidth
						variant="contained"
						sx={{ mt: 3, mb: 2 }}
						disabled={loading}
					>
						{loading ? "Logging in..." : "Login"}
					</Button>
				</Box>
			</Box>
		</Container>
	);
}
