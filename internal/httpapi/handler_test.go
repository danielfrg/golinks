package httpapi

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

	"golinks/internal/store"
)

func testApp(t *testing.T, readOnly bool) (http.Handler, *store.Store) {
	t.Helper()
	s, err := store.Open(filepath.Join(t.TempDir(), "app.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { s.Close() })
	return NewHandler(fstest.MapFS{"index.html": {Data: []byte("<main>Golinks</main>")}, "assets/app.js": {Data: []byte("console.log('app')")}}, s, Options{PublicURL: "http://go.test", ReadOnly: readOnly}), s
}
func request(h http.Handler, method, target, body string) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, target, strings.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("X-Golink", "1")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	return w
}
func TestManagement(t *testing.T) {
	h, _ := testApp(t, false)
	for _, tc := range []struct {
		method, path, body string
		status             int
	}{
		{"GET", "/_/api/links", "", 200},
		{"POST", "/_/api/links", `{"short":"Meeting-Notes","long":"https://host/"}`, 201},
		{"POST", "/_/api/links", `{"short":"meetingnotes","long":"https://other/"}`, 409},
		{"POST", "/_/api/links", `{"short":"a/b","long":"https://host/"}`, 422},
		{"PUT", "/_/api/links/meetingnotes", `{"long":"https://new/"}`, 200},
		{"GET", "/_/api/links/meetingnotes", "", 200},
		{"DELETE", "/_/api/links/meetingnotes", "", 204},
		{"GET", "/_/api/links/meetingnotes", "", 404},
	} {
		if w := request(h, tc.method, tc.path, tc.body); w.Code != tc.status {
			t.Fatalf("%s %s: %d %s", tc.method, tc.path, w.Code, w.Body.String())
		}
	}
	w := request(h, "GET", "/_/api/links", "")
	if !strings.Contains(w.Body.String(), `"links":[]`) || strings.Contains(w.Body.String(), "$schema") {
		t.Fatal(w.Body.String())
	}
	if w := request(h, "GET", "/_/api/openapi.json", ""); w.Code != 200 {
		t.Fatal(w.Code)
	}
}

func TestAPINamespace(t *testing.T) {
	h, _ := testApp(t, false)
	if w := request(h, "GET", "/_/api/docs", ""); w.Code != 200 || !strings.Contains(w.Body.String(), "/_/api/openapi") {
		t.Fatalf("docs: %d %s", w.Code, w.Body)
	}
	for _, target := range []string{"/_", "/_/nope", "/_/api", "/_/api/nope", "/_/api%2fnope", "/_/apiary", "/.api/links"} {
		if w := request(h, "GET", target, ""); w.Code != 404 {
			t.Errorf("%s: %d", target, w.Code)
		}
	}
	for _, name := range []string{"_", "_-", "_--"} {
		w := request(h, "POST", "/_/api/links", `{"short":"`+name+`","long":"https://host/"}`)
		if w.Code != 422 || !strings.Contains(w.Body.String(), "reserved") {
			t.Errorf("reserved %s: %d %s", name, w.Code, w.Body)
		}
	}
	if w := request(h, "POST", "/_/api/links", `{"short":"api","long":"https://host/"}`); w.Code != 201 {
		t.Fatalf("api short name unavailable: %d %s", w.Code, w.Body)
	}
	for _, target := range []string{"/api", "/API", "/a-p-i", "/api/docs", "/api/links"} {
		w := request(h, "GET", target, "")
		want := "https://host/"
		if target == "/api/docs" {
			want += "docs"
		}
		if target == "/api/links" {
			want += "links"
		}
		if w.Code != 302 || w.Header().Get("Location") != want {
			t.Errorf("%s should be a shortcut: %d %s", target, w.Code, w.Header().Get("Location"))
		}
	}
}

func TestLinkResponsesHaveNoIdentity(t *testing.T) {
	h, s := testApp(t, false)
	snapshot := `{"short":"wiki","long":"https://host/","created":"2026-10-13T00:00:00Z","lastEdit":"2026-10-13T00:00:00Z","clicks":0}`
	if _, err := s.Import(context.Background(), strings.NewReader(snapshot), nil); err != nil {
		t.Fatal(err)
	}
	for _, target := range []string{"/_/api/links", "/_/api/links/wiki", "/.detail/wiki", "/_/api/openapi.json"} {
		w := request(h, "GET", target, "")
		if w.Code != 200 || strings.Contains(strings.ToLower(w.Body.String()), `"owner"`) {
			t.Fatalf("identity in %s: %d %s", target, w.Code, w.Body)
		}
	}
}

func TestBrowserImport(t *testing.T) {
	h, s := testApp(t, false)
	snapshot := `{"short":"wiki","long":"https://host/","created":"2026-10-13T00:00:00Z","lastEdit":"2026-10-13T00:00:00Z","clicks":7}`
	body, _ := json.Marshal(map[string]string{"links": snapshot})
	for _, expected := range []string{`"imported":1`, `"skipped":1`} {
		w := request(h, "POST", "/_/api/import", string(body))
		if w.Code != 200 || !strings.Contains(w.Body.String(), expected) {
			t.Fatal(w.Code, w.Body.String())
		}
	}
	link, err := s.Get(context.Background(), "wiki")
	if err != nil || link.Clicks != 7 {
		t.Fatal(link, err)
	}
	bad, _ := json.Marshal(map[string]string{"links": strings.ReplaceAll(snapshot, "wiki", "docs") + "\n{broken"})
	if w := request(h, "POST", "/_/api/import", string(bad)); w.Code != 422 {
		t.Fatal(w.Code, w.Body.String())
	}
	list, _ := s.List(context.Background())
	if len(list) != 1 {
		t.Fatal("invalid import wrote partial data", list)
	}
	withStats, _ := json.Marshal(map[string]string{"links": strings.ReplaceAll(snapshot, "wiki", "history"), "stats": "history,1760313600,11\n"})
	if w := request(h, "POST", "/_/api/import", string(withStats)); w.Code != 200 {
		t.Fatal(w.Code, w.Body.String())
	}
	link, _ = s.Get(context.Background(), "history")
	if link.Clicks != 11 {
		t.Fatal(link.Clicks)
	}
	large, _ := json.Marshal(map[string]string{"links": `{"short":"large","long":"https://host/?q=` + strings.Repeat("&", 40000) + `"}`})
	if w := request(h, "POST", "/_/api/import", string(large)); w.Code != 200 {
		t.Fatal("export larger than 64 KiB failed", w.Code, w.Body.String())
	}
	readOnly, _ := testApp(t, true)
	if w := request(readOnly, "POST", "/_/api/import", string(body)); w.Code != 405 {
		t.Fatal(w.Code)
	}
	r := httptest.NewRequest("POST", "/_/api/import", strings.NewReader(string(body)))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != 403 {
		t.Fatal("import missing mutation guard", w.Code)
	}
	r = httptest.NewRequest("POST", "/_/api/import", strings.NewReader(string(body)))
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("X-Golink", "1")
	r.Header.Set("Origin", "https://other.test")
	w = httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != 403 {
		t.Fatal("cross-origin import accepted", w.Code)
	}
	if w := request(h, "POST", "/_/api/import", strings.Repeat(" ", int(maxImportBody)+1)); w.Code != 413 {
		t.Fatal("oversized import accepted", w.Code)
	}
	if w := request(h, "POST", "/_/api/import", `{"links":""}`); w.Code != 200 {
		t.Fatal("empty export cannot be imported", w.Code, w.Body.String())
	}
}
func TestRedirectsAndPaths(t *testing.T) {
	h, s := testApp(t, false)
	ctx := context.Background()
	s.Create(ctx, "meeting-notes", "https://host/base%2f/{{.Path}}?a=1")
	s.Create(ctx, "alias", "/meetingnotes")
	for _, name := range []string{"apiary", "api-client", "assets", "healthz"} {
		s.Create(ctx, name, "https://host/")
	}
	for _, tc := range []struct{ path, want string }{
		{"/MeetingNotes/a//b?x=2", "https://host/base%2f/a//b?a=1&x=2"},
		{"/meetingnotes/a/../b", "https://host/base%2f/a/../b?a=1"},
		{"/meetingnotes/a%2Fb", "https://host/base%2f/a/b?a=1"},
		{"/meetingnotes.", "https://host/base%2f/?a=1"},
		{"/alias", "/meetingnotes"}, {"/apiary", "https://host/"}, {"/api-client", "https://host/"}, {"/assets", "https://host/"}, {"/healthz", "https://host/"},
	} {
		w := request(h, "GET", tc.path, "")
		if w.Code != 302 || w.Header().Get("Location") != tc.want {
			t.Errorf("%s: %d %q want %q", tc.path, w.Code, w.Header().Get("Location"), tc.want)
		}
	}
	before, _ := s.Get(ctx, "meetingnotes")
	w := request(h, "HEAD", "/meetingnotes", "")
	after, _ := s.Get(ctx, "meetingnotes")
	if w.Code != 302 || before.Clicks != after.Clicks {
		t.Fatal("HEAD counted")
	}
	w = request(h, "GET", "/meetingnotes+", "")
	if w.Code != 302 || w.Header().Get("Location") != "/.detail/meetingnotes" {
		t.Fatal(w.Code, w.Header())
	}
	w = request(h, "GET", "/.detail/meetingnotes", "")
	var detail map[string]any
	if json.Unmarshal(w.Body.Bytes(), &detail) != nil || detail["Short"] != "meeting-notes" {
		t.Fatal(w.Body.String())
	}
	for _, target := range []string{"/_/api/nope", "/.assets/missing.js", "/.assets/.keep", "/_/api%2fnope"} {
		if w := request(h, "GET", target, ""); w.Code != 404 {
			t.Errorf("%s: %d", target, w.Code)
		}
	}
	w = request(h, "GET", "/missing", "")
	if w.Code != 404 {
		t.Fatal(w.Code)
	}
	r := httptest.NewRequest("GET", "/missing", nil)
	r.Header.Set("Accept", "text/html")
	w = httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != 404 || !strings.Contains(w.Body.String(), "Golinks") {
		t.Fatal(w.Code, w.Body)
	}
	w = request(h, "GET", "/.assets/assets/app.js", "")
	if w.Code != 200 {
		t.Fatal(w.Code)
	}
}
func TestMutationGuardsAndPreview(t *testing.T) {
	h, s := testApp(t, false)
	r := httptest.NewRequest("POST", "/_/api/links", strings.NewReader(`{"short":"a","long":"https://host/"}`))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != 403 {
		t.Fatal(w.Code)
	}
	r.Header.Set("X-Golink", "1")
	r.Header.Set("Origin", "https://evil.test")
	w = httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != 403 {
		t.Fatal(w.Code)
	}
	w = request(h, "POST", "/_/api/links", `{"short":"a","long":"https://host/`+strings.Repeat("a", 65536)+`"}`)
	if w.Code != 413 {
		t.Fatal(w.Code)
	}
	w = request(h, "POST", "/_/api/resolve", `{"long":"https://host/?q={{QueryEscape .Path}}","path":"a b","query":"x=1"}`)
	if w.Code != 200 || !strings.Contains(w.Body.String(), "a+b") {
		t.Fatal(w.Code, w.Body)
	}
	s.Create(context.Background(), "bad", "{{.Path}}:alert(1)")
	w = request(h, "GET", "/bad/javascript", "")
	if w.Code != 422 || w.Header().Get("Location") != "" {
		t.Fatal(w.Code, w.Header())
	}
	ro, _ := testApp(t, true)
	w = request(ro, "POST", "/_/api/links", `{"short":"a","long":"https://host/"}`)
	if w.Code != 405 {
		t.Fatal(w.Code)
	}
	for _, target := range []string{"/", "/.healthz", "/.readyz", "/_/api/settings"} {
		w = request(h, "GET", target, "")
		if w.Code != 200 {
			t.Fatalf("%s %d %s", target, w.Code, w.Body)
		}
	}
	_ = io.Discard
}
