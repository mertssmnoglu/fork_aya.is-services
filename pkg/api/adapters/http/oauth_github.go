package http

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/eser/ajan/datafx"
	"github.com/eser/aya.is-services/pkg/api/adapters/storage"
	"github.com/eser/aya.is-services/pkg/api/business/users"
	"github.com/golang-jwt/jwt/v5"
)

// Adapter for users.OAuthService for GitHub

type GitHubOAuthService struct {
	ClientID     string
	ClientSecret string
	RedirectBase string
	DataRegistry *datafx.Registry
}

func NewGitHubOAuthService(dataRegistry *datafx.Registry) *GitHubOAuthService {
	return &GitHubOAuthService{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		RedirectBase: os.Getenv("GITHUB_REDIRECT_BASE"),
		DataRegistry: dataRegistry,
	}
}

func (g *GitHubOAuthService) InitiateOAuth(ctx context.Context, redirectURI string) (string, users.OAuthState, error) {
	state := fmt.Sprintf("%d", time.Now().UnixNano()) // TODO: use secure random
	u := url.URL{
		Scheme: "https",
		Host:   "github.com",
		Path:   "/login/oauth/authorize",
	}
	q := u.Query()
	q.Set("client_id", g.ClientID)
	q.Set("redirect_uri", redirectURI)
	q.Set("state", state)
	q.Set("scope", "read:user user:email")
	u.RawQuery = q.Encode()
	return u.String(), users.OAuthState{State: state, RedirectURI: redirectURI}, nil
}

func (g *GitHubOAuthService) HandleOAuthCallback(ctx context.Context, code string, state string) (users.AuthResult, error) {
	// 1. Exchange code for access token
	tokenResp, err := http.PostForm("https://github.com/login/oauth/access_token", url.Values{
		"client_id":     {g.ClientID},
		"client_secret": {g.ClientSecret},
		"code":          {code},
	})
	if err != nil {
		return users.AuthResult{}, err
	}
	defer tokenResp.Body.Close()
	body, _ := io.ReadAll(tokenResp.Body)
	vals, _ := url.ParseQuery(string(body))
	accessToken := vals.Get("access_token")
	if accessToken == "" {
		return users.AuthResult{}, fmt.Errorf("failed to get access token")
	}

	// 2. Fetch user info from GitHub
	req, _ := http.NewRequestWithContext(ctx, "GET", "https://api.github.com/user", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	userResp, err := http.DefaultClient.Do(req)
	if err != nil {
		return users.AuthResult{}, err
	}
	defer userResp.Body.Close()
	var ghUser struct {
		ID     int64  `json:"id"`
		Login  string `json:"login"`
		Name   string `json:"name"`
		Email  string `json:"email"`
		Avatar string `json:"avatar_url"`
	}
	if err := json.NewDecoder(userResp.Body).Decode(&ghUser); err != nil {
		return users.AuthResult{}, err
	}

	// 3. Upsert user in DB
	repo, err := storage.NewRepositoryFromDefault(g.DataRegistry)
	if err != nil {
		return users.AuthResult{}, err
	}
	queries := repo.GetQueries()
	userID := fmt.Sprintf("github-%d", ghUser.ID)
	userParams := storage.CreateUserParams{
		Id:             userID,
		Kind:           "individual",
		Name:           ghUser.Name,
		Email:          sql.NullString{String: ghUser.Email, Valid: ghUser.Email != ""},
		GithubHandle:   sql.NullString{String: ghUser.Login, Valid: ghUser.Login != ""},
		GithubRemoteId: sql.NullString{String: fmt.Sprintf("%d", ghUser.ID), Valid: true},
	}
	_, _ = queries.CreateUser(ctx, userParams) // ignore error if already exists
	userRow, err := queries.GetUserById(ctx, storage.GetUserByIdParams{Id: userID})
	if err != nil {
		return users.AuthResult{}, err
	}
	user := &users.User{
		Id:           userRow.Id,
		Kind:         userRow.Kind,
		Name:         userRow.Name,
		Email:        &userRow.Email.String,
		GithubHandle: &userRow.GithubHandle.String,
		CreatedAt:    userRow.CreatedAt,
		UpdatedAt:    nil, // can map if needed
	}

	// 4. Create session in DB
	sessionID := fmt.Sprintf("sess-%s-%d", user.Id, time.Now().UnixNano())
	sessionParams := storage.CreateSessionParams{
		Id:                       sessionID,
		Status:                   "active",
		OauthRequestState:        state,
		OauthRequestCodeVerifier: "", // PKCE not used here
		OauthRedirectUri:         sql.NullString{Valid: false},
		LoggedInUserId:           sql.NullString{String: user.Id, Valid: true},
		LoggedInAt:               sql.NullTime{Time: time.Now(), Valid: true},
		ExpiresAt:                sql.NullTime{Time: time.Now().Add(24 * time.Hour), Valid: true},
	}
	_, err = queries.CreateSession(ctx, sessionParams)
	if err != nil {
		return users.AuthResult{}, err
	}

	// 5. Issue JWT
	claims := users.JWTClaims{
		UserID:    user.Id,
		SessionID: sessionID,
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    claims.UserID,
		"session_id": claims.SessionID,
		"exp":        claims.ExpiresAt,
	})
	secret := os.Getenv("JWT_SECRET")
	tokenString, err := jwtToken.SignedString([]byte(secret))
	if err != nil {
		return users.AuthResult{}, err
	}

	return users.AuthResult{
		User:      user,
		SessionID: sessionID,
		JWT:       tokenString,
	}, nil
}
