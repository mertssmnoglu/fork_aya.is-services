package users

import (
	"encoding/json"
	"time"

	"github.com/eser/aya.is-services/pkg/lib/vars"
)

func (s *User) MarshalJSON() ([]byte, error) {
	formatted := &struct {
		Id           string  `json:"id"`
		Kind         string  `json:"kind"`
		Name         string  `json:"name"`
		Email        *string `json:"email"`
		Phone        *string `json:"phone"`
		GithubHandle *string `json:"github_handle"`
		// GithubRemoteId      *string    `json:"github_remote_id"`
		BskyHandle *string `json:"bsky_handle"`
		// BskyRemoteId        *string    `json:"bsky_remote_id"`
		XHandle *string `json:"x_handle"`
		// XRemoteId           *string    `json:"x_remote_id"`
		IndividualProfileId *string    `json:"individual_profile_id"`
		CreatedAt           time.Time  `json:"created_at"`
		UpdatedAt           *time.Time `json:"updated_at"`
		DeletedAt           *time.Time `json:"deleted_at"`
	}{
		Id:           s.Id,
		Kind:         s.Kind,
		Name:         s.Name,
		Email:        vars.ToStringPtr(s.Email),
		Phone:        vars.ToStringPtr(s.Phone),
		GithubHandle: vars.ToStringPtr(s.GithubHandle),
		// GithubRemoteId:      &s.GithubRemoteId.String,
		BskyHandle: vars.ToStringPtr(s.BskyHandle),
		// BskyRemoteId:        &s.BskyRemoteId.String,
		XHandle: vars.ToStringPtr(s.XHandle),
		// XRemoteId:           &s.XRemoteId.String,
		IndividualProfileId: vars.ToStringPtr(s.IndividualProfileId),
		CreatedAt:           s.CreatedAt,
		UpdatedAt:           vars.ToTimePtr(s.UpdatedAt),
		DeletedAt:           vars.ToTimePtr(s.DeletedAt),
	}

	return json.Marshal(formatted) //nolint:wrapcheck
}
