package api

import "fmt"

type StoreErr struct {
	err error
}

func NewStoreErr(err error) *StoreErr {
	return &StoreErr{err: err}
}

func (e StoreErr) Error() string {
	if e.err != nil {
		return e.err.Error()
	}

	return "unknown store error"
}

func (e StoreErr) Unwrap() error { return e.err }

type ParamErr struct {
	err error
}

func NewParamErr(err error) *ParamErr {
	return &ParamErr{err: err}
}

func (e ParamErr) Error() string {
	if e.err != nil {
		return e.err.Error()
	}

	return "unknown param error"
}

func (e ParamErr) Unwrap() error { return e.err }

type HTTPErr struct {
	HTTPCode int
	Code     Code
	Err      error
}

func NewHTTPErr(httpCode int, code Code, err error) *HTTPErr {
	return &HTTPErr{
		HTTPCode: httpCode,
		Code:     code,
		Err:      err,
	}
}

func (e HTTPErr) Error() string {
	desc := fmt.Sprintf("HTTP: %d (%d) ", e.HTTPCode, e.Code)

	if e.Err != nil {
		desc += ": " + e.Err.Error()
	} else {
		desc += "unknown error"
	}

	return desc
}

func (e HTTPErr) Unwrap() error { return e.Err }
