package httpapi

import (
	"context"
	"errors"
	"log"
	"net/http"
	"net/url"

	"github.com/danielfrg/golinks/internal/links"
	"github.com/danielfrg/golinks/internal/linkservice"
	"github.com/danielfrg/golinks/internal/store"
	"github.com/danielgtaylor/huma/v2"
)

type linkInput struct {
	Short string `path:"short"`
}
type linkOutput struct{ Body links.Link }
type createInput struct {
	Body struct {
		Short string `json:"short"`
		Long  string `json:"long"`
	}
}
type updateInput struct {
	Short string `path:"short"`
	Body  struct {
		Long string `json:"long"`
	}
}
type listOutput struct {
	Body struct {
		Links []links.Link `json:"links"`
	}
}
type settingsOutput struct {
	Body struct {
		PublicURL string `json:"publicUrl"`
		ReadOnly  bool   `json:"readOnly"`
	}
}
type previewInput struct {
	Body struct {
		Long  string `json:"long"`
		Path  string `json:"path"`
		Query string `json:"query"`
	}
}
type previewOutput struct {
	Body struct {
		URL string `json:"url"`
	}
}

func apiError(err error) error {
	var validation *linkservice.ValidationError
	if errors.As(err, &validation) {
		return huma.Error422UnprocessableEntity(validation.Error())
	}
	if errors.Is(err, linkservice.ErrReadOnly) {
		return huma.Error405MethodNotAllowed(err.Error())
	}
	if errors.Is(err, store.ErrNotFound) {
		return huma.Error404NotFound(err.Error())
	}
	if errors.Is(err, store.ErrConflict) {
		return huma.Error409Conflict(err.Error())
	}
	log.Printf("database request: %v", err)
	return huma.Error500InternalServerError("Could not access links. Try again.")
}
func registerLinks(api huma.API, s *store.Store, options Options) {
	service := linkservice.New(s, linkservice.Options{PublicURL: options.PublicURL, ReadOnly: options.ReadOnly, Now: options.Now})
	huma.Register(api, huma.Operation{OperationID: "list-links", Method: "GET", Path: "/_/api/links", Summary: "List links"}, func(ctx context.Context, _ *struct{}) (*listOutput, error) {
		list, err := s.List(ctx)
		if err != nil {
			return nil, apiError(err)
		}
		out := &listOutput{}
		out.Body.Links = list
		return out, nil
	})
	huma.Register(api, huma.Operation{OperationID: "get-link", Method: "GET", Path: "/_/api/links/{short}", Summary: "Get a link"}, func(ctx context.Context, in *linkInput) (*linkOutput, error) {
		link, err := s.Get(ctx, in.Short)
		if err != nil {
			return nil, apiError(err)
		}
		return &linkOutput{Body: link}, nil
	})
	huma.Register(api, huma.Operation{OperationID: "create-link", Method: "POST", Path: "/_/api/links", Summary: "Create a link", DefaultStatus: 201}, func(ctx context.Context, in *createInput) (*linkOutput, error) {
		link, err := service.Create(ctx, in.Body.Short, in.Body.Long)
		if err != nil {
			return nil, apiError(err)
		}
		return &linkOutput{Body: link}, nil
	})
	huma.Register(api, huma.Operation{OperationID: "update-link", Method: "PUT", Path: "/_/api/links/{short}", Summary: "Update a destination"}, func(ctx context.Context, in *updateInput) (*linkOutput, error) {
		link, err := service.Update(ctx, in.Short, in.Body.Long)
		if err != nil {
			return nil, apiError(err)
		}
		return &linkOutput{Body: link}, nil
	})
	huma.Register(api, huma.Operation{OperationID: "delete-link", Method: "DELETE", Path: "/_/api/links/{short}", Summary: "Delete a link", DefaultStatus: 204}, func(ctx context.Context, in *linkInput) (*struct{}, error) {
		if err := s.Delete(ctx, in.Short); err != nil {
			return nil, apiError(err)
		}
		return nil, nil
	})
	huma.Register(api, huma.Operation{OperationID: "get-settings", Method: "GET", Path: "/_/api/settings", Summary: "Get UI settings"}, func(_ context.Context, _ *struct{}) (*settingsOutput, error) {
		out := &settingsOutput{}
		out.Body.PublicURL = options.PublicURL
		out.Body.ReadOnly = options.ReadOnly
		return out, nil
	})
	huma.Register(api, huma.Operation{OperationID: "preview-destination", Method: http.MethodPost, Path: "/_/api/resolve", Summary: "Preview a destination without visiting it"}, func(_ context.Context, in *previewInput) (*previewOutput, error) {
		query, err := url.ParseQuery(in.Body.Query)
		if err != nil {
			return nil, huma.Error422UnprocessableEntity("Query string is invalid")
		}
		u, err := service.Expand(in.Body.Long, in.Body.Path, query)
		if err != nil {
			return nil, huma.Error422UnprocessableEntity(err.Error())
		}
		out := &previewOutput{}
		out.Body.URL = u
		return out, nil
	})
}
