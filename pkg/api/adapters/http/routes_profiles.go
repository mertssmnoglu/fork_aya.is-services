package http

import (
	"net/http"
	"strings"

	"github.com/eser/ajan/datafx"
	"github.com/eser/ajan/httpfx"
	"github.com/eser/ajan/logfx"
	"github.com/eser/aya.is-services/pkg/api/adapters/storage"
	"github.com/eser/aya.is-services/pkg/api/business/profiles"
	"github.com/eser/aya.is-services/pkg/api/business/stories"
	"github.com/eser/aya.is-services/pkg/lib/cursors"
)

func RegisterHttpRoutesForProfiles( //nolint:funlen,cyclop,gocognit,maintidx
	routes *httpfx.Router,
	logger *logfx.Logger,
	dataRegistry *datafx.Registry,
) {
	routes.
		Route("GET /{locale}/profiles", func(ctx *httpfx.Context) httpfx.Result {
			// get variables from path
			localeParam := ctx.Request.PathValue("locale")
			cursor := cursors.NewCursorFromRequest(ctx.Request)

			filterKind, filterKindOk := cursor.Filters["kind"]
			if !filterKindOk {
				return ctx.Results.Error(
					http.StatusBadRequest,
					[]byte("filter_kind is required"),
				)
			}

			kinds := strings.SplitSeq(filterKind, ",")
			for kind := range kinds {
				if kind != "individual" && kind != "organization" && kind != "product" {
					return ctx.Results.Error(
						http.StatusBadRequest,
						[]byte("filter_kind is invalid"),
					)
				}
			}

			repository, err := storage.NewRepositoryFromDefault(dataRegistry)
			if err != nil {
				return ctx.Results.Error(http.StatusInternalServerError, []byte(err.Error()))
			}

			service := profiles.NewService(logger, repository)

			records, err := service.List(ctx.Request.Context(), localeParam, cursor)
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

			record, err := service.GetBySlug(ctx.Request.Context(), localeParam, slugParam)
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
		HasSummary("List profile pages by profile slug").
		HasDescription("List profile pages by profile slug.").
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
		HasSummary("List profile page by profile slug and page slug").
		HasDescription("List profile page by profile slug and page slug.").
		HasResponse(http.StatusOK)

	routes.
		Route("GET /{locale}/profiles/{slug}/links", func(ctx *httpfx.Context) httpfx.Result {
			// get variables from path
			localeParam := ctx.Request.PathValue("locale")
			slugParam := ctx.Request.PathValue("slug")

			repository, err := storage.NewRepositoryFromDefault(dataRegistry)
			if err != nil {
				return ctx.Results.Error(http.StatusInternalServerError, []byte(err.Error()))
			}

			service := profiles.NewService(logger, repository)

			record, err := service.ListLinksBySlug(
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
		HasSummary("List profile links by profile slug").
		HasDescription("List profile links by profile slug.").
		HasResponse(http.StatusOK)

	routes.
		Route("GET /{locale}/profiles/{slug}/stories", func(ctx *httpfx.Context) httpfx.Result {
			// get variables from path
			localeParam := ctx.Request.PathValue("locale")
			slugParam := ctx.Request.PathValue("slug")
			cursor := cursors.NewCursorFromRequest(ctx.Request)

			repository, err := storage.NewRepositoryFromDefault(dataRegistry)
			if err != nil {
				return ctx.Results.Error(http.StatusInternalServerError, []byte(err.Error()))
			}

			service := stories.NewService(logger, repository)

			records, err := service.ListByAuthorProfileSlug(
				ctx.Request.Context(),
				localeParam,
				slugParam,
				cursor,
			)
			if err != nil {
				return ctx.Results.Error(http.StatusInternalServerError, []byte(err.Error()))
			}

			return ctx.Results.Json(records)
		}).
		HasSummary("List stories authored by profile slug").
		HasDescription("List stories authored by profile slug.").
		HasResponse(http.StatusOK)

	routes.
		Route("GET /{locale}/profiles/{slug}/memberships", func(ctx *httpfx.Context) httpfx.Result {
			// get variables from path
			localeParam := ctx.Request.PathValue("locale")
			slugParam := ctx.Request.PathValue("slug")
			cursor := cursors.NewCursorFromRequest(ctx.Request)

			filterProfileKind, filterProfileKindOk := cursor.Filters["profile_kind"]
			if !filterProfileKindOk {
				return ctx.Results.Error(
					http.StatusBadRequest,
					[]byte("filter_profile_kind is required"),
				)
			}

			profileKinds := strings.Split(filterProfileKind, ",")
			for _, profileKind := range profileKinds {
				if profileKind != "individual" && profileKind != "organization" &&
					profileKind != "product" {
					return ctx.Results.Error(
						http.StatusBadRequest,
						[]byte("filter_profile_kind is invalid"),
					)
				}
			}

			repository, err := storage.NewRepositoryFromDefault(dataRegistry)
			if err != nil {
				return ctx.Results.Error(http.StatusInternalServerError, []byte(err.Error()))
			}

			service := profiles.NewService(logger, repository)

			records, err := service.ListProfileMembershipsBySlugAndKinds(
				ctx.Request.Context(),
				localeParam,
				slugParam,
				profileKinds,
				cursor,
			)
			if err != nil {
				return ctx.Results.Error(http.StatusInternalServerError, []byte(err.Error()))
			}

			return ctx.Results.Json(records)
		}).
		HasSummary("List profile memberships by profile slug and kind").
		HasDescription("List profile memberships by profile slug and kind.").
		HasResponse(http.StatusOK)
}
