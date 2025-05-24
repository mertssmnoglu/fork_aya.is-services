//nolint:dupl
package http

import (
	"net/http"

	"github.com/eser/ajan/httpfx"
	"github.com/eser/ajan/logfx"
	"github.com/eser/aya.is-services/pkg/api/business/users"
	"github.com/eser/aya.is-services/pkg/lib/cursors"
)

func RegisterHttpRoutesForUsers(
	routes *httpfx.Router,
	logger *logfx.Logger,
	usersService *users.Service,
) {
	routes.
		Route("GET /{locale}/users", func(ctx *httpfx.Context) httpfx.Result {
			// get variables from path
			cursor := cursors.NewCursorFromRequest(ctx.Request)

			records, err := usersService.List(ctx.Request.Context(), cursor)
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
			idParam := ctx.Request.PathValue("id")

			record, err := usersService.GetById(ctx.Request.Context(), idParam)
			if err != nil {
				return ctx.Results.Error(http.StatusInternalServerError, []byte(err.Error()))
			}

			wrappedResponse := cursors.WrapResponseWithCursor(record, nil)

			return ctx.Results.Json(wrappedResponse)
		}).
		HasSummary("Get user by ID").
		HasDescription("Get user by ID.").
		HasResponse(http.StatusOK)

	// --- Auth endpoints ---
	routes.
		Route("GET /{locale}/auth/{authProvider}/login", func(ctx *httpfx.Context) httpfx.Result {
			// get auth provider from path
			authProviderName := ctx.Request.PathValue("authProvider")
			authProvider := usersService.AuthProvider(authProviderName)

			if authProvider == nil {
				return ctx.Results.Error(http.StatusNotFound, []byte("OAuth service not found"))
			}

			// Initiate OAuth flow
			redirectURI := ctx.Request.URL.Query().Get("redirect_uri")
			authURL, _, err := authProvider.InitiateOAuth(ctx.Request.Context(), redirectURI)
			if err != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					[]byte("OAuth initiation failed"),
				)
			}

			// Optionally set state in cookie/session
			return ctx.Results.Redirect(authURL)
		}).
		HasSummary("Auth Login").
		HasDescription("Redirects to auth provider login.").
		HasResponse(http.StatusFound)

	routes.
		Route("GET /{locale}/auth/{authProider}/callback", func(ctx *httpfx.Context) httpfx.Result {
			// get auth provider from path
			authProviderName := ctx.Request.PathValue("authProvider")
			authProvider := usersService.AuthProvider(authProviderName)

			if authProvider == nil {
				return ctx.Results.Error(http.StatusNotFound, []byte("OAuth service not found"))
			}

			url := ctx.Request.URL
			queryString := url.Query()
			code := queryString.Get("code")
			state := queryString.Get("state")

			result, err := authProvider.HandleOAuthCallback(ctx.Request.Context(), code, state)
			if err != nil {
				return ctx.Results.Error(http.StatusUnauthorized, []byte("OAuth callback failed"))
			}

			// Set JWT as cookie or return in response
			return ctx.Results.Json(map[string]any{
				"token": result.JWT,
				"user":  result.User,
			})
		}).
		HasSummary("Auth Callback").
		HasDescription("Handles auth provider callback and returns JWT.").
		HasResponse(http.StatusOK)

	routes.
		Route("POST /{locale}/auth/logout", func(ctx *httpfx.Context) httpfx.Result {
			// Invalidate session logic (optional, e.g., remove session from DB)
			return ctx.Results.Json(map[string]string{"status": "logged out"})
		}).
		HasSummary("Logout").
		HasDescription("Logs out the user.").
		HasResponse(http.StatusOK)
}
