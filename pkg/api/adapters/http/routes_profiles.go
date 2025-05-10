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

func RegisterHttpRoutesForProfiles( //nolint:funlen,cyclop
	routes *httpfx.Router,
	logger *logfx.Logger,
	dataRegistry *datafx.Registry,
) {
	routes.
		Route("GET /{locale}/profiles", func(ctx *httpfx.Context) httpfx.Result {
			// get variables from path
			localeParam := ctx.Request.PathValue("locale")
			cursor := cursors.NewCursorFromRequest(ctx.Request)

			repository, err := storage.NewRepositoryFromDefault(dataRegistry)
			if err != nil {
				return ctx.Results.Error(http.StatusInternalServerError, []byte(err.Error()))
			}

			service := profiles.NewService(logger, repository)

			records, err := service.ListWithCursor(ctx.Request.Context(), localeParam, cursor)
			if err != nil {
				return ctx.Results.Error(http.StatusInternalServerError, []byte(err.Error()))
			}

			return ctx.Results.Json(records)
		}).
		HasSummary("List profiles").
		HasDescription("List profiles.").
		HasResponse(http.StatusOK)

	routes.
		Route("GET /{locale}/profiles/{slug}", func(ctx *httpfx.Context) httpfx.Result {
			// get variables from path
			localeParam := ctx.Request.PathValue("locale")
			slugParam := ctx.Request.PathValue("slug")

			repository, err := storage.NewRepositoryFromDefault(dataRegistry)
			if err != nil {
				return ctx.Results.Error(http.StatusInternalServerError, []byte(err.Error()))
			}

			service := profiles.NewService(logger, repository)

			record, err := service.GetBySlugEx(ctx.Request.Context(), localeParam, slugParam)
			if err != nil {
				return ctx.Results.Error(http.StatusInternalServerError, []byte(err.Error()))
			}

			wrappedResponse := cursors.WrapResponseWithCursor(record, nil)

			return ctx.Results.Json(wrappedResponse)
		}).
		HasSummary("Get profile by slug").
		HasDescription("Get profile by slug.").
		HasResponse(http.StatusOK)

	routes.
		Route("GET /{locale}/profiles/{slug}/pages", func(ctx *httpfx.Context) httpfx.Result {
			// get variables from path
			localeParam := ctx.Request.PathValue("locale")
			slugParam := ctx.Request.PathValue("slug")

			repository, err := storage.NewRepositoryFromDefault(dataRegistry)
			if err != nil {
				return ctx.Results.Error(http.StatusInternalServerError, []byte(err.Error()))
			}

			service := profiles.NewService(logger, repository)

			record, err := service.ListPagesBySlug(
				ctx.Request.Context(),
				localeParam,
				slugParam,
			)
			if err != nil {
				return ctx.Results.Error(http.StatusInternalServerError, []byte(err.Error()))
			}

			wrappedResponse := cursors.WrapResponseWithCursor(record, nil)

			return ctx.Results.Json(wrappedResponse)
		}).
		HasSummary("Get profile pages by profile slug").
		HasDescription("Get profile pages by profile slug.").
		HasResponse(http.StatusOK)

	routes.
		Route(
			"GET /{locale}/profiles/{slug}/pages/{pageSlug}",
			func(ctx *httpfx.Context) httpfx.Result {
				// get variables from path
				localeParam := ctx.Request.PathValue("locale")
				slugParam := ctx.Request.PathValue("slug")
				pageSlugParam := ctx.Request.PathValue("pageSlug")

				repository, err := storage.NewRepositoryFromDefault(dataRegistry)
				if err != nil {
					return ctx.Results.Error(http.StatusInternalServerError, []byte(err.Error()))
				}

				service := profiles.NewService(logger, repository)

				record, err := service.GetPageBySlug(
					ctx.Request.Context(),
					localeParam,
					slugParam,
					pageSlugParam,
				)
				if err != nil {
					return ctx.Results.Error(http.StatusInternalServerError, []byte(err.Error()))
				}

				wrappedResponse := cursors.WrapResponseWithCursor(record, nil)

				return ctx.Results.Json(wrappedResponse)
			},
		).
		HasSummary("Get profile page by profile slug and page slug").
		HasDescription("Get profile page by profile slug and page slug.").
		HasResponse(http.StatusOK)

	routes.
		Route("GET /{locale}/custom-domains/{domain}", func(ctx *httpfx.Context) httpfx.Result {
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
		}).
		HasSummary("Get profile by custom domain").
		HasDescription("Get profile by custom domain.").
		HasResponse(http.StatusOK)
}
