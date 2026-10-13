package httpapi

import (
	"errors"
	"io/fs"
	"log"
	"net/http"
	"strings"

	"golinks/internal/links"
	"golinks/internal/store"
	"golinks/internal/web"
)

func acceptHTML(r *http.Request) bool {
	return strings.Contains(strings.ToLower(r.Header.Get("Accept")), "text/html")
}

func serveRedirect(w http.ResponseWriter, r *http.Request, s *store.Store, assets fs.FS, options Options) {
	if r.Method != "GET" && r.Method != "HEAD" {
		w.Header().Set("Allow", "GET, HEAD")
		http.Error(w, "Method not allowed", 405)
		return
	}
	short, remainder, _ := strings.Cut(strings.TrimPrefix(r.URL.Path, "/"), "/")
	if strings.HasSuffix(short, "+") {
		w.Header().Set("Location", "/.detail/"+strings.TrimSuffix(short, "+"))
		w.WriteHeader(302)
		return
	}
	link, err := s.Get(r.Context(), short)
	if errors.Is(err, store.ErrNotFound) {
		trimmed := strings.TrimRight(short, ".,()[]{}")
		if trimmed != short {
			link, err = s.Get(r.Context(), trimmed)
		}
	}
	if errors.Is(err, store.ErrNotFound) {
		if acceptHTML(r) {
			web.Index(w, r, assets, 404)
		} else {
			http.NotFound(w, r)
		}
		return
	}
	if err != nil {
		log.Printf("resolve link: %v", err)
		http.Error(w, "Could not resolve link. Try again.", 500)
		return
	}
	query, err := r.URL.Query(), error(nil)
	destination, err := links.Expand(link.Long, links.Environment{Path: remainder, Query: query, Now: options.Now()})
	if err != nil {
		http.Error(w, err.Error(), 422)
		return
	}
	if r.Method == "GET" {
		if err := s.RecordClick(r.Context(), link.Short, options.Now()); err != nil {
			log.Printf("record click: %v", err)
		}
	}
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Location", destination.String())
	w.WriteHeader(302)
}
