// Code generated by ogen, DO NOT EDIT.

package openapi

import (
	"net/http"
	"net/url"
	"time"

	"github.com/go-faster/errors"

	"github.com/ogen-go/ogen/conv"
	"github.com/ogen-go/ogen/middleware"
	"github.com/ogen-go/ogen/ogenerrors"
	"github.com/ogen-go/ogen/uri"
	"github.com/ogen-go/ogen/validate"
)

// DeleteEntryParams is parameters of deleteEntry operation.
type DeleteEntryParams struct {
	// The path of the entry to delete.
	Path string
}

func unpackDeleteEntryParams(packed middleware.Parameters) (params DeleteEntryParams) {
	{
		key := middleware.ParameterKey{
			Name: "path",
			In:   "path",
		}
		params.Path = packed[key].(string)
	}
	return params
}

func decodeDeleteEntryParams(args [1]string, argsEscaped bool, r *http.Request) (params DeleteEntryParams, _ error) {
	// Decode path: path.
	if err := func() error {
		param := args[0]
		if argsEscaped {
			unescaped, err := url.PathUnescape(args[0])
			if err != nil {
				return errors.Wrap(err, "unescape path")
			}
			param = unescaped
		}
		if len(param) > 0 {
			d := uri.NewPathDecoder(uri.PathDecoderConfig{
				Param:   "path",
				Value:   param,
				Style:   uri.PathStyleSimple,
				Explode: false,
			})

			if err := func() error {
				val, err := d.DecodeValue()
				if err != nil {
					return err
				}

				c, err := conv.ToString(val)
				if err != nil {
					return err
				}

				params.Path = c
				return nil
			}(); err != nil {
				return err
			}
		} else {
			return validate.ErrFieldRequired
		}
		return nil
	}(); err != nil {
		return params, &ogenerrors.DecodeParamError{
			Name: "path",
			In:   "path",
			Err:  err,
		}
	}
	return params, nil
}

// GetEntryByDynamicPathParams is parameters of getEntryByDynamicPath operation.
type GetEntryByDynamicPathParams struct {
	// The path of the entry.
	Path string
}

func unpackGetEntryByDynamicPathParams(packed middleware.Parameters) (params GetEntryByDynamicPathParams) {
	{
		key := middleware.ParameterKey{
			Name: "path",
			In:   "path",
		}
		params.Path = packed[key].(string)
	}
	return params
}

func decodeGetEntryByDynamicPathParams(args [1]string, argsEscaped bool, r *http.Request) (params GetEntryByDynamicPathParams, _ error) {
	// Decode path: path.
	if err := func() error {
		param := args[0]
		if argsEscaped {
			unescaped, err := url.PathUnescape(args[0])
			if err != nil {
				return errors.Wrap(err, "unescape path")
			}
			param = unescaped
		}
		if len(param) > 0 {
			d := uri.NewPathDecoder(uri.PathDecoderConfig{
				Param:   "path",
				Value:   param,
				Style:   uri.PathStyleSimple,
				Explode: false,
			})

			if err := func() error {
				val, err := d.DecodeValue()
				if err != nil {
					return err
				}

				c, err := conv.ToString(val)
				if err != nil {
					return err
				}

				params.Path = c
				return nil
			}(); err != nil {
				return err
			}
		} else {
			return validate.ErrFieldRequired
		}
		return nil
	}(); err != nil {
		return params, &ogenerrors.DecodeParamError{
			Name: "path",
			In:   "path",
			Err:  err,
		}
	}
	return params, nil
}

// GetLatestEntriesParams is parameters of getLatestEntries operation.
type GetLatestEntriesParams struct {
	// Filter entries by the last edited date.
	LastLastEditedAt OptDateTime
}

func unpackGetLatestEntriesParams(packed middleware.Parameters) (params GetLatestEntriesParams) {
	{
		key := middleware.ParameterKey{
			Name: "last_last_edited_at",
			In:   "query",
		}
		if v, ok := packed[key]; ok {
			params.LastLastEditedAt = v.(OptDateTime)
		}
	}
	return params
}

func decodeGetLatestEntriesParams(args [0]string, argsEscaped bool, r *http.Request) (params GetLatestEntriesParams, _ error) {
	q := uri.NewQueryDecoder(r.URL.Query())
	// Decode query: last_last_edited_at.
	if err := func() error {
		cfg := uri.QueryParameterDecodingConfig{
			Name:    "last_last_edited_at",
			Style:   uri.QueryStyleForm,
			Explode: true,
		}

		if err := q.HasParam(cfg); err == nil {
			if err := q.DecodeParam(cfg, func(d uri.Decoder) error {
				var paramsDotLastLastEditedAtVal time.Time
				if err := func() error {
					val, err := d.DecodeValue()
					if err != nil {
						return err
					}

					c, err := conv.ToDateTime(val)
					if err != nil {
						return err
					}

					paramsDotLastLastEditedAtVal = c
					return nil
				}(); err != nil {
					return err
				}
				params.LastLastEditedAt.SetTo(paramsDotLastLastEditedAtVal)
				return nil
			}); err != nil {
				return err
			}
		}
		return nil
	}(); err != nil {
		return params, &ogenerrors.DecodeParamError{
			Name: "last_last_edited_at",
			In:   "query",
			Err:  err,
		}
	}
	return params, nil
}

// GetLinkPalletParams is parameters of getLinkPallet operation.
type GetLinkPalletParams struct {
	// The source entry path.
	Path string
}

func unpackGetLinkPalletParams(packed middleware.Parameters) (params GetLinkPalletParams) {
	{
		key := middleware.ParameterKey{
			Name: "path",
			In:   "path",
		}
		params.Path = packed[key].(string)
	}
	return params
}

func decodeGetLinkPalletParams(args [1]string, argsEscaped bool, r *http.Request) (params GetLinkPalletParams, _ error) {
	// Decode path: path.
	if err := func() error {
		param := args[0]
		if argsEscaped {
			unescaped, err := url.PathUnescape(args[0])
			if err != nil {
				return errors.Wrap(err, "unescape path")
			}
			param = unescaped
		}
		if len(param) > 0 {
			d := uri.NewPathDecoder(uri.PathDecoderConfig{
				Param:   "path",
				Value:   param,
				Style:   uri.PathStyleSimple,
				Explode: false,
			})

			if err := func() error {
				val, err := d.DecodeValue()
				if err != nil {
					return err
				}

				c, err := conv.ToString(val)
				if err != nil {
					return err
				}

				params.Path = c
				return nil
			}(); err != nil {
				return err
			}
		} else {
			return validate.ErrFieldRequired
		}
		return nil
	}(); err != nil {
		return params, &ogenerrors.DecodeParamError{
			Name: "path",
			In:   "path",
			Err:  err,
		}
	}
	return params, nil
}

// GetLinkedEntryPathsParams is parameters of getLinkedEntryPaths operation.
type GetLinkedEntryPathsParams struct {
	// The source entry path.
	Path string
}

func unpackGetLinkedEntryPathsParams(packed middleware.Parameters) (params GetLinkedEntryPathsParams) {
	{
		key := middleware.ParameterKey{
			Name: "path",
			In:   "path",
		}
		params.Path = packed[key].(string)
	}
	return params
}

func decodeGetLinkedEntryPathsParams(args [1]string, argsEscaped bool, r *http.Request) (params GetLinkedEntryPathsParams, _ error) {
	// Decode path: path.
	if err := func() error {
		param := args[0]
		if argsEscaped {
			unescaped, err := url.PathUnescape(args[0])
			if err != nil {
				return errors.Wrap(err, "unescape path")
			}
			param = unescaped
		}
		if len(param) > 0 {
			d := uri.NewPathDecoder(uri.PathDecoderConfig{
				Param:   "path",
				Value:   param,
				Style:   uri.PathStyleSimple,
				Explode: false,
			})

			if err := func() error {
				val, err := d.DecodeValue()
				if err != nil {
					return err
				}

				c, err := conv.ToString(val)
				if err != nil {
					return err
				}

				params.Path = c
				return nil
			}(); err != nil {
				return err
			}
		} else {
			return validate.ErrFieldRequired
		}
		return nil
	}(); err != nil {
		return params, &ogenerrors.DecodeParamError{
			Name: "path",
			In:   "path",
			Err:  err,
		}
	}
	return params, nil
}

// UpdateEntryBodyParams is parameters of updateEntryBody operation.
type UpdateEntryBodyParams struct {
	// The path of the entry to update.
	Path string
}

func unpackUpdateEntryBodyParams(packed middleware.Parameters) (params UpdateEntryBodyParams) {
	{
		key := middleware.ParameterKey{
			Name: "path",
			In:   "path",
		}
		params.Path = packed[key].(string)
	}
	return params
}

func decodeUpdateEntryBodyParams(args [1]string, argsEscaped bool, r *http.Request) (params UpdateEntryBodyParams, _ error) {
	// Decode path: path.
	if err := func() error {
		param := args[0]
		if argsEscaped {
			unescaped, err := url.PathUnescape(args[0])
			if err != nil {
				return errors.Wrap(err, "unescape path")
			}
			param = unescaped
		}
		if len(param) > 0 {
			d := uri.NewPathDecoder(uri.PathDecoderConfig{
				Param:   "path",
				Value:   param,
				Style:   uri.PathStyleSimple,
				Explode: false,
			})

			if err := func() error {
				val, err := d.DecodeValue()
				if err != nil {
					return err
				}

				c, err := conv.ToString(val)
				if err != nil {
					return err
				}

				params.Path = c
				return nil
			}(); err != nil {
				return err
			}
		} else {
			return validate.ErrFieldRequired
		}
		return nil
	}(); err != nil {
		return params, &ogenerrors.DecodeParamError{
			Name: "path",
			In:   "path",
			Err:  err,
		}
	}
	return params, nil
}

// UpdateEntryTitleParams is parameters of updateEntryTitle operation.
type UpdateEntryTitleParams struct {
	// The path of the entry to update.
	Path string
}

func unpackUpdateEntryTitleParams(packed middleware.Parameters) (params UpdateEntryTitleParams) {
	{
		key := middleware.ParameterKey{
			Name: "path",
			In:   "path",
		}
		params.Path = packed[key].(string)
	}
	return params
}

func decodeUpdateEntryTitleParams(args [1]string, argsEscaped bool, r *http.Request) (params UpdateEntryTitleParams, _ error) {
	// Decode path: path.
	if err := func() error {
		param := args[0]
		if argsEscaped {
			unescaped, err := url.PathUnescape(args[0])
			if err != nil {
				return errors.Wrap(err, "unescape path")
			}
			param = unescaped
		}
		if len(param) > 0 {
			d := uri.NewPathDecoder(uri.PathDecoderConfig{
				Param:   "path",
				Value:   param,
				Style:   uri.PathStyleSimple,
				Explode: false,
			})

			if err := func() error {
				val, err := d.DecodeValue()
				if err != nil {
					return err
				}

				c, err := conv.ToString(val)
				if err != nil {
					return err
				}

				params.Path = c
				return nil
			}(); err != nil {
				return err
			}
		} else {
			return validate.ErrFieldRequired
		}
		return nil
	}(); err != nil {
		return params, &ogenerrors.DecodeParamError{
			Name: "path",
			In:   "path",
			Err:  err,
		}
	}
	return params, nil
}
