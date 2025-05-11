//nolint:dupl
package http

import (
	"net/http"

	"github.com/eser/ajan/datafx"
	"github.com/eser/ajan/httpfx"
	"github.com/eser/ajan/logfx"
	"github.com/eser/aya.is-services/pkg/api/adapters/storage"
	"github.com/eser/aya.is-services/pkg/api/business/users"
	"github.com/eser/aya.is-services/pkg/lib/cursors"
)

func RegisterHttpRoutesForUsers(
	routes *httpfx.Router,
	logger *logfx.Logger,
	dataRegistry *datafx.Registry,
) {
	routes.
		Route("GET /{locale}/users", func(ctx *httpfx.Context) httpfx.Result {
			// get variables from path
			localeParam := ctx.Request.PathValue("locale")
			cursor := cursors.NewCursorFromRequest(ctx.Request)

			repository, err := storage.NewRepositoryFromDefault(dataRegistry)
			if err != nil {
				return ctx.Results.Error(http.StatusInternalServerError, []byte(err.Error()))
			}

			service := users.NewService(logger, repository)

			records, err := service.List(ctx.Request.Context(), localeParam, cursor)
			if err != nil {
				return ctx.Results.Error(http.StatusInternalServerError, []byte(err.Error()))
			}

			return ctx.Results.Json(records)
		}).
		HasSummary("List users").
		HasDescription("List users.").
		HasResponse(http.StatusOK)

	routes.
		Route("GET /{locale}/users/{id}", func(ctx *httpfx.Context) httpfx.Result {
			// get variables from path
			localeParam := ctx.Request.PathValue("locale")
			idParam := ctx.Request.PathValue("id")

			repository, err := storage.NewRepositoryFromDefault(dataRegistry)
			if err != nil {
				return ctx.Results.Error(http.StatusInternalServerError, []byte(err.Error()))
			}

			service := users.NewService(logger, repository)

			record, err := service.GetById(ctx.Request.Context(), localeParam, idParam)
			if err != nil {
				return ctx.Results.Error(http.StatusInternalServerError, []byte(err.Error()))
			}

			wrappedResponse := cursors.WrapResponseWithCursor(record, nil)

			return ctx.Results.Json(wrappedResponse)
		}).
		HasSummary("Get user by ID").
		HasDescription("Get user by ID.").
		HasResponse(http.StatusOK)
}
