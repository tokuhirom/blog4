// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: session.sql

package admindb

import (
	"context"
	"time"
)

const createSession = `-- name: CreateSession :exec
INSERT INTO admin_session (session_id, username, expires_at)
VALUES (?, ?, ?)
`

type CreateSessionParams struct {
	SessionID string
	Username  string
	ExpiresAt time.Time
}

func (q *Queries) CreateSession(ctx context.Context, arg CreateSessionParams) error {
	_, err := q.db.ExecContext(ctx, createSession, arg.SessionID, arg.Username, arg.ExpiresAt)
	return err
}

const deleteExpiredSessions = `-- name: DeleteExpiredSessions :exec
DELETE FROM admin_session
WHERE expires_at < NOW()
`

func (q *Queries) DeleteExpiredSessions(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, deleteExpiredSessions)
	return err
}

const deleteSession = `-- name: DeleteSession :exec
DELETE FROM admin_session
WHERE session_id = ?
`

func (q *Queries) DeleteSession(ctx context.Context, sessionID string) error {
	_, err := q.db.ExecContext(ctx, deleteSession, sessionID)
	return err
}

const getSession = `-- name: GetSession :one
SELECT session_id, username, expires_at, created_at, last_accessed_at
FROM admin_session
WHERE session_id = ? AND expires_at > NOW()
LIMIT 1
`

func (q *Queries) GetSession(ctx context.Context, sessionID string) (AdminSession, error) {
	row := q.db.QueryRowContext(ctx, getSession, sessionID)
	var i AdminSession
	err := row.Scan(
		&i.SessionID,
		&i.Username,
		&i.ExpiresAt,
		&i.CreatedAt,
		&i.LastAccessedAt,
	)
	return i, err
}

const updateSessionLastAccessed = `-- name: UpdateSessionLastAccessed :exec
UPDATE admin_session
SET last_accessed_at = NOW()
WHERE session_id = ?
`

func (q *Queries) UpdateSessionLastAccessed(ctx context.Context, sessionID string) error {
	_, err := q.db.ExecContext(ctx, updateSessionLastAccessed, sessionID)
	return err
}
