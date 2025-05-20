package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/eser/ajan/datafx"
	"github.com/eser/aya.is-services/pkg/api/adapters/storage"
	"github.com/eser/aya.is-services/pkg/api/business/users"
	"github.com/golang-jwt/jwt/v5"
)

// Adapter for users.OAuthService for GitHub

type GitHubOAuthService struct {
	DataRegistry *datafx.Registry
	ClientID     string
	ClientSecret string
	RedirectBase string
}

func NewGitHubOAuthService(dataRegistry *datafx.Registry) *GitHubOAuthService {
	return &GitHubOAuthService{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		RedirectBase: os.Getenv("GITHUB_REDIRECT_BASE"),
		DataRegistry: dataRegistry,
	}
}

func (g *GitHubOAuthService) InitiateOAuth(
	ctx context.Context,
	redirectURI string,
) (string, users.OAuthState, error) {
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

func (g *GitHubOAuthService) HandleOAuthCallback( //nolint:funlen
	ctx context.Context,
	code string,
	state string,
) (users.AuthResult, error) {
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
		Login  string `json:"login"`
		Name   string `json:"name"`
		Email  string `json:"email"`
		Avatar string `json:"avatar_url"`
		ID     int64  `json:"id"`
	}
	if err := json.NewDecoder(userResp.Body).Decode(&ghUser); err != nil {
		return users.AuthResult{}, err
	}

	// 3. Upsert user in DB
	repo, err := storage.NewRepositoryFromDefault(g.DataRegistry)
	if err != nil {
		return users.AuthResult{}, err
	}

	userId := fmt.Sprintf("github-%d", ghUser.ID)
	ghRemoteId := strconv.FormatInt(ghUser.ID, 10)
	now := time.Now()
	expiresAt := time.Now().Add(24 * time.Hour) //nolint:mnd

	user := users.User{
		Id:                  userId,
		Kind:                "regular",
		Name:                ghUser.Name,
		Email:               &ghUser.Email,
		Phone:               nil,
		GithubHandle:        &ghUser.Login,
		GithubRemoteId:      &ghRemoteId,
		BskyHandle:          nil,
		XHandle:             nil,
		IndividualProfileId: nil,
		CreatedAt:           now,
		UpdatedAt:           nil,
		DeletedAt:           nil,
	}
	err = repo.CreateUser(ctx, &user) // ignore error if already exists
	if err != nil {
		return users.AuthResult{}, err
	}

	// 4. Create session in DB
	session := users.Session{
		Id:                       fmt.Sprintf("sess-%s-%d", user.Id, time.Now().UnixNano()),
		Status:                   "active",
		OauthRequestState:        state,
		OauthRequestCodeVerifier: "", // PKCE not used here
		OauthRedirectUri:         nil,
		LoggedInUserId:           &user.Id,
		LoggedInAt:               &now,
		ExpiresAt:                &expiresAt,
		CreatedAt:                now,
		UpdatedAt:                nil,
	}
	err = repo.CreateSession(ctx, &session)
	if err != nil {
		return users.AuthResult{}, err
	}

	// 5. Issue JWT
	claims := users.JWTClaims{
		UserId:    user.Id,
		SessionId: session.Id,
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    claims.UserId,
		"session_id": claims.SessionId,
		"exp":        claims.ExpiresAt,
	})
	secret := os.Getenv("JWT_SECRET")
	tokenString, err := jwtToken.SignedString([]byte(secret))
	if err != nil {
		return users.AuthResult{}, err
	}

	return users.AuthResult{
		User:      &user,
		SessionId: session.Id,
		JWT:       tokenString,
	}, nil
}
