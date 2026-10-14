package store

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/danielfrg/golinks/internal/links"
)

var ErrInvalidSnapshot = errors.New("Invalid backup")

func invalidSnapshot(format string, args ...any) error {
	return fmt.Errorf("%w: %s", ErrInvalidSnapshot, fmt.Sprintf(format, args...))
}

type ImportResult struct {
	Imported     int `json:"imported"`
	Skipped      int `json:"skipped"`
	SkippedStats int `json:"skippedStats"`
}
type stat struct {
	ID              string
	Created, Clicks int64
}

func (s *Store) Export(ctx context.Context, w io.Writer) error {
	all, err := s.List(ctx)
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(w)
	for _, link := range all {
		if err := encoder.Encode(link); err != nil {
			return err
		}
	}
	return nil
}
func (s *Store) ExportStats(ctx context.Context, w io.Writer) error {
	rows, err := s.db.QueryContext(ctx, "SELECT ID,Created,Clicks FROM Stats ORDER BY Created,ID")
	if err != nil {
		return err
	}
	defer rows.Close()
	csvOut := csv.NewWriter(w)
	for rows.Next() {
		var value stat
		if err := rows.Scan(&value.ID, &value.Created, &value.Clicks); err != nil {
			return err
		}
		if err := csvOut.Write([]string{value.ID, strconv.FormatInt(value.Created, 10), strconv.FormatInt(value.Clicks, 10)}); err != nil {
			return err
		}
	}
	if err := rows.Err(); err != nil {
		return err
	}
	csvOut.Flush()
	return csvOut.Error()
}
func (s *Store) Import(ctx context.Context, linksJSONL io.Reader, statsCSV io.Reader) (ImportResult, error) {
	var result ImportResult
	var snapshots []links.Link
	names := map[string]bool{}
	decoder := json.NewDecoder(linksJSONL)
	line := 0
	for {
		line++
		var link links.Link
		err := decoder.Decode(&link)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return result, invalidSnapshot("links record %d: %v", line, err)
		}
		if err := links.ValidateShort(link.Short); err != nil {
			return result, invalidSnapshot("links record %d: %v", line, err)
		}
		if err := links.ValidateDestination(link.Long); err != nil {
			return result, invalidSnapshot("links record %d: %v", line, err)
		}
		if link.Clicks < 0 {
			return result, invalidSnapshot("links record %d: clicks cannot be negative", line)
		}
		id := links.Normalize(link.Short)
		if names[id] {
			return result, invalidSnapshot("links record %d: duplicate normalized name %q", line, link.Short)
		}
		names[id] = true
		snapshots = append(snapshots, link)
	}
	stats := map[string][]stat{}
	if statsCSV != nil {
		reader := csv.NewReader(statsCSV)
		reader.FieldsPerRecord = 3
		for row := 1; ; row++ {
			fields, err := reader.Read()
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				return result, invalidSnapshot("stats line %d: %v", row, err)
			}
			created, e1 := strconv.ParseInt(fields[1], 10, 64)
			clicks, e2 := strconv.ParseInt(fields[2], 10, 64)
			if e1 != nil || e2 != nil || clicks < 0 {
				return result, invalidSnapshot("stats line %d: invalid timestamp or click count", row)
			}
			if !names[fields[0]] {
				result.SkippedStats++
				continue
			}
			stats[fields[0]] = append(stats[fields[0]], stat{fields[0], created, clicks})
		}
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return result, err
	}
	defer tx.Rollback()
	for _, link := range snapshots {
		id := links.Normalize(link.Short)
		insert, err := tx.ExecContext(ctx, "INSERT INTO Links(ID,Short,Long,Created,LastEdit) VALUES(?,?,?,?,?) ON CONFLICT(ID) DO NOTHING", id, link.Short, link.Long, link.Created.Unix(), link.LastEdit.Unix())
		if err != nil {
			return ImportResult{}, err
		}
		count, err := insert.RowsAffected()
		if err != nil {
			return ImportResult{}, err
		}
		if count == 0 {
			result.Skipped++
			continue
		}
		result.Imported++
		values := stats[id]
		if statsCSV == nil && link.Clicks > 0 {
			values = []stat{{id, link.LastEdit.UTC().Truncate(time.Minute).Unix(), link.Clicks}}
		}
		for _, value := range values {
			if _, err := tx.ExecContext(ctx, "INSERT INTO Stats(ID,Created,Clicks) VALUES(?,?,?) ON CONFLICT(ID,Created) DO UPDATE SET Clicks=Clicks+excluded.Clicks", value.ID, value.Created, value.Clicks); err != nil {
				return ImportResult{}, err
			}
		}
	}
	if err := tx.Commit(); err != nil {
		return ImportResult{}, err
	}
	return result, nil
}
