// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: admin.sql

package admindb

import (
	"context"
	"database/sql"
)

const adminGetEntryByPath = `-- name: AdminGetEntryByPath :one
SELECT entry.path, entry.title, entry.body, entry.visibility, entry.format, entry.published_at, entry.last_edited_at, entry.created_at, entry.updated_at, entry_image.url AS image_url
FROM entry
    LEFT JOIN entry_image ON (entry.path = entry_image.path)
WHERE entry.path = ?
`

type AdminGetEntryByPathRow struct {
	Path         string
	Title        string
	Body         string
	Visibility   EntryVisibility
	Format       EntryFormat
	PublishedAt  sql.NullTime
	LastEditedAt sql.NullTime
	CreatedAt    sql.NullTime
	UpdatedAt    sql.NullTime
	ImageUrl     sql.NullString
}

func (q *Queries) AdminGetEntryByPath(ctx context.Context, path string) (AdminGetEntryByPathRow, error) {
	row := q.db.QueryRowContext(ctx, adminGetEntryByPath, path)
	var i AdminGetEntryByPathRow
	err := row.Scan(
		&i.Path,
		&i.Title,
		&i.Body,
		&i.Visibility,
		&i.Format,
		&i.PublishedAt,
		&i.LastEditedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.ImageUrl,
	)
	return i, err
}

const createEmptyEntry = `-- name: CreateEmptyEntry :execrows
INSERT INTO entry
           (path, title, body, visibility)
    VALUES (?,        ?, '',    'private')
`

type CreateEmptyEntryParams struct {
	Path  string
	Title string
}

func (q *Queries) CreateEmptyEntry(ctx context.Context, arg CreateEmptyEntryParams) (int64, error) {
	result, err := q.db.ExecContext(ctx, createEmptyEntry, arg.Path, arg.Title)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

const deleteEntry = `-- name: DeleteEntry :execrows
DELETE FROM entry WHERE path = ?
`

func (q *Queries) DeleteEntry(ctx context.Context, path string) (int64, error) {
	result, err := q.db.ExecContext(ctx, deleteEntry, path)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

const deleteEntryLink = `-- name: DeleteEntryLink :execrows
DELETE FROM entry_link WHERE src_path = ?
`

func (q *Queries) DeleteEntryLink(ctx context.Context, srcPath string) (int64, error) {
	result, err := q.db.ExecContext(ctx, deleteEntryLink, srcPath)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

const getAllEntryTitles = `-- name: GetAllEntryTitles :many
SELECT title
FROM entry
ORDER BY title ASC
`

func (q *Queries) GetAllEntryTitles(ctx context.Context) ([]string, error) {
	rows, err := q.db.QueryContext(ctx, getAllEntryTitles)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var title string
		if err := rows.Scan(&title); err != nil {
			return nil, err
		}
		items = append(items, title)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getEntriesByLinkedTitle = `-- name: GetEntriesByLinkedTitle :many
SELECT DISTINCT entry.path, entry.title, entry.body, entry.visibility, entry.format, entry.published_at, entry.last_edited_at, entry.created_at, entry.updated_at
FROM entry_link
    INNER JOIN entry ON (entry.path = entry_link.src_path)
WHERE entry_link.dst_title = ?
`

func (q *Queries) GetEntriesByLinkedTitle(ctx context.Context, dstTitle string) ([]Entry, error) {
	rows, err := q.db.QueryContext(ctx, getEntriesByLinkedTitle, dstTitle)
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

const getLatestEntries = `-- name: GetLatestEntries :many
SELECT entry.path, entry.title, entry.body, entry.visibility, entry.format, entry.published_at, entry.last_edited_at, entry.created_at, entry.updated_at, entry_image.url AS image_url
FROM entry
    LEFT JOIN entry_image ON (entry.path = entry_image.path)
WHERE (? IS NULL OR last_edited_at <= ?)
ORDER BY
    last_edited_at DESC
    , path DESC
LIMIT ?
`

type GetLatestEntriesParams struct {
	Column1      interface{}
	LastEditedAt sql.NullTime
	Limit        int32
}

type GetLatestEntriesRow struct {
	Path         string
	Title        string
	Body         string
	Visibility   EntryVisibility
	Format       EntryFormat
	PublishedAt  sql.NullTime
	LastEditedAt sql.NullTime
	CreatedAt    sql.NullTime
	UpdatedAt    sql.NullTime
	ImageUrl     sql.NullString
}

func (q *Queries) GetLatestEntries(ctx context.Context, arg GetLatestEntriesParams) ([]GetLatestEntriesRow, error) {
	rows, err := q.db.QueryContext(ctx, getLatestEntries, arg.Column1, arg.LastEditedAt, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetLatestEntriesRow
	for rows.Next() {
		var i GetLatestEntriesRow
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
			&i.ImageUrl,
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

const getLinkedEntries = `-- name: GetLinkedEntries :many
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
WHERE entry_link.src_path = ?
`

type GetLinkedEntriesRow struct {
	DstTitle   string
	Title      sql.NullString
	Path       sql.NullString
	Body       sql.NullString
	Visibility NullEntryVisibility
	Format     NullEntryFormat
	CreatedAt  sql.NullTime
	UpdatedAt  sql.NullTime
	ImageUrl   sql.NullString
}

func (q *Queries) GetLinkedEntries(ctx context.Context, srcPath string) ([]GetLinkedEntriesRow, error) {
	rows, err := q.db.QueryContext(ctx, getLinkedEntries, srcPath)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetLinkedEntriesRow
	for rows.Next() {
		var i GetLinkedEntriesRow
		if err := rows.Scan(
			&i.DstTitle,
			&i.Title,
			&i.Path,
			&i.Body,
			&i.Visibility,
			&i.Format,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.ImageUrl,
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

const getVisibility = `-- name: GetVisibility :one
SELECT visibility FROM entry WHERE path = ?
`

func (q *Queries) GetVisibility(ctx context.Context, path string) (EntryVisibility, error) {
	row := q.db.QueryRowContext(ctx, getVisibility, path)
	var visibility EntryVisibility
	err := row.Scan(&visibility)
	return visibility, err
}

const insertEntryLink = `-- name: InsertEntryLink :execrows
INSERT INTO entry_link (src_path, dst_title)
VALUES (?, ?)
`

type InsertEntryLinkParams struct {
	SrcPath  string
	DstTitle string
}

// TODO batch insert
func (q *Queries) InsertEntryLink(ctx context.Context, arg InsertEntryLinkParams) (int64, error) {
	result, err := q.db.ExecContext(ctx, insertEntryLink, arg.SrcPath, arg.DstTitle)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

const updateEntryBody = `-- name: UpdateEntryBody :execrows
UPDATE entry
SET body = ?, last_edited_at = NOW()
WHERE path = ?
`

type UpdateEntryBodyParams struct {
	Body string
	Path string
}

func (q *Queries) UpdateEntryBody(ctx context.Context, arg UpdateEntryBodyParams) (int64, error) {
	result, err := q.db.ExecContext(ctx, updateEntryBody, arg.Body, arg.Path)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

const updateEntryTitle = `-- name: UpdateEntryTitle :execrows
UPDATE entry
SET title = ?, last_edited_at = NOW()
WHERE path = ?
`

type UpdateEntryTitleParams struct {
	Title string
	Path  string
}

func (q *Queries) UpdateEntryTitle(ctx context.Context, arg UpdateEntryTitleParams) (int64, error) {
	result, err := q.db.ExecContext(ctx, updateEntryTitle, arg.Title, arg.Path)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

const updatePublishedAt = `-- name: UpdatePublishedAt :execrows
UPDATE entry SET published_at = ? WHERE path = ?
`

type UpdatePublishedAtParams struct {
	PublishedAt sql.NullTime
	Path        string
}

func (q *Queries) UpdatePublishedAt(ctx context.Context, arg UpdatePublishedAtParams) (int64, error) {
	result, err := q.db.ExecContext(ctx, updatePublishedAt, arg.PublishedAt, arg.Path)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

const updateVisibility = `-- name: UpdateVisibility :execrows
UPDATE entry SET visibility = ? WHERE path = ?
`

type UpdateVisibilityParams struct {
	Visibility EntryVisibility
	Path       string
}

func (q *Queries) UpdateVisibility(ctx context.Context, arg UpdateVisibilityParams) (int64, error) {
	result, err := q.db.ExecContext(ctx, updateVisibility, arg.Visibility, arg.Path)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
