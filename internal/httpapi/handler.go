package httpapi

import (
	"bytes"
	"io"
	"io/fs"
	"mime"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/danielfrg/golinks/internal/store"
	"github.com/danielfrg/golinks/internal/web"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
)

type Options struct {
	PublicURL string
	ReadOnly  bool
	Now       func() time.Time
}

func NewHandler(assets fs.FS, s *store.Store, options Options) http.Handler {
	if options.Now == nil {
		options.Now = time.Now
	}
	mux := http.NewServeMux()
	config := huma.DefaultConfig("Golinks API", "0.1.0")
	config.OpenAPIPath = "/_/api/openapi"
	config.DocsPath = "/_/api/docs"
	config.SchemasPath = ""
	config.CreateHooks = nil
	api := humago.New(mux, config)
	registerLinks(api, s, options)
	registerImport(api, s)
	mux.HandleFunc("GET /_/healthz", func(w http.ResponseWriter, r *http.Request) { _, _ = io.WriteString(w, "ok\n") })
	mux.HandleFunc("GET /_/readyz", func(w http.ResponseWriter, r *http.Request) {
		if err := s.Ping(r.Context()); err != nil {
			http.Error(w, "Database unavailable", 503)
			return
		}
		_, _ = io.WriteString(w, "ok\n")
	})
	mux.Handle("GET /_/assets/", web.Handler(assets))
	registerCompatibility(mux, s, options)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Referrer-Policy", "no-referrer")
		apiRequest := r.URL.Path == "/_/api" || strings.HasPrefix(r.URL.Path, "/_/api/")
		if apiRequest && r.Method != "GET" && r.Method != "HEAD" && r.Method != "OPTIONS" {
			if options.ReadOnly && r.URL.Path != "/_/api/resolve" {
				http.Error(w, "Golinks is read-only", 405)
				return
			}
			if r.Header.Get("X-Golink") != "1" {
				http.Error(w, "X-Golink header is required", 403)
				return
			}
			if origin := r.Header.Get("Origin"); origin != "" {
				u, err := url.Parse(origin)
				if err != nil || u.Host != r.Host || (u.Scheme != "http" && u.Scheme != "https") || u.User != nil || u.RawQuery != "" || u.Fragment != "" || (u.Path != "" && u.Path != "/") {
					http.Error(w, "Origin does not match the request host", 403)
					return
				}
			}
			media, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
			if err != nil || media != "application/json" {
				http.Error(w, "Use application/json", 415)
				return
			}
			limit := int64(64 * 1024)
			limitLabel := "64 KiB"
			if r.URL.Path == "/_/api/import" && r.Method == http.MethodPost {
				limit = maxImportBody
				limitLabel = "16 MiB; use the command-line import for larger backups"
			}
			data, err := io.ReadAll(http.MaxBytesReader(w, r.Body, limit))
			if err != nil {
				http.Error(w, "Request exceeds "+limitLabel, 413)
				return
			}
			r.Body = io.NopCloser(bytes.NewReader(data))
		}
		if r.URL.Path == "/" {
			if r.Method != "GET" && r.Method != "HEAD" {
				w.Header().Set("Allow", "GET, HEAD")
				http.Error(w, "Method not allowed", 405)
				return
			}
			web.Index(w, r, assets, 200)
			return
		}
		if r.URL.Path == "/_" || strings.HasPrefix(r.URL.Path, "/_/") {
			mux.ServeHTTP(w, r)
			return
		}
		serveRedirect(w, r, s, assets, options)
	})
}
