-- name: GetLatestEntries :many
SELECT entry.*, entry_image.url AS image_url
FROM entry
    LEFT JOIN entry_image ON (entry.path = entry_image.path)
WHERE (? IS NULL OR last_edited_at <= ?)
ORDER BY
    last_edited_at DESC
    , path DESC
LIMIT ?;

-- name: AdminGetEntryByPath :one
SELECT * FROM entry WHERE path = ?;

-- name: GetVisibility :one
SELECT visibility FROM entry WHERE path = ?;

-- name: UpdateVisibility :execrows
UPDATE entry SET visibility = ? WHERE path = ?;

-- name: UpdatePublishedAt :execrows
UPDATE entry SET published_at = ? WHERE path = ?;

-- name: UpdateEntryTitle :execrows
UPDATE entry
SET title = ?, last_edited_at = NOW()
WHERE path = ?;

-- name: UpdateEntryBody :execrows
UPDATE entry
SET body = ?, last_edited_at = NOW()
WHERE path = ?;

-- name: DeleteEntryLink :execrows
DELETE FROM entry_link WHERE src_path = ?;

-- name: InsertEntryLink :execrows
/* TODO batch insert */
INSERT INTO entry_link (src_path, dst_title)
VALUES (?, ?);



