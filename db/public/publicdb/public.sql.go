// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: public.sql

package publicdb

import (
	"context"
	"database/sql"
)

const getAsin = `-- name: GetAsin :one
SELECT asin, title, image_medium_url, link, created_at
FROM amazon_cache
WHERE asin = ?
`

func (q *Queries) GetAsin(ctx context.Context, asin string) (AmazonCache, error) {
	row := q.db.QueryRowContext(ctx, getAsin, asin)
	var i AmazonCache
	err := row.Scan(
		&i.Asin,
		&i.Title,
		&i.ImageMediumUrl,
		&i.Link,
		&i.CreatedAt,
	)
	return i, err
}

const getEntryByPath = `-- name: GetEntryByPath :one
SELECT path, title, body, visibility, format, published_at, last_edited_at, created_at, updated_at
FROM entry
WHERE path = ? AND visibility = 'public'
`

func (q *Queries) GetEntryByPath(ctx context.Context, path string) (Entry, error) {
	row := q.db.QueryRowContext(ctx, getEntryByPath, path)
	var i Entry
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
	)
	return i, err
}

const getEntryByTitle = `-- name: GetEntryByTitle :one
SELECT path, title, body, visibility, format, published_at, last_edited_at, created_at, updated_at
FROM entry
WHERE title = ? AND visibility = 'public'
`

func (q *Queries) GetEntryByTitle(ctx context.Context, title string) (Entry, error) {
	row := q.db.QueryRowContext(ctx, getEntryByTitle, title)
	var i Entry
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
	)
	return i, err
}

const getRelatedEntries1 = `-- name: GetRelatedEntries1 :many
SELECT dst_entry.path, dst_entry.title, dst_entry.body, dst_entry.visibility, dst_entry.format, dst_entry.published_at, dst_entry.last_edited_at, dst_entry.created_at, dst_entry.updated_at
FROM entry dst_entry
         INNER JOIN entry_link ON (dst_entry.title = entry_link.dst_title)
WHERE entry_link.src_path = ? AND dst_entry.visibility = 'public'
`

// 現在表示しているエントリがリンクしているページ
func (q *Queries) GetRelatedEntries1(ctx context.Context, srcPath string) ([]Entry, error) {
	rows, err := q.db.QueryContext(ctx, getRelatedEntries1, srcPath)
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

const getRelatedEntries2 = `-- name: GetRelatedEntries2 :many
SELECT src_entry.path, src_entry.title, src_entry.body, src_entry.visibility, src_entry.format, src_entry.published_at, src_entry.last_edited_at, src_entry.created_at, src_entry.updated_at
FROM entry src_entry
         INNER JOIN entry_link ON (src_entry.path = entry_link.src_path)
WHERE entry_link.dst_title = ? AND src_entry.visibility = 'public'
`

// 現在表示しているページにリンクしているページ
func (q *Queries) GetRelatedEntries2(ctx context.Context, dstTitle string) ([]Entry, error) {
	rows, err := q.db.QueryContext(ctx, getRelatedEntries2, dstTitle)
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

const getRelatedEntries3 = `-- name: GetRelatedEntries3 :many
SELECT dst_entry.path, dst_entry.title, dst_entry.body, dst_entry.visibility, dst_entry.format, dst_entry.published_at, dst_entry.last_edited_at, dst_entry.created_at, dst_entry.updated_at
FROM entry dst_entry
         INNER JOIN entry_link ON (dst_entry.title = entry_link.dst_title)
WHERE entry_link.src_path IN (SELECT src_entry.title
                              FROM entry src_entry
                                       INNER JOIN entry_link ON (src_entry.title = entry_link.src_path)
                              WHERE entry_link.dst_title = ?)
    AND dst_entry.visibility = 'public'
`

// 現在表示しているページにリンクしているページのリンク先
func (q *Queries) GetRelatedEntries3(ctx context.Context, dstTitle string) ([]Entry, error) {
	rows, err := q.db.QueryContext(ctx, getRelatedEntries3, dstTitle)
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

const searchEntries = `-- name: SearchEntries :many
SELECT entry.path, entry.title, entry.body, entry.visibility, entry.format, entry.published_at, entry.last_edited_at, entry.created_at, entry.updated_at, entry_image.url image_url
FROM entry
    LEFT JOIN entry_image ON (entry.path = entry_image.path)
WHERE visibility = 'public'
ORDER BY published_at DESC
LIMIT ? OFFSET ?
`

type SearchEntriesParams struct {
	Limit  int32
	Offset int32
}

type SearchEntriesRow struct {
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

func (q *Queries) SearchEntries(ctx context.Context, arg SearchEntriesParams) ([]SearchEntriesRow, error) {
	rows, err := q.db.QueryContext(ctx, searchEntries, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SearchEntriesRow
	for rows.Next() {
		var i SearchEntriesRow
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
