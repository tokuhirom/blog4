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
SELECT entry.*, entry_image.url AS image_url
FROM entry
    LEFT JOIN entry_image ON (entry.path = entry_image.path)
WHERE entry.path = ?;

-- name: UpdateEntryTitle :execrows
UPDATE entry
SET title = ?, last_edited_at = NOW()
WHERE path = ?;

-- name: UpdateEntryBody :execrows
UPDATE entry
SET body = ?, last_edited_at = NOW()
WHERE path = ?;

-- name: DeleteEntryLinkByPath :execrows
DELETE FROM entry_link WHERE src_path = ?;

-- name: InsertEntryLink :execrows
/* TODO batch insert */
INSERT INTO entry_link (src_path, dst_title)
VALUES (?, ?);

-- name: GetLinkedEntries :many
SELECT DISTINCT
    entry_link.dst_title AS dst_title,
    dest_entry.title title,
    dest_entry.path path,
    dest_entry.body,
    dest_entry.visibility,
    dest_entry.format,
    dest_entry.created_at,
    dest_entry.updated_at,
    entry_image.url AS image_url
FROM entry_link
    LEFT JOIN entry dest_entry ON (dest_entry.title = entry_link.dst_title)
    LEFT JOIN entry_image ON (dest_entry.path = entry_image.path)
WHERE entry_link.src_path = ?;

-- name: GetEntriesByLinkedTitle :many
SELECT DISTINCT entry.*
FROM entry_link
    INNER JOIN entry ON (entry.path = entry_link.src_path)
WHERE entry_link.dst_title = ?;

-- name: GetAllEntryTitles :many
SELECT title
FROM entry
ORDER BY title ASC;

-- name: CreateEmptyEntry :execrows
INSERT INTO entry
           (path, title, body, visibility)
    VALUES (?,        ?, '',    'private');

-- name: CreateEntryWithBody :execrows
INSERT INTO entry
           (path, title, body, visibility)
    VALUES (?,    ?,     ?,    'private');

-- name: DeleteEntry :execrows
DELETE FROM entry WHERE path = ?;
