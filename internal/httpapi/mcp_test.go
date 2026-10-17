package httpapi

import (
	"github.com/danielfrg/golinks/internal/store"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

func TestMCPEndpoint(t *testing.T) {
	db, err := store.Open(filepath.Join(t.TempDir(), "links.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	h := NewHandler(nil, db, Options{})
	for _, method := range []string{"GET", "HEAD", "PUT"} {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest(method, "http://localhost/_/mcp", nil))
		if w.Code != 405 || w.Header().Get("Cache-Control") != "no-store" {
			t.Fatalf("MCP routing %s: %d", method, w.Code)
		}
	}
}
