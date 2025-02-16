// Code generated by ogen, DO NOT EDIT.

package openapi

import (
	"fmt"
	"time"
)

func (s *ErrorResponseStatusCode) Error() string {
	return fmt.Sprintf("code %d: %+v", s.StatusCode, s.Response)
}

// Ref: #/components/schemas/CreateEntryRequest
type CreateEntryRequest struct {
	// The title of the new entry.
	Title OptString `json:"title"`
}

// GetTitle returns the value of Title.
func (s *CreateEntryRequest) GetTitle() OptString {
	return s.Title
}

// SetTitle sets the value of Title.
func (s *CreateEntryRequest) SetTitle(val OptString) {
	s.Title = val
}

// Ref: #/components/schemas/CreateEntryResponse
type CreateEntryResponse struct {
	// The path of the created entry.
	Path string `json:"path"`
}

// GetPath returns the value of Path.
func (s *CreateEntryResponse) GetPath() string {
	return s.Path
}

// SetPath sets the value of Path.
func (s *CreateEntryResponse) SetPath(val string) {
	s.Path = val
}

func (*CreateEntryResponse) createEntryRes() {}

// Ref: #/components/schemas/EmptyResponse
type EmptyResponse struct{}

func (*EmptyResponse) deleteEntryRes()      {}
func (*EmptyResponse) updateEntryBodyRes()  {}
func (*EmptyResponse) updateEntryTitleRes() {}

type EntryTitlesResponse []string

// Merged schema.
// Ref: #/components/schemas/EntryWithDestTitle
type EntryWithDestTitle struct {
	Path       string       `json:"path"`
	Title      string       `json:"title"`
	Body       string       `json:"body"`
	Visibility string       `json:"visibility"`
	Format     string       `json:"format"`
	ImageUrl   OptNilString `json:"imageUrl"`
	DstTitle   string       `json:"dstTitle"`
}

// GetPath returns the value of Path.
func (s *EntryWithDestTitle) GetPath() string {
	return s.Path
}

// GetTitle returns the value of Title.
func (s *EntryWithDestTitle) GetTitle() string {
	return s.Title
}

// GetBody returns the value of Body.
func (s *EntryWithDestTitle) GetBody() string {
	return s.Body
}

// GetVisibility returns the value of Visibility.
func (s *EntryWithDestTitle) GetVisibility() string {
	return s.Visibility
}

// GetFormat returns the value of Format.
func (s *EntryWithDestTitle) GetFormat() string {
	return s.Format
}

// GetImageUrl returns the value of ImageUrl.
func (s *EntryWithDestTitle) GetImageUrl() OptNilString {
	return s.ImageUrl
}

// GetDstTitle returns the value of DstTitle.
func (s *EntryWithDestTitle) GetDstTitle() string {
	return s.DstTitle
}

// SetPath sets the value of Path.
func (s *EntryWithDestTitle) SetPath(val string) {
	s.Path = val
}

// SetTitle sets the value of Title.
func (s *EntryWithDestTitle) SetTitle(val string) {
	s.Title = val
}

// SetBody sets the value of Body.
func (s *EntryWithDestTitle) SetBody(val string) {
	s.Body = val
}

// SetVisibility sets the value of Visibility.
func (s *EntryWithDestTitle) SetVisibility(val string) {
	s.Visibility = val
}

// SetFormat sets the value of Format.
func (s *EntryWithDestTitle) SetFormat(val string) {
	s.Format = val
}

// SetImageUrl sets the value of ImageUrl.
func (s *EntryWithDestTitle) SetImageUrl(val OptNilString) {
	s.ImageUrl = val
}

// SetDstTitle sets the value of DstTitle.
func (s *EntryWithDestTitle) SetDstTitle(val string) {
	s.DstTitle = val
}

// Ref: #/components/schemas/EntryWithImage
type EntryWithImage struct {
	Path       string       `json:"path"`
	Title      string       `json:"title"`
	Body       string       `json:"body"`
	Visibility string       `json:"visibility"`
	Format     string       `json:"format"`
	ImageUrl   OptNilString `json:"imageUrl"`
}

// GetPath returns the value of Path.
func (s *EntryWithImage) GetPath() string {
	return s.Path
}

// GetTitle returns the value of Title.
func (s *EntryWithImage) GetTitle() string {
	return s.Title
}

// GetBody returns the value of Body.
func (s *EntryWithImage) GetBody() string {
	return s.Body
}

// GetVisibility returns the value of Visibility.
func (s *EntryWithImage) GetVisibility() string {
	return s.Visibility
}

// GetFormat returns the value of Format.
func (s *EntryWithImage) GetFormat() string {
	return s.Format
}

// GetImageUrl returns the value of ImageUrl.
func (s *EntryWithImage) GetImageUrl() OptNilString {
	return s.ImageUrl
}

// SetPath sets the value of Path.
func (s *EntryWithImage) SetPath(val string) {
	s.Path = val
}

// SetTitle sets the value of Title.
func (s *EntryWithImage) SetTitle(val string) {
	s.Title = val
}

// SetBody sets the value of Body.
func (s *EntryWithImage) SetBody(val string) {
	s.Body = val
}

// SetVisibility sets the value of Visibility.
func (s *EntryWithImage) SetVisibility(val string) {
	s.Visibility = val
}

// SetFormat sets the value of Format.
func (s *EntryWithImage) SetFormat(val string) {
	s.Format = val
}

// SetImageUrl sets the value of ImageUrl.
func (s *EntryWithImage) SetImageUrl(val OptNilString) {
	s.ImageUrl = val
}

// Ref: #/components/schemas/ErrorResponse
type ErrorResponse struct {
	Message OptString `json:"message"`
	Error   OptString `json:"error"`
}

// GetMessage returns the value of Message.
func (s *ErrorResponse) GetMessage() OptString {
	return s.Message
}

// GetError returns the value of Error.
func (s *ErrorResponse) GetError() OptString {
	return s.Error
}

// SetMessage sets the value of Message.
func (s *ErrorResponse) SetMessage(val OptString) {
	s.Message = val
}

// SetError sets the value of Error.
func (s *ErrorResponse) SetError(val OptString) {
	s.Error = val
}

func (*ErrorResponse) createEntryRes()           {}
func (*ErrorResponse) deleteEntryRes()           {}
func (*ErrorResponse) getLinkedEntryPathsRes()   {}
func (*ErrorResponse) updateEntryBodyRes()       {}
func (*ErrorResponse) updateEntryVisibilityRes() {}

// ErrorResponseStatusCode wraps ErrorResponse with StatusCode.
type ErrorResponseStatusCode struct {
	StatusCode int
	Response   ErrorResponse
}

// GetStatusCode returns the value of StatusCode.
func (s *ErrorResponseStatusCode) GetStatusCode() int {
	return s.StatusCode
}

// GetResponse returns the value of Response.
func (s *ErrorResponseStatusCode) GetResponse() ErrorResponse {
	return s.Response
}

// SetStatusCode sets the value of StatusCode.
func (s *ErrorResponseStatusCode) SetStatusCode(val int) {
	s.StatusCode = val
}

// SetResponse sets the value of Response.
func (s *ErrorResponseStatusCode) SetResponse(val ErrorResponse) {
	s.Response = val
}

// Ref: #/components/schemas/GetLatestEntriesRow
type GetLatestEntriesRow struct {
	Path         OptString      `json:"Path"`
	Title        OptString      `json:"Title"`
	Body         OptString      `json:"Body"`
	Visibility   OptString      `json:"Visibility"`
	Format       OptString      `json:"Format"`
	PublishedAt  OptNilDateTime `json:"PublishedAt"`
	LastEditedAt OptNilDateTime `json:"LastEditedAt"`
	CreatedAt    OptNilDateTime `json:"CreatedAt"`
	UpdatedAt    OptNilDateTime `json:"UpdatedAt"`
	ImageUrl     OptNilString   `json:"ImageUrl"`
}

// GetPath returns the value of Path.
func (s *GetLatestEntriesRow) GetPath() OptString {
	return s.Path
}

// GetTitle returns the value of Title.
func (s *GetLatestEntriesRow) GetTitle() OptString {
	return s.Title
}

// GetBody returns the value of Body.
func (s *GetLatestEntriesRow) GetBody() OptString {
	return s.Body
}

// GetVisibility returns the value of Visibility.
func (s *GetLatestEntriesRow) GetVisibility() OptString {
	return s.Visibility
}

// GetFormat returns the value of Format.
func (s *GetLatestEntriesRow) GetFormat() OptString {
	return s.Format
}

// GetPublishedAt returns the value of PublishedAt.
func (s *GetLatestEntriesRow) GetPublishedAt() OptNilDateTime {
	return s.PublishedAt
}

// GetLastEditedAt returns the value of LastEditedAt.
func (s *GetLatestEntriesRow) GetLastEditedAt() OptNilDateTime {
	return s.LastEditedAt
}

// GetCreatedAt returns the value of CreatedAt.
func (s *GetLatestEntriesRow) GetCreatedAt() OptNilDateTime {
	return s.CreatedAt
}

// GetUpdatedAt returns the value of UpdatedAt.
func (s *GetLatestEntriesRow) GetUpdatedAt() OptNilDateTime {
	return s.UpdatedAt
}

// GetImageUrl returns the value of ImageUrl.
func (s *GetLatestEntriesRow) GetImageUrl() OptNilString {
	return s.ImageUrl
}

// SetPath sets the value of Path.
func (s *GetLatestEntriesRow) SetPath(val OptString) {
	s.Path = val
}

// SetTitle sets the value of Title.
func (s *GetLatestEntriesRow) SetTitle(val OptString) {
	s.Title = val
}

// SetBody sets the value of Body.
func (s *GetLatestEntriesRow) SetBody(val OptString) {
	s.Body = val
}

// SetVisibility sets the value of Visibility.
func (s *GetLatestEntriesRow) SetVisibility(val OptString) {
	s.Visibility = val
}

// SetFormat sets the value of Format.
func (s *GetLatestEntriesRow) SetFormat(val OptString) {
	s.Format = val
}

// SetPublishedAt sets the value of PublishedAt.
func (s *GetLatestEntriesRow) SetPublishedAt(val OptNilDateTime) {
	s.PublishedAt = val
}

// SetLastEditedAt sets the value of LastEditedAt.
func (s *GetLatestEntriesRow) SetLastEditedAt(val OptNilDateTime) {
	s.LastEditedAt = val
}

// SetCreatedAt sets the value of CreatedAt.
func (s *GetLatestEntriesRow) SetCreatedAt(val OptNilDateTime) {
	s.CreatedAt = val
}

// SetUpdatedAt sets the value of UpdatedAt.
func (s *GetLatestEntriesRow) SetUpdatedAt(val OptNilDateTime) {
	s.UpdatedAt = val
}

// SetImageUrl sets the value of ImageUrl.
func (s *GetLatestEntriesRow) SetImageUrl(val OptNilString) {
	s.ImageUrl = val
}

// Ref: #/components/schemas/LinkPalletData
type LinkPalletData struct {
	// Array of potential new link titles.
	NewLinks []string `json:"newLinks"`
	// Array of directly linked entries.
	Links []EntryWithImage `json:"links"`
	// Array of two-hop link relationships.
	Twohops []TwoHopLink `json:"twohops"`
}

// GetNewLinks returns the value of NewLinks.
func (s *LinkPalletData) GetNewLinks() []string {
	return s.NewLinks
}

// GetLinks returns the value of Links.
func (s *LinkPalletData) GetLinks() []EntryWithImage {
	return s.Links
}

// GetTwohops returns the value of Twohops.
func (s *LinkPalletData) GetTwohops() []TwoHopLink {
	return s.Twohops
}

// SetNewLinks sets the value of NewLinks.
func (s *LinkPalletData) SetNewLinks(val []string) {
	s.NewLinks = val
}

// SetLinks sets the value of Links.
func (s *LinkPalletData) SetLinks(val []EntryWithImage) {
	s.Links = val
}

// SetTwohops sets the value of Twohops.
func (s *LinkPalletData) SetTwohops(val []TwoHopLink) {
	s.Twohops = val
}

// Object where keys are lowercase destination entry titles and values are their paths (null if entry
// doesn't exist)..
// Ref: #/components/schemas/LinkedEntryPathsResponse
type LinkedEntryPathsResponse map[string]string

func (s *LinkedEntryPathsResponse) init() LinkedEntryPathsResponse {
	m := *s
	if m == nil {
		m = map[string]string{}
		*s = m
	}
	return m
}

func (*LinkedEntryPathsResponse) getLinkedEntryPathsRes() {}

// NewOptDateTime returns new OptDateTime with value set to v.
func NewOptDateTime(v time.Time) OptDateTime {
	return OptDateTime{
		Value: v,
		Set:   true,
	}
}

// OptDateTime is optional time.Time.
type OptDateTime struct {
	Value time.Time
	Set   bool
}

// IsSet returns true if OptDateTime was set.
func (o OptDateTime) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptDateTime) Reset() {
	var v time.Time
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *OptDateTime) SetTo(v time.Time) {
	o.Set = true
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptDateTime) Get() (v time.Time, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptDateTime) Or(d time.Time) time.Time {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

// NewOptNilDateTime returns new OptNilDateTime with value set to v.
func NewOptNilDateTime(v time.Time) OptNilDateTime {
	return OptNilDateTime{
		Value: v,
		Set:   true,
	}
}

// OptNilDateTime is optional nullable time.Time.
type OptNilDateTime struct {
	Value time.Time
	Set   bool
	Null  bool
}

// IsSet returns true if OptNilDateTime was set.
func (o OptNilDateTime) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptNilDateTime) Reset() {
	var v time.Time
	o.Value = v
	o.Set = false
	o.Null = false
}

// SetTo sets value to v.
func (o *OptNilDateTime) SetTo(v time.Time) {
	o.Set = true
	o.Null = false
	o.Value = v
}

// IsSet returns true if value is Null.
func (o OptNilDateTime) IsNull() bool { return o.Null }

// SetNull sets value to null.
func (o *OptNilDateTime) SetToNull() {
	o.Set = true
	o.Null = true
	var v time.Time
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptNilDateTime) Get() (v time.Time, ok bool) {
	if o.Null {
		return v, false
	}
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptNilDateTime) Or(d time.Time) time.Time {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

// NewOptNilString returns new OptNilString with value set to v.
func NewOptNilString(v string) OptNilString {
	return OptNilString{
		Value: v,
		Set:   true,
	}
}

// OptNilString is optional nullable string.
type OptNilString struct {
	Value string
	Set   bool
	Null  bool
}

// IsSet returns true if OptNilString was set.
func (o OptNilString) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptNilString) Reset() {
	var v string
	o.Value = v
	o.Set = false
	o.Null = false
}

// SetTo sets value to v.
func (o *OptNilString) SetTo(v string) {
	o.Set = true
	o.Null = false
	o.Value = v
}

// IsSet returns true if value is Null.
func (o OptNilString) IsNull() bool { return o.Null }

// SetNull sets value to null.
func (o *OptNilString) SetToNull() {
	o.Set = true
	o.Null = true
	var v string
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptNilString) Get() (v string, ok bool) {
	if o.Null {
		return v, false
	}
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptNilString) Or(d string) string {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

// NewOptString returns new OptString with value set to v.
func NewOptString(v string) OptString {
	return OptString{
		Value: v,
		Set:   true,
	}
}

// OptString is optional string.
type OptString struct {
	Value string
	Set   bool
}

// IsSet returns true if OptString was set.
func (o OptString) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptString) Reset() {
	var v string
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *OptString) SetTo(v string) {
	o.Set = true
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptString) Get() (v string, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptString) Or(d string) string {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

// Ref: #/components/schemas/TwoHopLink
type TwoHopLink struct {
	Src   EntryWithDestTitle `json:"src"`
	Links []EntryWithImage   `json:"links"`
}

// GetSrc returns the value of Src.
func (s *TwoHopLink) GetSrc() EntryWithDestTitle {
	return s.Src
}

// GetLinks returns the value of Links.
func (s *TwoHopLink) GetLinks() []EntryWithImage {
	return s.Links
}

// SetSrc sets the value of Src.
func (s *TwoHopLink) SetSrc(val EntryWithDestTitle) {
	s.Src = val
}

// SetLinks sets the value of Links.
func (s *TwoHopLink) SetLinks(val []EntryWithImage) {
	s.Links = val
}

// Ref: #/components/schemas/UpdateEntryBodyRequest
type UpdateEntryBodyRequest struct {
	// The new content of the entry.
	Body string `json:"body"`
}

// GetBody returns the value of Body.
func (s *UpdateEntryBodyRequest) GetBody() string {
	return s.Body
}

// SetBody sets the value of Body.
func (s *UpdateEntryBodyRequest) SetBody(val string) {
	s.Body = val
}

type UpdateEntryTitleConflict ErrorResponse

func (*UpdateEntryTitleConflict) updateEntryTitleRes() {}

type UpdateEntryTitleNotFound ErrorResponse

func (*UpdateEntryTitleNotFound) updateEntryTitleRes() {}

// Ref: #/components/schemas/UpdateEntryTitleRequest
type UpdateEntryTitleRequest struct {
	// The new title for the entry.
	Title string `json:"title"`
}

// GetTitle returns the value of Title.
func (s *UpdateEntryTitleRequest) GetTitle() string {
	return s.Title
}

// SetTitle sets the value of Title.
func (s *UpdateEntryTitleRequest) SetTitle(val string) {
	s.Title = val
}

// Ref: #/components/schemas/UpdateVisibilityRequest
type UpdateVisibilityRequest struct {
	// The new visibility status for the entry.
	Visibility string `json:"visibility"`
}

// GetVisibility returns the value of Visibility.
func (s *UpdateVisibilityRequest) GetVisibility() string {
	return s.Visibility
}

// SetVisibility sets the value of Visibility.
func (s *UpdateVisibilityRequest) SetVisibility(val string) {
	s.Visibility = val
}

// Ref: #/components/schemas/UpdateVisibilityResponse
type UpdateVisibilityResponse struct {
	// The new visibility status for the entry.
	Visibility string `json:"visibility"`
}

// GetVisibility returns the value of Visibility.
func (s *UpdateVisibilityResponse) GetVisibility() string {
	return s.Visibility
}

// SetVisibility sets the value of Visibility.
func (s *UpdateVisibilityResponse) SetVisibility(val string) {
	s.Visibility = val
}

func (*UpdateVisibilityResponse) updateEntryVisibilityRes() {}
