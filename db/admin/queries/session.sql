-- name: CreateSession :exec
INSERT INTO admin_session (session_id, username, expires_at)
VALUES (?, ?, ?);

-- name: GetSession :one
SELECT session_id, username, expires_at, created_at, last_accessed_at
FROM admin_session
WHERE session_id = ? AND expires_at > NOW()
LIMIT 1;

-- name: UpdateSessionLastAccessed :exec
UPDATE admin_session
SET last_accessed_at = NOW()
WHERE session_id = ?;

-- name: DeleteSession :exec
DELETE FROM admin_session
WHERE session_id = ?;

-- name: DeleteExpiredSessions :exec
DELETE FROM admin_session
WHERE expires_at < NOW();