package store

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"
	"time"
)

func oldDatabase(t *testing.T, clicks string) string {
	t.Helper()
	filename := filepath.Join(t.TempDir(), "golink.db")
	db, err := sql.Open("sqlite", filename)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	_, err = db.Exec(`CREATE TABLE Links(ID TEXT PRIMARY KEY,Short TEXT NOT NULL,Long TEXT NOT NULL,Created INTEGER NOT NULL,LastEdit INTEGER NOT NULL,Owner TEXT NOT NULL DEFAULT '');
CREATE TABLE Stats(ID TEXT NOT NULL,Created INTEGER NOT NULL,Clicks INTEGER);
INSERT INTO Links VALUES('wiki','Wiki','https://example.org/',1654294500,1654294560,'');
INSERT INTO Stats VALUES('wiki',1654294500,3),('wiki',1654294500,4),('wiki',1654294560,NULL),('deleted',1654294500,5);`)
	if err != nil {
		t.Fatal(err)
	}
	if clicks != "" {
		if _, err := db.Exec("INSERT INTO Stats VALUES('wiki',1654294620,?)", clicks); err != nil {
			t.Fatal(err)
		}
	}
	return filename
}

func TestExistingDatabaseUpgrade(t *testing.T) {
	filename := oldDatabase(t, "")
	s, err := Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	got, err := s.Get(ctx, "wiki")
	if err != nil || got.Clicks != 7 || got.Created.Unix() != 1654294500 || got.LastEdit.Unix() != 1654294560 {
		t.Fatal(got, err)
	}
	if err := s.RecordClick(ctx, "wiki", time.Unix(1654294500, 0)); err != nil {
		t.Fatalf("click after upgrade: %v", err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	s, err = Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	got, err = s.Get(ctx, "wiki")
	if err != nil || got.Clicks != 8 {
		t.Fatal(got, err)
	}
	var archived int64
	if err := s.db.QueryRow("SELECT SUM(Clicks) FROM StatsOrphans_v1").Scan(&archived); err != nil || archived != 5 {
		t.Fatal("orphan history was lost", archived, err)
	}
	if err := s.Delete(ctx, "wiki"); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Create(ctx, "wiki", "https://example.org/new"); err != nil {
		t.Fatal(err)
	}
	got, _ = s.Get(ctx, "wiki")
	if got.Clicks != 0 {
		t.Fatal("deleted history returned", got.Clicks)
	}
}

func TestDatabaseUpgradeFailureLeavesDataIntact(t *testing.T) {
	filename := oldDatabase(t, "invalid")
	if s, err := Open(filename); err == nil {
		s.Close()
		t.Fatal("accepted invalid click data")
	}
	db, err := sql.Open("sqlite", filename)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM Stats").Scan(&count); err != nil || count != 5 {
		t.Fatal(count, err)
	}
}
