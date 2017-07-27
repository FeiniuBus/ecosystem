package errwrap

import (
	"errors"
	"fmt"
	"net/http"
)

// ErrorWrap is a http error wrap
type ErrorWrap struct {
	StatusCode int
	What       error
}

func (w *ErrorWrap) MakeError() error {
	return fmt.Errorf("%s: %s", http.StatusText(w.StatusCode), w.What.Error())
}

// Wrap defines a ErrorWrap
func Wrap(code int, err string) *ErrorWrap {
	return &ErrorWrap{
		StatusCode: code,
		What:       errors.New(err),
	}
}

// Wrapf defines a ErrorWrap
func Wrapf(code int, format string, a ...interface{}) *ErrorWrap {
	err := fmt.Errorf(format, a...)
	return &ErrorWrap{
		StatusCode: code,
		What:       err,
	}
}

// WrapError defines a ErrorWrap
func WrapError(code int, err error) *ErrorWrap {
	return &ErrorWrap{
		StatusCode: code,
		What:       err,
	}
}
