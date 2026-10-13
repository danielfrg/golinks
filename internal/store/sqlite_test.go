package store

import (
	"context"
	"errors"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func openTest(t *testing.T) *Store {
	t.Helper()
	s, err := Open(filepath.Join(t.TempDir(), "links.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { s.Close() })
	return s
}

func TestPersistenceAndConflicts(t *testing.T) {
	ctx := context.Background()
	filename := filepath.Join(t.TempDir(), "links.db")
	s, err := Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	list, err := s.List(ctx)
	if err != nil || list == nil || len(list) != 0 {
		t.Fatalf("empty: %v %v", list, err)
	}
	created, err := s.Create(ctx, "Meeting-Notes", "https://host/")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Create(ctx, "MEETINGNOTES", "https://other/"); !errors.Is(err, ErrConflict) {
		t.Fatal(err)
	}
	updated, err := s.Update(ctx, "meetingnotes", "https://new/")
	if err != nil || !updated.Created.Equal(created.Created) {
		t.Fatalf("%+v %v", updated, err)
	}
	var wg sync.WaitGroup
	for range 20 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := s.RecordClick(ctx, created.Short, time.Now()); err != nil {
				t.Error(err)
			}
		}()
	}
	wg.Wait()
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	s, err = Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	got, err := s.Get(ctx, "MEETINGNOTES")
	if err != nil || got.Long != "https://new/" || got.Clicks != 20 {
		t.Fatalf("%+v %v", got, err)
	}
	if err := s.Delete(ctx, got.Short); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Get(ctx, got.Short); !errors.Is(err, ErrNotFound) {
		t.Fatal(err)
	}
	var count int
	if err := s.db.QueryRow("SELECT count(*) FROM Stats").Scan(&count); err != nil || count != 0 {
		t.Fatalf("stats: %d %v", count, err)
	}
}

func TestDatabaseSettings(t *testing.T) {
	s := openTest(t)
	var mode string
	var timeout, syncMode int
	s.db.QueryRow("PRAGMA journal_mode").Scan(&mode)
	s.db.QueryRow("PRAGMA busy_timeout").Scan(&timeout)
	s.db.QueryRow("PRAGMA synchronous").Scan(&syncMode)
	if mode != "delete" || timeout != 5000 || syncMode != 2 {
		t.Fatalf("%s %d %d", mode, timeout, syncMode)
	}
	if _, err := Open(t.TempDir()); err == nil {
		t.Fatal("accepted directory as database")
	}
}
