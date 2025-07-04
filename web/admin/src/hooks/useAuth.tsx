import type React from "react";
import {
	createContext,
	useCallback,
	useContext,
	useEffect,
	useState,
} from "react";
import { useLocation, useNavigate } from "react-router-dom";
import { authCheck, authLogout } from "../generated-client";

interface AuthContextType {
	isAuthenticated: boolean;
	username: string | null;
	loading: boolean;
	logout: () => Promise<void>;
	checkAuth: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | null>(null);

export function AuthProvider({ children }: { children: React.ReactNode }) {
	const [isAuthenticated, setIsAuthenticated] = useState(false);
	const [username, setUsername] = useState<string | null>(null);
	const [loading, setLoading] = useState(true);
	const navigate = useNavigate();
	const location = useLocation();

	const checkAuth = useCallback(async () => {
		try {
			setLoading(true);
			const response = await authCheck();
			if (response.status === 200 && response.data.authenticated) {
				setIsAuthenticated(true);
				setUsername(response.data.username || null);
			} else {
				setIsAuthenticated(false);
				setUsername(null);
			}
		} catch (_error) {
			// Session expired or network error
			setIsAuthenticated(false);
			setUsername(null);
		} finally {
			setLoading(false);
		}
	}, []);

	const logout = useCallback(async () => {
		try {
			await authLogout();
			// Clear auth state immediately
			setIsAuthenticated(false);
			setUsername(null);
			// Force navigation to login page
			navigate("/login", { replace: true });
		} catch (error) {
			console.error("Logout error:", error);
			// Even if logout fails, clear client state
			setIsAuthenticated(false);
			setUsername(null);
			navigate("/login", { replace: true });
		}
	}, [navigate]);

	useEffect(() => {
		checkAuth();
	}, [checkAuth]);

	useEffect(() => {
		if (!loading && !isAuthenticated && location.pathname !== "/login") {
			navigate("/login");
		}
	}, [loading, isAuthenticated, location.pathname, navigate]);

	return (
		<AuthContext.Provider
			value={{ isAuthenticated, username, loading, logout, checkAuth }}
		>
			{children}
		</AuthContext.Provider>
	);
}

export function useAuth() {
	const context = useContext(AuthContext);
	if (!context) {
		throw new Error("useAuth must be used within an AuthProvider");
	}
	return context;
}
