package buffering

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
)

// ErrNotSupport error
var ErrNotSupport = errors.New("not support")

// RequestBody is a request body that cat be reuseable
type RequestBody struct {
	Buffer *bytes.Reader
	body   io.ReadCloser
}

// Read method
func (b *RequestBody) Read(p []byte) (int, error) {
	if b.Buffer != nil {
		return b.Buffer.Read(p)
	}

	return 0, io.EOF
}

// Seek method
func (b *RequestBody) Seek(offset int64, whence int) (int64, error) {
	if b.Buffer != nil {
		return b.Buffer.Seek(offset, whence)
	}

	return 0, ErrNotSupport
}

// Close method
func (b *RequestBody) Close() error {
	return b.body.Close()
}

// NewRequestBody returns a new RequestBody pointer
func NewRequestBody(r *http.Request) *RequestBody {
	body := &RequestBody{
		body: r.Body,
	}

	if !mayContainRequestBody(r.Method) {
		body.Buffer = nil
		return body
	}

	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		body.Buffer = nil
	} else {
		body.Buffer = bytes.NewReader(content)
	}

	return body
}

// UseBuffering method
func UseBuffering(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b := r.Body
		r.Body = NewRequestBody(r)

		h.ServeHTTP(w, r)

		r.Body = b
	})
}

func mayContainRequestBody(method string) bool {
	return method == "POST" || method == "PUT" || method == "PATCH"
}
