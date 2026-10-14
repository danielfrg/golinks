package main

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"testing/fstest"
	"time"

	"github.com/danielfrg/golinks/internal/httpapi"
	"github.com/danielfrg/golinks/internal/store"
)

func TestAppBackupRoundTrip(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	source, err := store.Open(filepath.Join(dir, "source.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer source.Close()
	for short, destination := range map[string]string{
		"Meeting-Notes": "https://example.org/notes?x=1&x=2",
		"search":        "https://example.org/search?q={{QueryEscape .Path}}",
		"alias":         "/meetingnotes",
	} {
		if _, err := source.Create(ctx, short, destination); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := source.Update(ctx, "meetingnotes", "https://example.org/current?x=1&x=2"); err != nil {
		t.Fatal(err)
	}
	for _, minute := range []int64{1760000000, 1760000060, 1760000060} {
		if err := source.RecordClick(ctx, "meetingnotes", time.Unix(minute, 0)); err != nil {
			t.Fatal(err)
		}
	}
	assets := fstest.MapFS{"index.html": {Data: []byte("Golinks")}}
	sourceHTTP := httpapi.NewHandler(assets, source, httpapi.Options{})
	paths := []string{"/MeetingNotes", "/search/space%20here", "/alias"}
	locations := make([]string, len(paths))
	for i, path := range paths {
		response := httptest.NewRecorder()
		sourceHTTP.ServeHTTP(response, httptest.NewRequest(http.MethodGet, path, nil))
		if response.Code != http.StatusFound {
			t.Fatal(response.Code, response.Body.String())
		}
		locations[i] = response.Header().Get("Location")
	}
	expected, err := source.List(ctx)
	if err != nil {
		t.Fatal(err)
	}
	exported := make(map[string][]byte)
	for path, name := range map[string]string{"/_/export": "links.jsonl", "/_/export-stats": "stats.csv"} {
		response := httptest.NewRecorder()
		sourceHTTP.ServeHTTP(response, httptest.NewRequest(http.MethodGet, path, nil))
		if response.Code != 200 || !strings.Contains(response.Header().Get("Content-Disposition"), "attachment") {
			t.Fatal(response.Code, response.Header(), response.Body.String())
		}
		exported[name] = bytes.Clone(response.Body.Bytes())
		if err := os.WriteFile(filepath.Join(dir, name), exported[name], 0600); err != nil {
			t.Fatal(err)
		}
	}
	for _, withStats := range []bool{false, true} {
		t.Run(map[bool]string{false: "links only", true: "with history"}[withStats], func(t *testing.T) {
			database := filepath.Join(t.TempDir(), "restored.db")
			args := []string{"--links", filepath.Join(dir, "links.jsonl")}
			if withStats {
				args = append(args, "--stats", filepath.Join(dir, "stats.csv"))
			}
			var output bytes.Buffer
			if err := runImport(ctx, args, database, &output); err != nil || !strings.Contains(output.String(), "Imported 3 links; skipped 0") {
				t.Fatal(output.String(), err)
			}
			output.Reset()
			if err := runImport(ctx, args, database, &output); err != nil || !strings.Contains(output.String(), "Imported 0 links; skipped 3") {
				t.Fatal(output.String(), err)
			}
			restored, err := store.Open(database)
			if err != nil {
				t.Fatal(err)
			}
			defer restored.Close()
			actual, err := restored.List(ctx)
			if err != nil || !reflect.DeepEqual(actual, expected) {
				t.Fatalf("restored %+v; expected %+v: %v", actual, expected, err)
			}
			if withStats {
				var history bytes.Buffer
				if err := restored.ExportStats(ctx, &history); err != nil || !bytes.Equal(history.Bytes(), exported["stats.csv"]) {
					t.Fatalf("history changed: %s: %v", history.String(), err)
				}
			}
			handler := httpapi.NewHandler(assets, restored, httpapi.Options{})
			for i, path := range paths {
				response := httptest.NewRecorder()
				handler.ServeHTTP(response, httptest.NewRequest(http.MethodHead, path, nil))
				if response.Code != http.StatusFound || response.Header().Get("Location") != locations[i] {
					t.Fatalf("restored %s: %d %s", path, response.Code, response.Header().Get("Location"))
				}
			}
		})
	}
}
