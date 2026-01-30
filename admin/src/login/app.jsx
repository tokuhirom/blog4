import { useCallback, useState } from 'preact/hooks';
import { login } from './api.js';

export function App() {
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [rememberMe, setRememberMe] = useState(false);
    const [error, setError] = useState('');
    const [submitting, setSubmitting] = useState(false);

    const handleSubmit = useCallback(
        async (e) => {
            e.preventDefault();
            if (submitting) return;
            setError('');
            setSubmitting(true);
            try {
                const data = await login(username, password, rememberMe);
                if (data.ok && data.redirect) {
                    window.location.href = data.redirect;
                    return;
                }
                setError(data.error || 'Login failed');
            } catch {
                setError('Network error. Please try again.');
            }
            setSubmitting(false);
        },
        [username, password, rememberMe, submitting],
    );

    return (
        <div class="login-box">
            <h1>Admin Login</h1>

            {error && (
                <div class="error-message">
                    <p>{error}</p>
                </div>
            )}

            <form onSubmit={handleSubmit}>
                <div class="form-group">
                    <label for="username">Username</label>
                    <input
                        type="text"
                        id="username"
                        name="username"
                        required
                        autoComplete="username"
                        value={username}
                        onInput={(e) => setUsername(e.currentTarget.value)}
                    />
                </div>

                <div class="form-group">
                    <label for="password">Password</label>
                    <input
                        type="password"
                        id="password"
                        name="password"
                        required
                        autoComplete="current-password"
                        value={password}
                        onInput={(e) => setPassword(e.currentTarget.value)}
                    />
                </div>

                <div class="form-group checkbox-group">
                    <label>
                        <input
                            type="checkbox"
                            checked={rememberMe}
                            onChange={(e) => setRememberMe(e.currentTarget.checked)}
                        />
                        <span>Remember me for 30 days</span>
                    </label>
                </div>

                <button type="submit" class="btn btn-primary btn-login" disabled={submitting}>
                    {submitting ? 'Logging in...' : 'Login'}
                </button>
            </form>
        </div>
    );
}
