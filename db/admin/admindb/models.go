// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0

package admindb

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"
)

type EntryFormat string

const (
	EntryFormatHtml EntryFormat = "html"
	EntryFormatMkdn EntryFormat = "mkdn"
)

func (e *EntryFormat) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = EntryFormat(s)
	case string:
		*e = EntryFormat(s)
	default:
		return fmt.Errorf("unsupported scan type for EntryFormat: %T", src)
	}
	return nil
}

type NullEntryFormat struct {
	EntryFormat EntryFormat
	Valid       bool // Valid is true if EntryFormat is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullEntryFormat) Scan(value interface{}) error {
	if value == nil {
		ns.EntryFormat, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.EntryFormat.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullEntryFormat) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.EntryFormat), nil
}

type EntryVisibility string

const (
	EntryVisibilityPrivate EntryVisibility = "private"
	EntryVisibilityPublic  EntryVisibility = "public"
)

func (e *EntryVisibility) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = EntryVisibility(s)
	case string:
		*e = EntryVisibility(s)
	default:
		return fmt.Errorf("unsupported scan type for EntryVisibility: %T", src)
	}
	return nil
}

type NullEntryVisibility struct {
	EntryVisibility EntryVisibility
	Valid           bool // Valid is true if EntryVisibility is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullEntryVisibility) Scan(value interface{}) error {
	if value == nil {
		ns.EntryVisibility, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.EntryVisibility.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullEntryVisibility) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.EntryVisibility), nil
}

type AdminSession struct {
	SessionID      string
	Username       string
	ExpiresAt      time.Time
	CreatedAt      sql.NullTime
	LastAccessedAt sql.NullTime
}

type AmazonCache struct {
	Asin           string
	Title          sql.NullString
	ImageMediumUrl sql.NullString
	Link           string
	CreatedAt      sql.NullTime
}

type Entry struct {
	Path        string
	Title       string
	Body        string
	Visibility  EntryVisibility
	Format      EntryFormat
	PublishedAt sql.NullTime
	// last manualy edited at
	LastEditedAt sql.NullTime
	CreatedAt    sql.NullTime
	UpdatedAt    sql.NullTime
}

type EntryImage struct {
	Path      string
	Url       sql.NullString
	CreatedAt sql.NullTime
}

type EntryLink struct {
	SrcPath  string
	DstTitle string
}
