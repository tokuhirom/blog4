-- name: GetEntryVisibility :one
SELECT visibility, published_at
FROM entry
WHERE path = ? FOR UPDATE;

-- name: UpdateVisibility :exec
UPDATE entry
SET visibility = ?
WHERE path = ?;

-- name: UpdatePublishedAt :exec
UPDATE entry
SET published_at = NOW()
WHERE path = ? AND published_at IS NULL;
