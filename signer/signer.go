package signer

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	timeFormat = "20060102T150405Z"
	// emptyStringSHA256 is a SHA256 of an empty string
	emptyStringSHA256 = `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855`
)

var noEscape [256]bool

func init() {
	for i := 0; i < len(noEscape); i++ {
		noEscape[i] = (i >= 'A' && i <= 'Z') ||
			(i >= 'a' && i <= 'z') ||
			(i >= '0' && i <= '9') ||
			i == '-' ||
			i == '.' ||
			i == '_' ||
			i == '~'
	}
}

// UseSignature method
func UseSignature(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if len(auth) == 0 {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("请求未签名"))
			return
		}

		if !strings.HasPrefix(auth, "FNSIGN") {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("签名格式不正确"))
			return
		}

		token := auth[len("FNSIGN ")-1:]
		segments := strings.Split(token, ",")
		if len(segments) != 2 {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("签名数据不符合规定格式"))
			return
		}

		accessKeyID := strings.Split(segments[0], "=")[1]
		if len(accessKeyID) == 0 {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("无法识别用户的access_key"))
			return
		}

		h.ServeHTTP(w, r)
	})
}

type credentials struct {
	AccessKeyID     string
	SecretAccessKey string
}

type signingCtx struct {
	Request *http.Request
	Body    io.ReadSeeker
	Query   url.Values

	formattedTime string
	credValues    credentials

	bodyDigest       string
	canonicalString  string
	credentialString string
	stringToSign     string
	signature        string
}

func getSecretAccessKey(accessKeyID string) (string, error) {
	return "", nil
}

func (ctx *signingCtx) buildTime() error {
	stime := ctx.Request.Header.Get("x-feiniubus-date")
	if stime == "" {
		return errors.New("Can't resolve signature time")
	}

	ctx.formattedTime = stime
	return nil
}

func (ctx *signingCtx) buildBodyDigest() {
	var hash string
	if ctx.Body == nil {
		hash = emptyStringSHA256
	} else {
		hash = hex.EncodeToString(makeSha256Reader(ctx.Body))
	}

	ctx.bodyDigest = hash
}

func (ctx *signingCtx) buildCanonicalString() {
	ctx.Request.URL.RawQuery = strings.Replace(ctx.Query.Encode(), "+", "%20", -1)

	uri := getURIPath(ctx.Request.URL)
	uri = escapePath(uri, false)

	ctx.canonicalString = strings.Join([]string{
		ctx.Request.Method,
		uri,
		ctx.Request.URL.RawQuery,
		ctx.bodyDigest,
	}, "\n")
}

func (ctx *signingCtx) buildStringToSign() {
	var prefix = strings.Join([]string{
		"HMAC-SHA256",
		ctx.formattedTime,
	}, "-")

	ctx.stringToSign = strings.Join([]string{
		prefix,
		ctx.credValues.AccessKeyID,
		hex.EncodeToString(makeSha256([]byte(ctx.canonicalString))),
	}, "\n")
}

func (ctx *signingCtx) buildSignature() {
	secret := ctx.credValues.SecretAccessKey
	date := makeHmac([]byte("FNSIGN"+secret), []byte(ctx.formattedTime))
	credentials := makeHmac(date, []byte("feiniubus_request"))
	signature := makeHmac(credentials, []byte(ctx.stringToSign))
	ctx.signature = hex.EncodeToString(signature)
}

func escapePath(path string, encodeSep bool) string {
	var buf bytes.Buffer
	for i := 0; i < len(path); i++ {
		c := path[i]
		if noEscape[c] || (c == '/' && !encodeSep) {
			buf.WriteByte(c)
		} else {
			fmt.Fprintf(&buf, "%%%02X", c)
		}
	}
	return buf.String()
}

func getURIPath(u *url.URL) string {
	var uri string

	if len(u.Opaque) > 0 {
		uri = "/" + strings.Join(strings.Split(u.Opaque, "/")[3:], "/")
	} else {
		uri = u.EscapedPath()
	}

	if len(uri) == 0 {
		uri = "/"
	}

	return uri
}

func makeHmac(key []byte, data []byte) []byte {
	hash := hmac.New(sha256.New, key)
	hash.Write(data)
	return hash.Sum(nil)
}

func makeSha256(data []byte) []byte {
	hash := sha256.New()
	hash.Write(data)
	return hash.Sum(nil)
}

func makeSha256Reader(reader io.ReadSeeker) []byte {
	hash := sha256.New()
	start, _ := reader.Seek(0, 1)
	defer reader.Seek(start, 0)

	io.Copy(hash, reader)
	return hash.Sum(nil)
}
