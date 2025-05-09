package http

import (
	"net/http"

	"github.com/eser/ajan/datafx"
	"github.com/eser/ajan/httpfx"
	"github.com/eser/ajan/logfx"
	"github.com/eser/aya.is-services/pkg/api/adapters/storage"
	"github.com/eser/aya.is-services/pkg/api/business/profiles"
)

func RegisterHttpRoutesForProfiles(routes *httpfx.Router, logger *logfx.Logger, dataRegistry *datafx.Registry) {
	routes.
		Route("GET /{locale}/profiles", func(ctx *httpfx.Context) httpfx.Result {
			// get variables from path
			// localeParam := ctx.Request.PathValue("locale")

			store, err := storage.NewFromDefault(dataRegistry)
			if err != nil {
				return ctx.Results.Error(http.StatusInternalServerError, []byte(err.Error()))
			}

			service := profiles.NewService(logger, store)

			records, err := service.List(ctx.Request.Context())
			if err != nil {
				return ctx.Results.Error(http.StatusInternalServerError, []byte(err.Error()))
			}

			return ctx.Results.Json(records)
		}).
		HasSummary("List profiles").
		HasDescription("List profiles.").
		HasResponse(http.StatusOK)
}
