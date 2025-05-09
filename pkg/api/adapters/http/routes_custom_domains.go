package http

import (
	"net/http"

	"github.com/eser/ajan/datafx"
	"github.com/eser/ajan/httpfx"
	"github.com/eser/ajan/logfx"
	"github.com/eser/aya.is-services/pkg/api/adapters/storage"
	"github.com/eser/aya.is-services/pkg/api/business/custom_domains"
)

func RegisterHttpRoutesForCustomDomains(routes *httpfx.Router, logger *logfx.Logger, dataRegistry *datafx.Registry) {
	routes.
		Route("GET /{locale}/custom-domains/{domain}", func(ctx *httpfx.Context) httpfx.Result {
			// get variables from path
			// localeParam := ctx.Request.PathValue("locale")
			domainParam := ctx.Request.PathValue("domain")

			store, err := storage.NewFromDefault(dataRegistry)
			if err != nil {
				return ctx.Results.Error(http.StatusInternalServerError, []byte(err.Error()))
			}

			service := custom_domains.NewService(logger, store)

			records, err := service.GetByDomain(ctx.Request.Context(), domainParam)
			if err != nil {
				return ctx.Results.Error(http.StatusInternalServerError, []byte(err.Error()))
			}

			return ctx.Results.Json(records)
		}).
		HasSummary("Get custom domain by domain").
		HasDescription("Get custom domain by domain.").
		HasResponse(http.StatusOK)
}
