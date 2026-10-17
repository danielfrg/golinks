package mcp

import (
	"context"
	"errors"
	"github.com/danielfrg/golinks/internal/links"
	"github.com/danielfrg/golinks/internal/linkservice"
	"github.com/danielfrg/golinks/internal/store"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

type Link struct {
	links.Link
	PublicURL string `json:"public_url"`
}
type SearchInput struct {
	Paging
	Query string `json:"query,omitempty"`
	Sort  string `json:"sort,omitempty" jsonschema:"name (default) or clicks"`
}
type SearchOutput struct {
	Links      []Link `json:"links"`
	Total      int    `json:"total"`
	NextOffset *int   `json:"next_offset"`
}
type GetInput struct {
	Short string `json:"short"`
}
type WriteInput struct {
	Short       string `json:"short"`
	Destination string `json:"destination"`
}
type ResolveInput struct {
	Short string `json:"short"`
	Path  string `json:"path,omitempty"`
	Query string `json:"query,omitempty"`
}
type ResolveOutput struct {
	URL string `json:"url"`
}

func annotation(read bool) *sdk.ToolAnnotations {
	no := false
	return &sdk.ToolAnnotations{ReadOnlyHint: read, IdempotentHint: read, DestructiveHint: &no, OpenWorldHint: &no}
}
func safeError(err error) error {
	if err == nil {
		return nil
	}
	var validation *linkservice.ValidationError
	switch {
	case errors.Is(err, store.ErrNotFound):
		return errors.New("Link not found")
	case errors.Is(err, store.ErrConflict):
		return errors.New("Short link already exists")
	case errors.Is(err, linkservice.ErrReadOnly):
		return linkservice.ErrReadOnly
	case errors.As(err, &validation):
		return errors.New("Invalid short link or destination")
	case errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded):
		return errors.New("Link request canceled")
	default:
		return errors.New("Links unavailable; try again")
	}
}
func NewHandler(service *linkservice.Service, hosts []string) http.Handler {
	server := sdk.NewServer(&sdk.Implementation{Name: "golinks", Version: "0.1.0"}, nil)
	present := func(v links.Link) Link { return Link{Link: v, PublicURL: service.PublicURL(v.Short)} }
	sdk.AddTool(server, &sdk.Tool{Name: "search_links", Description: "Search short links and destinations without visiting them. Sort by name or descending clicks.", Annotations: annotation(true)}, func(ctx context.Context, _ *sdk.CallToolRequest, in SearchInput) (*sdk.CallToolResult, SearchOutput, error) {
		out := SearchOutput{Links: []Link{}}
		if _, _, _, err := in.Paging.bounds(0); err != nil {
			return nil, out, err
		}
		if in.Sort != "" && in.Sort != "name" && in.Sort != "clicks" {
			return nil, out, errors.New("sort must be name or clicks")
		}
		all, err := service.List(ctx)
		if err != nil {
			return nil, out, safeError(err)
		}
		filtered := []links.Link{}
		q := strings.ToLower(in.Query)
		for _, v := range all {
			if strings.Contains(strings.ToLower(v.Short+" "+v.Long), q) {
				filtered = append(filtered, v)
			}
		}
		sort.Slice(filtered, func(i, j int) bool {
			a, b := filtered[i], filtered[j]
			if in.Sort == "clicks" && a.Clicks != b.Clicks {
				return a.Clicks > b.Clicks
			}
			return a.Short < b.Short
		})
		start, end, next, _ := in.Paging.bounds(len(filtered))
		out.Total = len(filtered)
		out.NextOffset = next
		for _, v := range filtered[start:end] {
			out.Links = append(out.Links, present(v))
		}
		return nil, out, nil
	})
	sdk.AddTool(server, &sdk.Tool{Name: "get_link", Description: "Get a short link and its stored destination without visiting it.", Annotations: annotation(true)}, func(ctx context.Context, _ *sdk.CallToolRequest, in GetInput) (*sdk.CallToolResult, Link, error) {
		v, err := service.Get(ctx, in.Short)
		return nil, present(v), safeError(err)
	})
	sdk.AddTool(server, &sdk.Tool{Name: "resolve_link", Description: "Preview the immediate redirect, including relative aliases, without following it or recording a visit. query is a URL query string.", Annotations: annotation(true)}, func(ctx context.Context, _ *sdk.CallToolRequest, in ResolveInput) (*sdk.CallToolResult, ResolveOutput, error) {
		q, err := url.ParseQuery(in.Query)
		if err != nil {
			return nil, ResolveOutput{}, errors.New("Invalid query string")
		}
		value, err := service.Resolve(ctx, in.Short, in.Path, q)
		return nil, ResolveOutput{value}, safeError(err)
	})
	if !service.Options.ReadOnly {
		sdk.AddTool(server, &sdk.Tool{Name: "create_link", Description: "Create a short link. Fails if its normalized name already exists.", Annotations: annotation(false)}, func(ctx context.Context, _ *sdk.CallToolRequest, in WriteInput) (*sdk.CallToolResult, Link, error) {
			v, err := service.Create(ctx, in.Short, in.Destination)
			return nil, present(v), safeError(err)
		})
		sdk.AddTool(server, &sdk.Tool{Name: "update_link", Description: "Replace the destination of an existing short link.", Annotations: annotation(false)}, func(ctx context.Context, _ *sdk.CallToolRequest, in WriteInput) (*sdk.CallToolResult, Link, error) {
			v, err := service.Update(ctx, in.Short, in.Destination)
			return nil, present(v), safeError(err)
		})
	}
	return transport(server, hosts)
}
