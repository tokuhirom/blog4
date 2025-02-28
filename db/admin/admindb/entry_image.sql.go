// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: entry_image.sql

package admindb

import (
	"context"
	"database/sql"
)

const deleteEntryImageByPath = `-- name: DeleteEntryImageByPath :execrows
DELETE FROM entry_image WHERE path = ?
`

func (q *Queries) DeleteEntryImageByPath(ctx context.Context, path string) (int64, error) {
	result, err := q.db.ExecContext(ctx, deleteEntryImageByPath, path)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

const getEntryImageNotProcessedEntries = `-- name: GetEntryImageNotProcessedEntries :many
SELECT entry.path, entry.title, entry.body, entry.visibility, entry.format, entry.published_at, entry.last_edited_at, entry.created_at, entry.updated_at
FROM entry
    LEFT JOIN entry_image ON (entry.path = entry_image.path)
WHERE entry_image.path IS NULL
    AND (body LIKE '%[asin:%' OR body LIKE '%.png%' OR body LIKE '%.jpg%')
ORDER BY updated_at DESC
`

func (q *Queries) GetEntryImageNotProcessedEntries(ctx context.Context) ([]Entry, error) {
	rows, err := q.db.QueryContext(ctx, getEntryImageNotProcessedEntries)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Entry
	for rows.Next() {
		var i Entry
		if err := rows.Scan(
			&i.Path,
			&i.Title,
			&i.Body,
			&i.Visibility,
			&i.Format,
			&i.PublishedAt,
			&i.LastEditedAt,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const insertEntryImage = `-- name: InsertEntryImage :execrows
INSERT INTO entry_image (path, url)
VALUES (?, ?)
`

type InsertEntryImageParams struct {
	Path string
	Url  sql.NullString
}

func (q *Queries) InsertEntryImage(ctx context.Context, arg InsertEntryImageParams) (int64, error) {
	result, err := q.db.ExecContext(ctx, insertEntryImage, arg.Path, arg.Url)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
