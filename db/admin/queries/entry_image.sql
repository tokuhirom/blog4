-- name: GetEntryImageNotProcessedEntries :many
SELECT entry.*
FROM entry
    LEFT JOIN entry_image ON (entry.path = entry_image.path)
WHERE entry_image.path IS NULL
    AND (body LIKE '%[asin:%' OR body LIKE '%.png%' OR body LIKE '%.jpg%')
ORDER BY updated_at DESC;

-- name: InsertEntryImage :execrows
INSERT INTO entry_image (path, url)
VALUES (?, ?);

-- name: DeleteEntryImageByPath :execrows
DELETE FROM entry_image WHERE path = ?;

-- name: GetEntryImageByPath :one
SELECT path, url, created_at
FROM entry_image
WHERE path = ?;
