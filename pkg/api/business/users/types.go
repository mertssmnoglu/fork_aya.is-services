package users

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type RecordID string

type RecordIDGenerator func() RecordID

func DefaultIDGenerator() RecordID {
	return RecordID(ulid.Make().String())
}

type User struct {
	CreatedAt      time.Time `json:"created_at"`
	Email          *string   `json:"email"`
	Phone          *string   `json:"phone"`
	GithubHandle   *string   `json:"github_handle"`
	GithubRemoteId *string   `json:"github_remote_id"`
	BskyHandle     *string   `json:"bsky_handle"`
	// BskyRemoteId        *string    `json:"bsky_remote_id"`
	XHandle *string `json:"x_handle"`
	// XRemoteId           *string    `json:"x_remote_id"`
	IndividualProfileId *string    `json:"individual_profile_id"`
	UpdatedAt           *time.Time `json:"updated_at"`
	DeletedAt           *time.Time `json:"deleted_at"`
	Id                  string     `json:"id"`
	Kind                string     `json:"kind"`
	Name                string     `json:"name"`
}

type Session struct {
	CreatedAt                time.Time  `json:"created_at"`
	OauthRedirectUri         *string    `json:"oauth_redirect_uri"`
	LoggedInUserId           *string    `json:"logged_in_user_id"`
	LoggedInAt               *time.Time `json:"logged_in_at"`
	ExpiresAt                *time.Time `json:"expires_at"`
	UpdatedAt                *time.Time `json:"updated_at"`
	Id                       string     `json:"id"`
	Status                   string     `json:"status"`
	OauthRequestState        string     `json:"oauth_request_state"`
	OauthRequestCodeVerifier string     `json:"oauth_request_code_verifier"`
}

// --- OAuth & Auth types ---

type OAuthProvider string

const (
	OAuthProviderGitHub OAuthProvider = "github"
)

type OAuthState struct {
	State       string
	RedirectURI string
}

type GitHubUserInfo struct {
	ID     string
	Login  string
	Name   string
	Email  string
	Avatar string
}

type AuthResult struct {
	User      *User
	SessionId string
	JWT       string
}

type JWTClaims struct {
	UserId    string
	SessionId string
	ExpiresAt int64
}
