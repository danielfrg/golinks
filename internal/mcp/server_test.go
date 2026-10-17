package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/danielfrg/golinks/internal/linkservice"
	"github.com/danielfrg/golinks/internal/store"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

func openSession(t *testing.T, ro bool) (*sdk.ClientSession, *store.Store) {
	t.Helper()
	db, err := store.Open(filepath.Join(t.TempDir(), "links.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	server := httptest.NewServer(NewHandler(linkservice.New(db, linkservice.Options{PublicURL: "https://go.example.net", ReadOnly: ro}), nil))
	t.Cleanup(server.Close)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	t.Cleanup(cancel)
	session, err := sdk.NewClient(&sdk.Implementation{Name: "test", Version: "1"}, nil).Connect(ctx, &sdk.StreamableClientTransport{Endpoint: server.URL, DisableStandaloneSSE: true}, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { session.Close() })
	return session, db
}
func call(t *testing.T, s *sdk.ClientSession, name string, args any, wantError bool) map[string]any {
	t.Helper()
	v, err := s.CallTool(context.Background(), &sdk.CallToolParams{Name: name, Arguments: args})
	if err != nil {
		t.Fatal(err)
	}
	if v.IsError != wantError {
		t.Fatalf("%s error=%v content=%v", name, v.IsError, v.Content)
	}
	if wantError {
		return nil
	}
	data, _ := json.Marshal(v.StructuredContent)
	var out map[string]any
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatal(err)
	}
	return out
}
func TestToolsAndContracts(t *testing.T) {
	s, db := openSession(t, false)
	tools, err := s.ListTools(context.Background(), nil)
	if err != nil || len(tools.Tools) != 5 {
		t.Fatalf("tools: %v %v", tools, err)
	}
	call(t, s, "create_link", map[string]any{"short": "Docs-Home", "destination": "https://example.net/docs"}, false)
	call(t, s, "create_link", map[string]any{"short": "docshome", "destination": "https://example.net/other"}, true)
	call(t, s, "create_link", map[string]any{"short": "bad", "destination": "javascript:alert(1)"}, true)
	got := call(t, s, "get_link", map[string]any{"short": "docshome"}, false)
	if got["public_url"] != "https://go.example.net/Docs-Home" {
		t.Fatalf("URL: %v", got)
	}
	call(t, s, "update_link", map[string]any{"short": "DocsHome", "destination": "https://example.net/new"}, false)
	page := call(t, s, "search_links", map[string]any{"query": "NEW", "limit": 1}, false)
	if page["total"] != float64(1) {
		t.Fatalf("page: %v", page)
	}
	for _, n := range []int{0, -1, 101} {
		call(t, s, "search_links", map[string]any{"limit": n}, true)
	}
	call(t, s, "search_links", map[string]any{"offset": -1}, true)
	call(t, s, "get_link", map[string]any{"short": "missing"}, true)
	call(t, s, "create_link", map[string]any{"short": "alias", "destination": "/Docs-Home"}, false)
	resolved := call(t, s, "resolve_link", map[string]any{"short": "alias"}, false)
	if resolved["url"] != "/Docs-Home" {
		t.Fatal(resolved)
	}
	call(t, s, "resolve_link", map[string]any{"short": "alias", "query": "%zz"}, true)
	v, _ := db.Get(context.Background(), "alias")
	if v.Clicks != 0 {
		t.Fatal("resolve counted click")
	}
	ro, _ := openSession(t, true)
	listed, _ := ro.ListTools(context.Background(), nil)
	if len(listed.Tools) != 3 {
		t.Fatal("readonly tools")
	}
	if _, err := ro.CallTool(context.Background(), &sdk.CallToolParams{Name: "create_link", Arguments: map[string]any{"short": "x", "destination": "https://example.net"}}); err == nil {
		t.Fatal("readonly write tool exists")
	}
}

func TestPagingAndSanitizedFailures(t *testing.T) {
	s, db := openSession(t, false)
	for i := 0; i < 23; i++ {
		if _, err := db.Create(context.Background(), fmt.Sprintf("link%02d", i), "https://example.net"); err != nil {
			t.Fatal(err)
		}
	}
	first := call(t, s, "search_links", map[string]any{}, false)
	if first["total"] != float64(23) || len(first["links"].([]any)) != 20 || first["next_offset"] != float64(20) {
		t.Fatal(first)
	}
	last := call(t, s, "search_links", map[string]any{"offset": 20, "sort": "clicks"}, false)
	if len(last["links"].([]any)) != 3 || last["next_offset"] != nil {
		t.Fatal(last)
	}
	empty := call(t, s, "search_links", map[string]any{"offset": 1000000}, false)
	if len(empty["links"].([]any)) != 0 {
		t.Fatal(empty)
	}
	call(t, s, "search_links", map[string]any{"sort": "bad"}, true)
	db.Close()
	v, err := s.CallTool(context.Background(), &sdk.CallToolParams{Name: "get_link", Arguments: map[string]any{"short": "link00"}})
	data, _ := json.Marshal(v.Content)
	if err != nil || !v.IsError || !strings.Contains(string(data), "Links unavailable") || strings.Contains(string(data), "database") {
		t.Fatalf("unsafe failure %v %v", v, err)
	}
}
