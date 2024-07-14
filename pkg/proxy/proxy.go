package proxy

import (
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/vietanhduong/bucket-proxy/pkg/bucket"
	"github.com/vietanhduong/bucket-proxy/pkg/bucket/types"
	"github.com/vietanhduong/bucket-proxy/pkg/logging"
)

var log = logging.WithField("pkg", "pkg/proxy")

type Proxy struct {
	bucket bucket.Interface
}

func New(bucket bucket.Interface) *Proxy {
	return &Proxy{bucket: bucket}
}

func (p *Proxy) HttpHandler() (string, http.Handler) {
	return "/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { p.proxy(w, r) })
}

func (p *Proxy) proxy(w http.ResponseWriter, r *http.Request) {
	l := log.WithField("path", r.URL.Path)
	attrs, err := p.bucket.ObjectMetadata(r.Context(), strings.TrimPrefix(r.URL.Path, "/"))
	if err != nil {
		l.WithError(err).Error("could not get object metadata")
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if attrs == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	if last := handleModifySince(r); !last.IsZero() {
		if !attrs.Updated.Truncate(time.Second).After(last) {
			w.WriteHeader(304)
			return
		}
	}

	opts := types.DownloadOptions{
		AcceptCompress: strings.Contains(r.Header.Get(HttpAcceptEncoding), "gzip"),
	}

	reader, err := p.bucket.Download(r.Context(), attrs.Name, opts)
	if err != nil {
		l.WithError(err).Error("could not download object")
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	setHeader(w, HttpLastModified, attrs.Updated)
	setHeader(w, HttpContentType, attrs.ContentType)
	setHeader(w, HttpContentLanguage, attrs.ContentLanguage)
	setHeader(w, HttpCacheControl, attrs.CacheControl)
	setHeader(w, HttpContentEncoding, reader.ContentEncoding)
	setHeader(w, HttpContentDisposition, attrs.ContentDisposition)
	setHeader(w, HttpContentLength, reader.Size)

	io.Copy(w, reader)
}

func handleModifySince(r *http.Request) time.Time {
	lastStrs, ok := r.Header[HttpIfModifiedSince]
	if !ok || len(lastStrs) == 0 {
		return time.Time{}
	}
	last, err := http.ParseTime(lastStrs[0])
	if err != nil {
		log.WithError(err).Errorf("failed to parse %s header", HttpIfModifiedSince)
		return time.Time{}
	}
	return last
}

func setHeader(w http.ResponseWriter, key string, value any) {
	switch v := value.(type) {
	case string:
		if v != "" {
			w.Header().Add(key, v)
		}
	case int64:
		if v > 0 {
			w.Header().Add(key, strconv.FormatInt(v, 10))
		}
	case time.Time:
		if !v.IsZero() {
			w.Header().Add(key, v.UTC().Format(http.TimeFormat))
		}
	default:
		log.Errorf("unsupported type %v", reflect.TypeOf(value))
	}
}
