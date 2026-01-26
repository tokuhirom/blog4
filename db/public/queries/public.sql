-- name: SearchEntries :many
SELECT entry.*, entry_image.url image_url
FROM entry
    LEFT JOIN entry_image ON (entry.path = entry_image.path)
WHERE visibility = 'public'
ORDER BY published_at DESC
LIMIT ? OFFSET ?;

-- name: GetEntryByPath :one
SELECT entry.*, entry_image.url image_url
FROM entry
    LEFT JOIN entry_image ON (entry.path = entry_image.path)
WHERE entry.path = ? AND visibility = 'public';

-- name: GetEntryByTitle :one
SELECT *
FROM entry
WHERE title = ? AND visibility = 'public';

-- name: GetAsin :one
SELECT *
FROM amazon_cache
WHERE asin = ?;

-- name: GetRelatedEntries1 :many
/* 現在表示しているエントリがリンクしているページ */
SELECT dst_entry.*
FROM entry dst_entry
         INNER JOIN entry_link ON (dst_entry.title = entry_link.dst_title)
WHERE entry_link.src_path = ? AND dst_entry.visibility = 'public';

-- name: GetRelatedEntries2 :many
/* 現在表示しているページにリンクしているページ */
SELECT src_entry.*
FROM entry src_entry
         INNER JOIN entry_link ON (src_entry.path = entry_link.src_path)
WHERE entry_link.dst_title = ? AND src_entry.visibility = 'public';

-- name: GetRelatedEntries3 :many
/* 現在表示しているページにリンクしているページのリンク先 */
SELECT dst_entry.*
FROM entry dst_entry
         INNER JOIN entry_link ON (dst_entry.title = entry_link.dst_title)
WHERE entry_link.src_path IN (SELECT src_entry.title
                              FROM entry src_entry
                                       INNER JOIN entry_link ON (src_entry.title = entry_link.src_path)
                              WHERE entry_link.dst_title = ?)
    AND dst_entry.visibility = 'public';

-- name: FullTextSearchEntries :many
SELECT entry.*, entry_image.url image_url,
       MATCH(entry.title, entry.body) AGAINST(? IN NATURAL LANGUAGE MODE) AS relevance
FROM entry
    LEFT JOIN entry_image ON (entry.path = entry_image.path)
WHERE visibility = 'public'
  AND MATCH(entry.title, entry.body) AGAINST(? IN NATURAL LANGUAGE MODE)
ORDER BY relevance DESC
LIMIT ? OFFSET ?;
