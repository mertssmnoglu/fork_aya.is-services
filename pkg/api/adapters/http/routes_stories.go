//nolint:dupl
package http

import (
	"net/http"

	"github.com/eser/ajan/datafx"
	"github.com/eser/ajan/httpfx"
	"github.com/eser/ajan/logfx"
	"github.com/eser/aya.is-services/pkg/api/adapters/storage"
	"github.com/eser/aya.is-services/pkg/api/business/stories"
)

func RegisterHttpRoutesForStories(routes *httpfx.Router, logger *logfx.Logger, dataRegistry *datafx.Registry) {
	routes.
		Route("GET /{locale}/stories", func(ctx *httpfx.Context) httpfx.Result {
			// get variables from path
			// localeParam := ctx.Request.PathValue("locale")

			store, err := storage.NewFromDefault(dataRegistry)
			if err != nil {
				return ctx.Results.Error(http.StatusInternalServerError, []byte(err.Error()))
			}

			service := stories.NewService(logger, store)

			records, err := service.List(ctx.Request.Context())
			if err != nil {
				return ctx.Results.Error(http.StatusInternalServerError, []byte(err.Error()))
			}

			return ctx.Results.Json(records)
		}).
		HasSummary("List stories").
		HasDescription("List stories.").
		HasResponse(http.StatusOK)

	routes.
		Route("GET /{locale}/stories/{slug}", func(ctx *httpfx.Context) httpfx.Result {
			// get variables from path
			// localeParam := ctx.Request.PathValue("locale")
			slugParam := ctx.Request.PathValue("slug")

			store, err := storage.NewFromDefault(dataRegistry)
			if err != nil {
				return ctx.Results.Error(http.StatusInternalServerError, []byte(err.Error()))
			}

			service := stories.NewService(logger, store)

			record, err := service.GetBySlug(ctx.Request.Context(), slugParam)
			if err != nil {
				return ctx.Results.Error(http.StatusInternalServerError, []byte(err.Error()))
			}

			return ctx.Results.Json(record)
		}).
		HasSummary("Get story by slug").
		HasDescription("Get story by slug.").
		HasResponse(http.StatusOK)
}
