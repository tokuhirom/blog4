-- name: SearchEntries :many
SELECT *
FROM entry
WHERE visibility = 'public'
ORDER BY published_at DESC
LIMIT ? OFFSET ?;

-- name: GetEntryByPath :one
SELECT *
FROM entry
WHERE path = ? AND visibility = 'public';

-- name: GetEntryByTitle :one
SELECT *
FROM entry
WHERE title = ? AND visibility = 'public';

-- name: GetAsin :one
SELECT *
FROM amazon_cache
WHERE asin = ?;
