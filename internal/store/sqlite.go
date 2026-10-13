package store

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"golinks/internal/links"
	_ "modernc.org/sqlite"
)

var (
	ErrNotFound = errors.New("Link not found")
	ErrConflict = errors.New("A link with that name already exists (case and hyphens are ignored)")
)

//go:embed schema.sql
var schema string

type Store struct{ db *sql.DB }

func Open(filename string) (*Store, error) {
	if filename == "" {
		return nil, errors.New("Database path cannot be empty")
	}
	if filename != ":memory:" {
		if err := os.MkdirAll(filepath.Dir(filename), 0700); err != nil {
			return nil, fmt.Errorf("create database directory: %w", err)
		}
	}
	db, err := sql.Open("sqlite", filename)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	fail := func(err error) (*Store, error) { db.Close(); return nil, fmt.Errorf("open database: %w", err) }
	if err := db.Ping(); err != nil {
		return fail(err)
	}
	for _, statement := range []string{"PRAGMA busy_timeout=5000", "PRAGMA journal_mode=DELETE", "PRAGMA synchronous=FULL", "PRAGMA foreign_keys=ON", schema} {
		if _, err := db.Exec(statement); err != nil {
			return fail(err)
		}
	}
	return &Store{db: db}, nil
}

func (s *Store) Close() error                   { return s.db.Close() }
func (s *Store) Ping(ctx context.Context) error { return s.db.PingContext(ctx) }

const selectLinks = `SELECT l.Short,l.Long,l.Created,l.LastEdit,
 COALESCE((SELECT SUM(Clicks) FROM Stats WHERE ID=l.ID),0) FROM Links l`

type scanner interface{ Scan(...any) error }

func scanLink(row scanner) (links.Link, error) {
	var link links.Link
	var created, edited int64
	err := row.Scan(&link.Short, &link.Long, &created, &edited, &link.Clicks)
	if errors.Is(err, sql.ErrNoRows) {
		return link, ErrNotFound
	}
	link.Created = time.Unix(created, 0).UTC()
	link.LastEdit = time.Unix(edited, 0).UTC()
	return link, err
}

func (s *Store) List(ctx context.Context) ([]links.Link, error) {
	rows, err := s.db.QueryContext(ctx, selectLinks+" ORDER BY l.Short COLLATE NOCASE,l.Short")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]links.Link, 0)
	for rows.Next() {
		link, err := scanLink(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, link)
	}
	return out, rows.Err()
}
func (s *Store) Get(ctx context.Context, short string) (links.Link, error) {
	return scanLink(s.db.QueryRowContext(ctx, selectLinks+" WHERE l.ID=?", links.Normalize(short)))
}

func (s *Store) Create(ctx context.Context, short, destination string) (links.Link, error) {
	if err := links.ValidateShort(short); err != nil {
		return links.Link{}, err
	}
	if err := links.ValidateDestination(destination); err != nil {
		return links.Link{}, err
	}
	now := time.Now().UTC().Unix()
	result, err := s.db.ExecContext(ctx, "INSERT INTO Links(ID,Short,Long,Created,LastEdit) VALUES(?,?,?,?,?) ON CONFLICT(ID) DO NOTHING", links.Normalize(short), short, destination, now, now)
	if err != nil {
		return links.Link{}, err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return links.Link{}, err
	}
	if n == 0 {
		return links.Link{}, ErrConflict
	}
	return s.Get(ctx, short)
}
func (s *Store) Update(ctx context.Context, short, destination string) (links.Link, error) {
	if err := links.ValidateDestination(destination); err != nil {
		return links.Link{}, err
	}
	result, err := s.db.ExecContext(ctx, "UPDATE Links SET Long=?,LastEdit=? WHERE ID=?", destination, time.Now().UTC().Unix(), links.Normalize(short))
	if err != nil {
		return links.Link{}, err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return links.Link{}, err
	}
	if n == 0 {
		return links.Link{}, ErrNotFound
	}
	return s.Get(ctx, short)
}
func (s *Store) Delete(ctx context.Context, short string) error {
	result, err := s.db.ExecContext(ctx, "DELETE FROM Links WHERE ID=?", links.Normalize(short))
	if err != nil {
		return err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}
func (s *Store) RecordClick(ctx context.Context, short string, at time.Time) error {
	_, err := s.db.ExecContext(ctx, "INSERT INTO Stats(ID,Created,Clicks) VALUES(?,?,1) ON CONFLICT(ID,Created) DO UPDATE SET Clicks=Clicks+1", links.Normalize(short), at.UTC().Truncate(time.Minute).Unix())
	return err
}
