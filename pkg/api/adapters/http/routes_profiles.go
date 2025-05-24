package http

import (
	"net/http"
	"strings"

	"github.com/eser/ajan/httpfx"
	"github.com/eser/ajan/logfx"
	"github.com/eser/aya.is-services/pkg/api/business/profiles"
	"github.com/eser/aya.is-services/pkg/api/business/stories"
	"github.com/eser/aya.is-services/pkg/lib/cursors"
)

func RegisterHttpRoutesForProfiles( //nolint:funlen,cyclop,gocognit,maintidx
	routes *httpfx.Router,
	logger *logfx.Logger,
	profilesService *profiles.Service,
	storiesService *stories.Service,
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

			records, err := profilesService.List(ctx.Request.Context(), localeParam, cursor)
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

			record, err := profilesService.GetBySlugEx(
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
		HasSummary("Get profile by slug").
		HasDescription("Get profile by slug.").
		HasResponse(http.StatusOK)

	routes.
		Route("GET /{locale}/profiles/{slug}/pages", func(ctx *httpfx.Context) httpfx.Result {
			// get variables from path
			localeParam := ctx.Request.PathValue("locale")
			slugParam := ctx.Request.PathValue("slug")

			record, err := profilesService.ListPagesBySlug(
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

				record, err := profilesService.GetPageBySlug(
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

			record, err := profilesService.ListLinksBySlug(
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

			records, err := storiesService.ListByPublicationProfileSlug(
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
		HasSummary("List stories published to profile slug").
		HasDescription("List stories published to profile slug.").
		HasResponse(http.StatusOK)

	routes.
		Route(
			"GET /{locale}/profiles/{slug}/contributions",
			func(ctx *httpfx.Context) httpfx.Result {
				// get variables from path
				localeParam := ctx.Request.PathValue("locale")
				slugParam := ctx.Request.PathValue("slug")
				cursor := cursors.NewCursorFromRequest(ctx.Request)

				records, err := profilesService.ListProfileContributionsBySlug(
					ctx.Request.Context(),
					localeParam,
					slugParam,
					cursor,
				)
				if err != nil {
					return ctx.Results.Error(http.StatusInternalServerError, []byte(err.Error()))
				}

				return ctx.Results.Json(records)
			},
		).
		HasSummary("List profile contributions by profile slug").
		HasDescription("List profile contributions by profile slug.").
		HasResponse(http.StatusOK)

	routes.
		Route(
			"GET /{locale}/profiles/{slug}/members",
			func(ctx *httpfx.Context) httpfx.Result {
				// get variables from path
				localeParam := ctx.Request.PathValue("locale")
				slugParam := ctx.Request.PathValue("slug")
				cursor := cursors.NewCursorFromRequest(ctx.Request)

				records, err := profilesService.ListProfileMembersBySlug(
					ctx.Request.Context(),
					localeParam,
					slugParam,
					cursor,
				)
				if err != nil {
					return ctx.Results.Error(http.StatusInternalServerError, []byte(err.Error()))
				}

				return ctx.Results.Json(records)
			},
		).
		HasSummary("List profile members by profile slug").
		HasDescription("List profile members by profile slug.").
		HasResponse(http.StatusOK)
}
