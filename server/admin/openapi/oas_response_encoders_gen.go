// Code generated by ogen, DO NOT EDIT.

package openapi

import (
	"net/http"

	"github.com/go-faster/errors"
	"github.com/go-faster/jx"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	ht "github.com/ogen-go/ogen/http"
)

func encodeCreateEntryResponse(response CreateEntryRes, w http.ResponseWriter, span trace.Span) error {
	switch response := response.(type) {
	case *CreateEntryResponse:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(201)
		span.SetStatus(codes.Ok, http.StatusText(201))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *ErrorResponse:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(409)
		span.SetStatus(codes.Error, http.StatusText(409))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	default:
		return errors.Errorf("unexpected response type: %T", response)
	}
}

func encodeGetAllEntryTitlesResponse(response EntryTitlesResponse, w http.ResponseWriter, span trace.Span) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)
	span.SetStatus(codes.Ok, http.StatusText(200))

	e := new(jx.Encoder)
	response.Encode(e)
	if _, err := e.WriteTo(w); err != nil {
		return errors.Wrap(err, "write")
	}

	return nil
}

func encodeGetEntryByDynamicPathResponse(response *GetLatestEntriesRow, w http.ResponseWriter, span trace.Span) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)
	span.SetStatus(codes.Ok, http.StatusText(200))

	e := new(jx.Encoder)
	response.Encode(e)
	if _, err := e.WriteTo(w); err != nil {
		return errors.Wrap(err, "write")
	}

	return nil
}

func encodeGetLatestEntriesResponse(response []GetLatestEntriesRow, w http.ResponseWriter, span trace.Span) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)
	span.SetStatus(codes.Ok, http.StatusText(200))

	e := new(jx.Encoder)
	e.ArrStart()
	for _, elem := range response {
		elem.Encode(e)
	}
	e.ArrEnd()
	if _, err := e.WriteTo(w); err != nil {
		return errors.Wrap(err, "write")
	}

	return nil
}

func encodeGetLinkedEntryPathsResponse(response LinkedEntriesResponse, w http.ResponseWriter, span trace.Span) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)
	span.SetStatus(codes.Ok, http.StatusText(200))

	e := new(jx.Encoder)
	response.Encode(e)
	if _, err := e.WriteTo(w); err != nil {
		return errors.Wrap(err, "write")
	}

	return nil
}

func encodeUpdateEntryBodyResponse(response UpdateEntryBodyRes, w http.ResponseWriter, span trace.Span) error {
	switch response := response.(type) {
	case *EmptyResponse:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(200)
		span.SetStatus(codes.Ok, http.StatusText(200))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *ErrorResponse:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(404)
		span.SetStatus(codes.Error, http.StatusText(404))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	default:
		return errors.Errorf("unexpected response type: %T", response)
	}
}

func encodeUpdateEntryTitleResponse(response UpdateEntryTitleRes, w http.ResponseWriter, span trace.Span) error {
	switch response := response.(type) {
	case *EmptyResponse:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(200)
		span.SetStatus(codes.Ok, http.StatusText(200))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *UpdateEntryTitleNotFound:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(404)
		span.SetStatus(codes.Error, http.StatusText(404))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *UpdateEntryTitleConflict:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(409)
		span.SetStatus(codes.Error, http.StatusText(409))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	default:
		return errors.Errorf("unexpected response type: %T", response)
	}
}

func encodeErrorResponse(response *ErrorResponseStatusCode, w http.ResponseWriter, span trace.Span) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	code := response.StatusCode
	if code == 0 {
		// Set default status code.
		code = http.StatusOK
	}
	w.WriteHeader(code)
	if st := http.StatusText(code); code >= http.StatusBadRequest {
		span.SetStatus(codes.Error, st)
	} else {
		span.SetStatus(codes.Ok, st)
	}

	e := new(jx.Encoder)
	response.Response.Encode(e)
	if _, err := e.WriteTo(w); err != nil {
		return errors.Wrap(err, "write")
	}

	if code >= http.StatusInternalServerError {
		return errors.Wrapf(ht.ErrInternalServerErrorResponse, "code: %d, message: %s", code, http.StatusText(code))
	}
	return nil

}
