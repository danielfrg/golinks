package httpapi

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golinks/internal/links"
	"golinks/internal/store"
)

type legacyLink struct {
	Short, Long       string
	Created, LastEdit time.Time
	Clicks            int64
}

func legacy(link links.Link) legacyLink {
	return legacyLink{link.Short, link.Long, link.Created, link.LastEdit, link.Clicks}
}

func registerCompatibility(mux *http.ServeMux, s *store.Store, options Options) {
	mux.HandleFunc("GET /.export", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-ndjson")
		w.Header().Set("Content-Disposition", `attachment; filename="golinks.jsonl"`)
		if err := s.Export(r.Context(), w); err != nil {
			http.Error(w, "Could not export links", 500)
		}
	})
	mux.HandleFunc("GET /.export-stats", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", `attachment; filename="golinks-stats.csv"`)
		if err := s.ExportStats(r.Context(), w); err != nil {
			http.Error(w, "Could not export click counts", 500)
		}
	})
	mux.HandleFunc("GET /.detail/{short}", func(w http.ResponseWriter, r *http.Request) {
		link, err := s.Get(r.Context(), r.PathValue("short"))
		if err != nil {
			if err == store.ErrNotFound {
				http.NotFound(w, r)
			} else {
				http.Error(w, "Could not load link", 500)
			}
			return
		}
		if acceptHTML(r) {
			http.Redirect(w, r, "/?edit="+url.QueryEscape(link.Short), 302)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(legacy(link))
	})
	mux.HandleFunc("GET /.all", func(w http.ResponseWriter, r *http.Request) { http.Redirect(w, r, "/", 302) })
	mux.HandleFunc("GET /.search", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/?q="+url.QueryEscape(r.URL.Query().Get("q")), 302)
	})
	mux.HandleFunc("GET /.help", func(w http.ResponseWriter, r *http.Request) { http.Redirect(w, r, "/?help=1", 302) })
	mux.HandleFunc("GET /.opensearch", func(w http.ResponseWriter, r *http.Request) {
		type descriptor struct {
			XMLName     xml.Name `xml:"OpenSearchDescription"`
			XMLNS       string   `xml:"xmlns,attr"`
			Name        string   `xml:"ShortName"`
			Description string
			Encoding    string `xml:"InputEncoding"`
			URL         struct {
				Type     string `xml:"type,attr"`
				Template string `xml:"template,attr"`
			} `xml:"Url"`
		}
		d := descriptor{XMLNS: "http://a9.com/-/spec/opensearch/1.1/", Name: "Golinks", Description: "Homelab short links", Encoding: "UTF-8"}
		d.URL.Type = "text/html"
		d.URL.Template = strings.TrimSuffix(options.PublicURL, "/") + "/{searchTerms}"
		w.Header().Set("Content-Type", "application/opensearchdescription+xml")
		_, _ = w.Write([]byte(xml.Header))
		_ = xml.NewEncoder(w).Encode(d)
	})
}
