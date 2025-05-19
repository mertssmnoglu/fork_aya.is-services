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
	oauthService users.OAuthService, // Injected dependency
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

	// --- GitHub OAuth endpoints ---
	routes.Route("GET /auth/github/login", func(ctx *httpfx.Context) httpfx.Result {
		// Initiate OAuth flow
		redirectURI := ctx.Request.URL.Query().Get("redirect_uri")
		authURL, _, err := oauthService.InitiateOAuth(ctx.Request.Context(), redirectURI)
		if err != nil {
			return ctx.Results.Error(http.StatusInternalServerError, []byte("OAuth initiation failed"))
		}
		// Optionally set state in cookie/session
		return ctx.Results.Redirect(authURL)
	}).HasSummary("GitHub OAuth Login").HasDescription("Redirects to GitHub OAuth login.")

	routes.Route("GET /auth/github/callback", func(ctx *httpfx.Context) httpfx.Result {
		code := ctx.Request.URL.Query().Get("code")
		state := ctx.Request.URL.Query().Get("state")
		result, err := oauthService.HandleOAuthCallback(ctx.Request.Context(), code, state)
		if err != nil {
			return ctx.Results.Error(http.StatusUnauthorized, []byte("OAuth callback failed"))
		}
		// Set JWT as cookie or return in response
		return ctx.Results.Json(map[string]interface{}{
			"token": result.JWT,
			"user":  result.User,
		})
	}).HasSummary("GitHub OAuth Callback").HasDescription("Handles GitHub OAuth callback and returns JWT.")

	routes.Route("POST /auth/logout", func(ctx *httpfx.Context) httpfx.Result {
		// Invalidate session logic (optional, e.g., remove session from DB)
		return ctx.Results.Json(map[string]string{"status": "logged out"})
	}).HasSummary("Logout").HasDescription("Logs out the user.")
}
