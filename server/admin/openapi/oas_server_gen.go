// Code generated by ogen, DO NOT EDIT.

package openapi

import (
	"context"
)

// Handler handles operations described by OpenAPI v3 specification.
type Handler interface {
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
	// GetLinkedEntryPaths implements getLinkedEntryPaths operation.
	//
	// Get linked entry paths.
	//
	// GET /entries/{path}/links
	GetLinkedEntryPaths(ctx context.Context, params GetLinkedEntryPathsParams) (LinkedEntriesResponse, error)
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
