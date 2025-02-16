// Code generated by ogen, DO NOT EDIT.

package openapi

import (
	"context"
)

// Handler handles operations described by OpenAPI v3 specification.
type Handler interface {
	// CreateEntry implements createEntry operation.
	//
	// Create a new entry.
	//
	// POST /entries
	CreateEntry(ctx context.Context, req *CreateEntryRequest) (CreateEntryRes, error)
	// DeleteEntry implements deleteEntry operation.
	//
	// Delete an entry.
	//
	// DELETE /entries/{path}
	DeleteEntry(ctx context.Context, params DeleteEntryParams) (DeleteEntryRes, error)
	// GetAllEntryTitles implements getAllEntryTitles operation.
	//
	// Get all entry titles.
	//
	// GET /entries/titles
	GetAllEntryTitles(ctx context.Context) (EntryTitlesResponse, error)
	// GetEntryByDynamicPath implements getEntryByDynamicPath operation.
	//
	// Get entry by dynamic path.
	//
	// GET /entries/{path}
	GetEntryByDynamicPath(ctx context.Context, params GetEntryByDynamicPathParams) (*GetLatestEntriesRow, error)
	// GetLatestEntries implements getLatestEntries operation.
	//
	// Get latest entries.
	//
	// GET /entries
	GetLatestEntries(ctx context.Context, params GetLatestEntriesParams) ([]GetLatestEntriesRow, error)
	// GetLinkPallet implements getLinkPallet operation.
	//
	// Get linked entry paths.
	//
	// GET /entries/{path}/link-pallet
	GetLinkPallet(ctx context.Context, params GetLinkPalletParams) (*LinkPalletData, error)
	// GetLinkedEntryPaths implements getLinkedEntryPaths operation.
	//
	// Get linked entry paths.
	//
	// GET /entries/{path}/linked-paths
	GetLinkedEntryPaths(ctx context.Context, params GetLinkedEntryPathsParams) (GetLinkedEntryPathsRes, error)
	// UpdateEntryBody implements updateEntryBody operation.
	//
	// Update entry body.
	//
	// PUT /entries/{path}/body
	UpdateEntryBody(ctx context.Context, req *UpdateEntryBodyRequest, params UpdateEntryBodyParams) (UpdateEntryBodyRes, error)
	// UpdateEntryTitle implements updateEntryTitle operation.
	//
	// Update entry title.
	//
	// PUT /entries/{path}/title
	UpdateEntryTitle(ctx context.Context, req *UpdateEntryTitleRequest, params UpdateEntryTitleParams) (UpdateEntryTitleRes, error)
	// NewError creates *ErrorResponseStatusCode from error returned by handler.
	//
	// Used for common default response.
	NewError(ctx context.Context, err error) *ErrorResponseStatusCode
}

// Server implements http server based on OpenAPI v3 specification and
// calls Handler to handle requests.
type Server struct {
	h Handler
	baseServer
}

// NewServer creates new Server.
func NewServer(h Handler, opts ...ServerOption) (*Server, error) {
	s, err := newServerConfig(opts...).baseServer()
	if err != nil {
		return nil, err
	}
	return &Server{
		h:          h,
		baseServer: s,
	}, nil
}
