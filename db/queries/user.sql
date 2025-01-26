-- name: SearchEntries :many
SELECT *
FROM entry
WHERE visibility = 'public'
ORDER BY published_at DESC
LIMIT 60 OFFSET ?;

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
