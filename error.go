package es

import (
	xerr "github.com/goclub/error"
)

func newDecodeSearchResultError(decodeError error, searchHit []byte, element any, decodeFailMessage string) (err error) {
	return xerr.WithStack(&DecodeSearchResultError{
		SearchHit:         searchHit,
		Element:           element,
		DecodeFailMessage: decodeFailMessage,
		DecodeError:       decodeError,
	})
}

type DecodeSearchResultError struct {
	SearchHit         []byte
	Element           any
	DecodeFailMessage string
	DecodeError       error
}

func (e *DecodeSearchResultError) Error() string {
	return `goclub/es: SearchSlice7 json Decode fail:` + e.DecodeFailMessage
}
func (e *DecodeSearchResultError) Unwrap() error { return e.DecodeError }

func AsDecodeSearchResultError(err error) (decodeSearchResultError *DecodeSearchResultError, as bool) {
	var ptr *DecodeSearchResultError
	if xerr.As(err, &ptr) {
		return ptr, true
	} else {
		return nil, false
	}
}
