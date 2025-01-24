-- name: SearchEntries :many
SELECT *
FROM entry
ORDER BY published_at DESC
LIMIT 60 OFFSET ?;

-- name: GetEntryByPath :one
SELECT *
FROM entry
WHERE path = ? AND visibility = 'public';
