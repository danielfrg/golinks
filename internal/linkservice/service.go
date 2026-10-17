// Package linkservice shares validated link operations across HTTP and MCP.
package linkservice

import (
	"context"
	"errors"
	"net/url"
	"strings"
	"time"

	"github.com/danielfrg/golinks/internal/links"
	"github.com/danielfrg/golinks/internal/store"
)

var ErrReadOnly = errors.New("Golinks is read-only")

type ValidationError struct{ Message string }

func (e *ValidationError) Error() string { return e.Message }

type Options struct {
	PublicURL string
	ReadOnly  bool
	Now       func() time.Time
}
type Service struct {
	store   *store.Store
	Options Options
}

func New(s *store.Store, options Options) *Service {
	if options.Now == nil {
		options.Now = time.Now
	}
	return &Service{store: s, Options: options}
}
func (s *Service) List(ctx context.Context) ([]links.Link, error) { return s.store.List(ctx) }
func (s *Service) Get(ctx context.Context, short string) (links.Link, error) {
	return s.store.Get(ctx, short)
}
func (s *Service) Create(ctx context.Context, short, destination string) (links.Link, error) {
	if s.Options.ReadOnly {
		return links.Link{}, ErrReadOnly
	}
	if err := links.ValidateShort(short); err != nil {
		return links.Link{}, &ValidationError{err.Error()}
	}
	if err := links.ValidateDestination(destination); err != nil {
		return links.Link{}, &ValidationError{err.Error()}
	}
	return s.store.Create(ctx, short, destination)
}
func (s *Service) Update(ctx context.Context, short, destination string) (links.Link, error) {
	if s.Options.ReadOnly {
		return links.Link{}, ErrReadOnly
	}
	if err := links.ValidateDestination(destination); err != nil {
		return links.Link{}, &ValidationError{err.Error()}
	}
	return s.store.Update(ctx, short, destination)
}
func (s *Service) Expand(destination, path string, query url.Values) (string, error) {
	value, err := links.Expand(destination, links.Environment{Path: path, Query: query, Now: s.Options.Now()})
	if err != nil {
		return "", &ValidationError{err.Error()}
	}
	return value.String(), nil
}
func (s *Service) Resolve(ctx context.Context, short, path string, query url.Values) (string, error) {
	link, err := s.Get(ctx, short)
	if errors.Is(err, store.ErrNotFound) {
		trimmed := strings.TrimRight(short, ".,()[]{}")
		if trimmed != short {
			link, err = s.Get(ctx, trimmed)
		}
	}
	if err != nil {
		return "", err
	}
	return s.Expand(link.Long, path, query)
}
func (s *Service) PublicURL(short string) string {
	return strings.TrimRight(s.Options.PublicURL, "/") + "/" + url.PathEscape(short)
}
