package store

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"
)

const sampleSnapshot = `{"short":"Meeting-Notes","long":"https://host/","created":"2022-06-03T22:15:29Z","lastEdit":"2022-06-04T22:15:29Z","clicks":7}`

func TestSnapshotsRoundTripAndIdempotence(t *testing.T) {
	s := openTest(t)
	ctx := context.Background()
	stats := "meetingnotes,1654294500,3\nmeetingnotes,1654294560,5\n"
	result, err := s.Import(ctx, strings.NewReader(sampleSnapshot), strings.NewReader(stats))
	if err != nil || result.Imported != 1 {
		t.Fatal(result, err)
	}
	link, err := s.Get(ctx, "meetingnotes")
	if err != nil || link.Clicks != 8 || link.Created.Unix() != 1654294529 {
		t.Fatal(link, err)
	}
	result, err = s.Import(ctx, strings.NewReader(sampleSnapshot), strings.NewReader(stats))
	if err != nil || result.Skipped != 1 {
		t.Fatal(result, err)
	}
	link, _ = s.Get(ctx, "meetingnotes")
	if link.Clicks != 8 {
		t.Fatal(link.Clicks)
	}
	var snapshot, csv bytes.Buffer
	if err := s.Export(ctx, &snapshot); err != nil {
		t.Fatal(err)
	}
	if err := s.ExportStats(ctx, &csv); err != nil {
		t.Fatal(err)
	}
	other := openTest(t)
	if _, err := other.Import(ctx, &snapshot, &csv); err != nil {
		t.Fatal(err)
	}
	got, _ := other.Get(ctx, link.Short)
	if got != link {
		t.Fatal(got, link)
	}
	updated, err := s.Update(ctx, link.Short, "https://new/")
	if err != nil || updated.Created != link.Created {
		t.Fatal(updated, err)
	}
}

func TestExportUsesApplicationFields(t *testing.T) {
	ctx := context.Background()
	s := openTest(t)
	if _, err := s.Create(ctx, "wiki", "https://example.org/"); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	if err := s.Export(ctx, &out); err != nil {
		t.Fatal(err)
	}
	var record map[string]json.RawMessage
	if err := json.Unmarshal(bytes.TrimSpace(out.Bytes()), &record); err != nil {
		t.Fatal(err)
	}
	if len(record) != 5 {
		t.Fatalf("expected only link fields, got %s", out.Bytes())
	}
	for _, field := range []string{"short", "long", "created", "lastEdit", "clicks"} {
		if _, ok := record[field]; !ok {
			t.Errorf("missing %s in %s", field, out.Bytes())
		}
	}
}

func TestLinksOnlyExportImportPreservesData(t *testing.T) {
	ctx := context.Background()
	source := openTest(t)
	if _, err := source.Import(ctx, strings.NewReader(sampleSnapshot), nil); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	if err := source.Export(ctx, &out); err != nil {
		t.Fatal(err)
	}
	target := openTest(t)
	if _, err := target.Import(ctx, bytes.NewReader(out.Bytes()), nil); err != nil {
		t.Fatal(err)
	}
	expected, _ := source.Get(ctx, "meetingnotes")
	actual, _ := target.Get(ctx, "meetingnotes")
	if actual != expected {
		t.Fatalf("restored %+v, want %+v", actual, expected)
	}
	if result, err := target.Import(ctx, bytes.NewReader(out.Bytes()), nil); err != nil || result.Imported != 0 || result.Skipped != 1 {
		t.Fatal(result, err)
	}
	actual, _ = target.Get(ctx, "meetingnotes")
	if actual.Clicks != expected.Clicks {
		t.Fatal("repeat import duplicated clicks")
	}
}
func TestSnapshotFailuresAreAtomic(t *testing.T) {
	for _, input := range []string{sampleSnapshot + "\n{broken", sampleSnapshot + "\n" + strings.ReplaceAll(sampleSnapshot, "Meeting-Notes", "meetingnotes"), strings.ReplaceAll(sampleSnapshot, "https://host/", "javascript:alert(1)")} {
		s := openTest(t)
		if _, err := s.Import(context.Background(), strings.NewReader(input), nil); err == nil {
			t.Fatal("accepted invalid snapshot")
		}
		list, _ := s.List(context.Background())
		if len(list) != 0 {
			t.Fatal("partial import")
		}
	}
	s := openTest(t)
	if _, err := s.Import(context.Background(), strings.NewReader(sampleSnapshot), strings.NewReader("meetingnotes,not-a-number,4")); err == nil {
		t.Fatal("accepted invalid CSV")
	}
	list, _ := s.List(context.Background())
	if len(list) != 0 {
		t.Fatal("partial CSV import")
	}
	if _, err := s.Import(context.Background(), strings.NewReader(sampleSnapshot), nil); err != nil {
		t.Fatal(err)
	}
	link, _ := s.Get(context.Background(), "meetingnotes")
	if link.Clicks != 7 {
		t.Fatal(link.Clicks)
	}
}

func TestSnapshotRoundTripWithJSONEscaping(t *testing.T) {
	ctx := context.Background()
	s := openTest(t)
	destination := "https://host.test/?q=" + strings.Repeat("&", 40000)
	if _, err := s.Create(ctx, "large", destination); err != nil {
		t.Fatal(err)
	}
	var snapshot bytes.Buffer
	if err := s.Export(ctx, &snapshot); err != nil {
		t.Fatal(err)
	}
	other := openTest(t)
	if _, err := other.Import(ctx, &snapshot, nil); err != nil {
		t.Fatal(err)
	}
	got, err := other.Get(ctx, "large")
	if err != nil || got.Long != destination {
		t.Fatal("destination did not round trip", err)
	}
}

func TestImportSkipsOrphanStatsAndValidatesThem(t *testing.T) {
	ctx := context.Background()
	s := openTest(t)
	result, err := s.Import(ctx, strings.NewReader(sampleSnapshot), strings.NewReader("meetingnotes,1654294500,3\ndeleted,1654294500,7\n"))
	if err != nil || result.Imported != 1 || result.SkippedStats != 1 {
		t.Fatal(result, err)
	}
	got, _ := s.Get(ctx, "meetingnotes")
	if got.Clicks != 3 {
		t.Fatal(got.Clicks)
	}
	other := openTest(t)
	if _, err := other.Import(ctx, strings.NewReader(sampleSnapshot), strings.NewReader("deleted,not-a-time,7\n")); err == nil {
		t.Fatal("malformed orphan CSV must still fail")
	}
}
