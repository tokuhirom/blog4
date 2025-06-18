package admin

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strings"

	"github.com/tokuhirom/blog4/db/admin/admindb"
	"github.com/tokuhirom/blog4/internal/admin/openapi"
)

func getLinkPalletData(ctx context.Context, db *sql.DB, queries *admindb.Queries, targetPath string, targetTitle string) (*openapi.LinkPalletData, error) {
	// このエントリがリンクしているページのリストを取得
	links, err := queries.GetLinkedEntries(ctx, targetPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get linked entries for path %s: %w", targetPath, err)
	}

	// このエントリにリンクしているページのリストを取得
	reverseLinks, err := queries.GetEntriesByLinkedTitle(ctx, targetTitle)
	if err != nil {
		return nil, fmt.Errorf("failed to get entries by linked title %s: %w", targetTitle, err)
	}

	// links の指す先のタイトルにリンクしているエントリのリストを取得
	twohopEntries, err := getEntriesByLinkedTitles(
		ctx,
		db,
		targetPath,
		links,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get entries by linked titles: %w", err)
	}

	return buildLinkPalletData(links, reverseLinks, twohopEntries, targetPath), nil
}

type GetEntriesByLinkedTitlesRow struct {
	DstTitle   string
	Path       string
	Title      string
	Body       string
	Visibility admindb.EntryVisibility
	Format     admindb.EntryFormat
	CreatedAt  sql.NullTime
	UpdatedAt  sql.NullTime
	ImageUrl   sql.NullString
}

func getEntriesByLinkedTitles(ctx context.Context, db *sql.DB, targetPath string, links []admindb.GetLinkedEntriesRow) ([]*GetEntriesByLinkedTitlesRow, error) {
	if len(links) == 0 {
		return []*GetEntriesByLinkedTitlesRow{}, nil
	}

	var linkedTitles []string
	for _, link := range links {
		linkedTitles = append(linkedTitles, strings.ToLower(link.DstTitle))
	}

	// プレースホルダーを動的に生成
	placeholders := make([]string, len(linkedTitles))
	args := make([]interface{}, 0, len(linkedTitles)+1)
	args = append(args, targetPath)

	for i, title := range linkedTitles {
		placeholders[i] = "LOWER(?)"
		args = append(args, title)
	}

	query := fmt.Sprintf(`
		SELECT DISTINCT
			entry_link.dst_title,
			entry.path,
			entry.title,
			entry.body,
			entry.visibility,
			entry.format,
			entry.created_at,
			entry.updated_at,
			entry_image.url AS image_url
		FROM entry_link
			INNER JOIN entry ON (entry.path = entry_link.src_path)
			LEFT JOIN entry_image ON (entry.path = entry_image.path)
		WHERE entry.path != ?
			AND LOWER(entry_link.dst_title) IN (%s)
	`, strings.Join(placeholders, ", "))

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			slog.Error("Cannot close rows", slog.Any("error", err))
		}
	}(rows)

	var results []*GetEntriesByLinkedTitlesRow
	for rows.Next() {
		var row GetEntriesByLinkedTitlesRow
		if err := rows.Scan(
			&row.DstTitle,
			&row.Path,
			&row.Title,
			&row.Body,
			&row.Visibility,
			&row.Format,
			&row.CreatedAt,
			&row.UpdatedAt,
			&row.ImageUrl,
		); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		results = append(results, &row)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return results, nil
}

func buildLinkPalletData(links []admindb.GetLinkedEntriesRow, reverseLinks []admindb.Entry, twohopEntries []*GetEntriesByLinkedTitlesRow, targetPath string) *openapi.LinkPalletData {
	// twohopEntries を dst_title でグループ化
	twohopEntriesByTitle := make(map[string][]openapi.EntryWithImage)
	for _, entry := range twohopEntries {
		lowerDstTitle := strings.ToLower(entry.DstTitle)
		twohopEntriesByTitle[lowerDstTitle] = append(twohopEntriesByTitle[lowerDstTitle], openapi.EntryWithImage{
			Path:       entry.Path,
			Title:      entry.Title,
			Body:       entry.Body,
			Visibility: string(entry.Visibility),
			Format:     string(entry.Format),
			ImageUrl:   toOptNilString(entry.ImageUrl),
		})
	}

	var resultLinks []*openapi.EntryWithImage
	var resultTwoHops []openapi.TwoHopLink
	var newLinks []string
	seenPath := make(map[string]struct{})
	seenPath[targetPath] = struct{}{}

	// twohopEntries に入っているエントリのリストを作成
	for _, link := range links {
		if link.Path.Valid {
			seenPath[link.Path.String] = struct{}{}
		}

		lowerDstTitle := strings.ToLower(link.DstTitle)
		if entries, exists := twohopEntriesByTitle[lowerDstTitle]; exists {
			resultTwoHops = append(resultTwoHops, openapi.TwoHopLink{
				Src: openapi.EntryWithDestTitle{
					Path:       link.Path.String,
					Title:      link.Title.String,
					Body:       link.Body.String,
					Visibility: string(link.Visibility.EntryVisibility),
					Format:     string(link.Format.EntryFormat),
					ImageUrl:   toOptNilString(link.ImageUrl),
				},
				Links: entries,
			})
			for _, entry := range entries {
				seenPath[entry.Path] = struct{}{}
			}
		} else {
			if link.Body.Valid && link.Body.String != "" {
				resultLinks = append(resultLinks, &openapi.EntryWithImage{
					Path:       link.Path.String,
					Title:      link.Title.String,
					Body:       link.Body.String,
					Visibility: string(link.Visibility.EntryVisibility),
					Format:     string(link.Format.EntryFormat),
					ImageUrl:   toOptNilString(link.ImageUrl),
				})
			} else {
				newLinks = append(newLinks, link.DstTitle)
			}
		}
	}

	for _, reverseLink := range reverseLinks {
		if _, exists := seenPath[reverseLink.Path]; !exists {
			resultLinks = append(resultLinks, &openapi.EntryWithImage{
				Path:       reverseLink.Path,
				Title:      reverseLink.Title,
				Body:       reverseLink.Body,
				Visibility: string(reverseLink.Visibility),
				Format:     string(reverseLink.Format),
				ImageUrl:   openapi.OptNilString{Null: true},
			})
		}
	}

	return &openapi.LinkPalletData{
		NewLinks: uniqueStrings(newLinks),
		Links:    uniqueEntries(resultLinks),
		Twohops:  resultTwoHops,
	}
}

func toOptNilString(src sql.NullString) openapi.OptNilString {
	if src.Valid {
		return openapi.OptNilString{Null: false, Value: src.String}
	}
	return openapi.OptNilString{Null: true}
}

// uniqueStrings returns a slice of unique strings.
func uniqueStrings(input []string) []string {
	uniqueMap := make(map[string]struct{})
	var result []string
	for _, str := range input {
		if _, exists := uniqueMap[str]; !exists {
			uniqueMap[str] = struct{}{}
			result = append(result, str)
		}
	}
	return result
}

// uniqueEntries returns a slice of unique EntryWithImage based on their Path.
func uniqueEntries(input []*openapi.EntryWithImage) []openapi.EntryWithImage {
	uniqueMap := make(map[string]*openapi.EntryWithImage)
	for _, entry := range input {
		if _, exists := uniqueMap[entry.Path]; !exists {
			uniqueMap[entry.Path] = entry
		}
	}
	var result []openapi.EntryWithImage
	for _, entry := range uniqueMap {
		result = append(result, *entry)
	}
	return result
}
