package http

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/eser/ajan/datafx"
	"github.com/eser/ajan/httpfx"
	"github.com/eser/aya.is-services/pkg/api/adapters/storage"
	"github.com/golang-jwt/jwt/v5"
)

const AuthHeader = "Authorization"

func AuthMiddleware(dataRegistry *datafx.Registry) httpfx.Handler {
	return func(ctx *httpfx.Context) httpfx.Result {
		// FIXME(@eser) no need to check if the header is specified
		auth := ctx.Request.Header.Get(AuthHeader)

		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			return ctx.Results.Error(http.StatusUnauthorized, []byte("Unauthorized"))
		}

		tokenStr := strings.TrimPrefix(auth, "Bearer ")
		secret := os.Getenv("JWT_SECRET")
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			return ctx.Results.Error(http.StatusUnauthorized, []byte("Invalid token"))
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return ctx.Results.Error(http.StatusUnauthorized, []byte("Invalid claims"))
		}

		sessionId, _ := claims["session_id"].(string)
		if sessionId == "" {
			return ctx.Results.Error(http.StatusUnauthorized, []byte("No session"))
		}

		repository, err := storage.NewRepositoryFromDefault(dataRegistry)
		if err != nil {
			return ctx.Results.Error(http.StatusInternalServerError, []byte(err.Error()))
		}

		// Load session from DB
		session, err := repository.GetSessionById(ctx.Request.Context(), sessionId)
		if err != nil || session.Status != "active" {
			return ctx.Results.Error(http.StatusUnauthorized, []byte("Session invalid"))
		}

		// Update logged_in_at
		_ = repository.UpdateSessionLoggedInAt(ctx.Request.Context(), sessionId, time.Now())

		result := ctx.Next()

		return result
	}
}
