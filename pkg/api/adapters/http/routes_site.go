package http

import (
	"net/http"

	"github.com/eser/ajan/datafx"
	"github.com/eser/ajan/httpfx"
	"github.com/eser/ajan/logfx"
	"github.com/eser/aya.is-services/pkg/api/adapters/storage"
	"github.com/eser/aya.is-services/pkg/api/business/profiles"
	"github.com/eser/aya.is-services/pkg/lib/cursors"
)

func RegisterHttpRoutesForSite( //nolint:funlen
	routes *httpfx.Router,
	logger *logfx.Logger,
	dataRegistry *datafx.Registry,
) {
	routes.
		Route(
			"GET /{locale}/site/custom-domains/{domain}",
			func(ctx *httpfx.Context) httpfx.Result {
				// get variables from path
				localeParam := ctx.Request.PathValue("locale")
				domainParam := ctx.Request.PathValue("domain")

				repository, err := storage.NewRepositoryFromDefault(dataRegistry)
				if err != nil {
					return ctx.Results.Error(http.StatusInternalServerError, []byte(err.Error()))
				}

				service := profiles.NewService(logger, repository)

				records, err := service.GetByCustomDomain(
					ctx.Request.Context(),
					localeParam,
					domainParam,
				)
				if err != nil {
					return ctx.Results.Error(http.StatusInternalServerError, []byte(err.Error()))
				}

				wrappedResponse := cursors.WrapResponseWithCursor(records, nil)

				return ctx.Results.Json(wrappedResponse)
			},
		).
		HasSummary("Get profile by a custom domain").
		HasDescription("Get profile by a custom domain.").
		HasResponse(http.StatusOK)

	routes.
		Route("GET /{locale}/site/spotlight", func(ctx *httpfx.Context) httpfx.Result {
			// get variables from path
			localeParam := ctx.Request.PathValue("locale")

			repository, err := storage.NewRepositoryFromDefault(dataRegistry)
			if err != nil {
				return ctx.Results.Error(http.StatusInternalServerError, []byte(err.Error()))
			}

			service := profiles.NewService(logger, repository)

			records, err := service.List(
				ctx.Request.Context(),
				localeParam,
				cursors.NewCursor(0, nil),
			)
			if err != nil {
				return ctx.Results.Error(http.StatusInternalServerError, []byte(err.Error()))
			}

			return ctx.Results.Json(records)
		}).
		HasSummary("Gets spotlight metadata").
		HasDescription("Gets spotlight metadata.").
		HasResponse(http.StatusOK)
}
