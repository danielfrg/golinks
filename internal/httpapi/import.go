package httpapi

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/danielfrg/golinks/internal/store"
	"github.com/danielgtaylor/huma/v2"
)

const maxImportBody int64 = 16 << 20

type importInput struct {
	Body struct {
		Links string  `json:"links" doc:"JSONL file downloaded from Export links"`
		Stats *string `json:"stats,omitempty" doc:"Optional CSV file downloaded from /_/export-stats"`
	}
}

type importOutput struct{ Body store.ImportResult }

func registerImport(api huma.API, s *store.Store) {
	huma.Register(api, huma.Operation{
		OperationID: "import-links", Method: http.MethodPost, Path: "/_/api/import",
		Summary: "Import a link backup", MaxBodyBytes: maxImportBody,
	}, func(ctx context.Context, in *importInput) (*importOutput, error) {
		var stats io.Reader
		if in.Body.Stats != nil {
			stats = strings.NewReader(*in.Body.Stats)
		}
		result, err := s.Import(ctx, strings.NewReader(in.Body.Links), stats)
		if errors.Is(err, store.ErrInvalidSnapshot) {
			return nil, huma.Error422UnprocessableEntity(err.Error())
		}
		if err != nil {
			return nil, apiError(err)
		}
		return &importOutput{Body: result}, nil
	})
}
