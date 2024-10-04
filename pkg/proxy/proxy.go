package proxy

import (
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/vietanhduong/bucket-proxy/pkg/bucket"
	"github.com/vietanhduong/bucket-proxy/pkg/bucket/types"
	"github.com/vietanhduong/bucket-proxy/pkg/logging"
)

var log = logging.WithField("pkg", "pkg/proxy")

type Proxy struct {
	bucket       bucket.Interface
	webMode      bool
	indexPage    string
	notFoundPage string
}

func New(bucket bucket.Interface, opt ...Option) *Proxy {
	instance := &Proxy{
		bucket:       bucket,
		indexPage:    "index.html",
		notFoundPage: "404.html",
	}
	for _, o := range opt {
		o(instance)
	}
	log.WithFields(logrus.Fields{
		"web_mode":       instance.webMode,
		"index_page":     instance.indexPage,
		"not_found_page": instance.notFoundPage,
	}).Info("init proxy instance with values")
	return instance
}

func (p *Proxy) HttpHandler() (string, http.Handler) {
	return "/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { p.proxy(w, r) })
}

func (p *Proxy) proxy(w http.ResponseWriter, r *http.Request) {
	l := log.WithField("path", r.URL.Path)
	path := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/"), "/")
	if p.webMode && path == "" {
		p.handleIndexPage(w, r, "")
		return
	}
	attrs, err := p.bucket.ObjectMetadata(r.Context(), path)
	if err != nil {
		l.WithError(err).Error("could not get object metadata")
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if attrs == nil {
		p.handleNotFoundPage(w, r)
		return
	}

	if attrs.IsDirectory {
		p.handleIndexPage(w, r, path)
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

func (p *Proxy) handleNotFoundPage(w http.ResponseWriter, r *http.Request) {
	if p.webMode {
		// try to download the not found page, response a not found message in-case download failed
		reader, err := p.bucket.Download(r.Context(), p.notFoundPage, types.DownloadOptions{})
		if err == nil {
			setHeader(w, HttpLastModified, reader.Updated)
			setHeader(w, HttpContentType, reader.ContentType)
			setHeader(w, HttpContentLanguage, reader.ContentLanguage)
			setHeader(w, HttpCacheControl, reader.CacheControl)
			setHeader(w, HttpContentEncoding, reader.ContentEncoding)
			setHeader(w, HttpContentDisposition, reader.ContentDisposition)
			setHeader(w, HttpContentLength, reader.Size)
			io.Copy(w, reader)
			return
		}
		log.WithError(err).Trace("failed to download not found page")
	}
	http.Error(w, "not found", http.StatusNotFound)
}

func (p *Proxy) handleIndexPage(w http.ResponseWriter, r *http.Request, dir string) {
	if p.webMode {
		path := p.indexPage
		if dir != "" {
			path = dir + "/" + p.indexPage
		}
		reader, err := p.bucket.Download(r.Context(), path, types.DownloadOptions{})
		if err == nil {
			setHeader(w, HttpLastModified, reader.Updated)
			setHeader(w, HttpContentType, reader.ContentType)
			setHeader(w, HttpContentLanguage, reader.ContentLanguage)
			setHeader(w, HttpCacheControl, reader.CacheControl)
			setHeader(w, HttpContentEncoding, reader.ContentEncoding)
			setHeader(w, HttpContentDisposition, reader.ContentDisposition)
			setHeader(w, HttpContentLength, reader.Size)
			io.Copy(w, reader)
			return
		}
		// if go to here, it means the index page is not found
		// we should return 404
		log.WithField("path", path).WithError(err).Trace("failed to download index page")
	}
	p.handleNotFoundPage(w, r)
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
