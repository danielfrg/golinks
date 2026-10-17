package linkservice

import (
	"context"
	"errors"
	"net/url"
	"path/filepath"
	"testing"
	"time"

	"github.com/danielfrg/golinks/internal/store"
)

func TestLinkServiceContracts(t *testing.T) {
	ctx := context.Background()
	db, err := store.Open(filepath.Join(t.TempDir(), "links.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	s := New(db, Options{PublicURL: "https://go.example.net", Now: func() time.Time { return time.Date(2026, 10, 16, 12, 0, 0, 0, time.UTC) }})
	link, err := s.Create(ctx, "Meeting-Notes", "https://example.net/notes")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = s.Get(ctx, "meetingnotes"); err != nil {
		t.Fatal(err)
	}
	if _, err = s.Create(ctx, "meetingnotes", "https://example.net/other"); !errors.Is(err, store.ErrConflict) {
		t.Fatalf("collision: %v", err)
	}
	for _, destination := range []string{"", "javascript:alert(1)"} {
		if _, err = s.Update(ctx, link.Short, destination); err == nil {
			t.Fatal("invalid destination accepted")
		}
	}
	if _, err = s.Create(ctx, "_", "https://example.net"); err == nil {
		t.Fatal("reserved short accepted")
	}
	if _, err = s.Create(ctx, "alias", "/Meeting-Notes"); err != nil {
		t.Fatal(err)
	}
	resolved, err := s.Resolve(ctx, "alias", "", nil)
	if err != nil || resolved != "/Meeting-Notes" {
		t.Fatalf("alias %q %v", resolved, err)
	}
	if _, err = s.Create(ctx, "search", "https://example.net/?q={{QueryEscape .Path}}"); err != nil {
		t.Fatal(err)
	}
	resolved, err = s.Resolve(ctx, "search", "a/b c", url.Values{"page": {"2"}})
	if err != nil || resolved != "https://example.net/?page=2&q=a%2Fb+c" {
		t.Fatalf("template %q %v", resolved, err)
	}
	for _, short := range []string{"search", "alias"} {
		v, _ := s.Get(ctx, short)
		if v.Clicks != 0 {
			t.Fatal("resolution counted visit")
		}
	}
	ro := New(db, Options{ReadOnly: true})
	if _, err = ro.Create(ctx, "new", "https://example.net"); !errors.Is(err, ErrReadOnly) {
		t.Fatalf("readonly create: %v", err)
	}
	if _, err = ro.Update(ctx, link.Short, "https://example.net"); !errors.Is(err, ErrReadOnly) {
		t.Fatalf("readonly update: %v", err)
	}
	canceled, cancel := context.WithCancel(ctx)
	cancel()
	if _, err = s.Get(canceled, link.Short); err == nil {
		t.Fatal("cancellation ignored")
	}
}
