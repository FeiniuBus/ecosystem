package compression

import (
	"net/http"
	"sort"
	"strings"
)

// MimeTypes is default MIME types to compress response for.
var MimeTypes = []string{
	"text/plain",
	"text/css",
	"application/javascript",
	"text/html",
	"application/xml",
	"text/xml",
	"application/json",
	"text/json",
}

func init() {
	sort.Strings(MimeTypes)
}

type compressionProvider struct {
	mimeTypes []string
}

func (c *compressionProvider) shouldCompressResponse(w http.ResponseWriter) bool {
	contentRange := w.Header().Get("Content-Range")
	if len(contentRange) != 0 {
		return false
	}

	mimeType := w.Header().Get("Content-Type")
	if len(mimeType) == 0 {
		return false
	}

	separator := strings.Index(mimeType, ";")
	if separator > 0 {
		mimeType = mimeType[0:separator]
		mimeType = strings.TrimSpace(mimeType)
	}

	i := sort.SearchStrings(MimeTypes, mimeType)
	return i < len(MimeTypes)
}
