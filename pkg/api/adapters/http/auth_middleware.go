package http

import (
	"net/http"

	"github.com/eser/ajan/httpfx"
)

const AuthHeader = "Authorization"

func AuthMiddleware() httpfx.Handler {
	return func(ctx *httpfx.Context) httpfx.Result {
		// FIXME(@eser) no need to check if the header is specified
		auth := ctx.Request.Header.Get(AuthHeader)
		if auth == "" {
			return ctx.Results.Error(http.StatusUnauthorized, []byte("Unauthorized"))
		}

		result := ctx.Next()

		return result
	}
}
