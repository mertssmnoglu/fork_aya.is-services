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
		GithubHandle *string `json:"githubHandle"`
		// GithubRemoteId      *string    `json:"githubRemoteId"`
		BskyHandle *string `json:"bskyHandle"`
		// BskyRemoteId        *string    `json:"bskyRemoteId"`
		XHandle *string `json:"xHandle"`
		// XRemoteId           *string    `json:"xRemoteId"`
		IndividualProfileId *string    `json:"individualProfileId"`
		CreatedAt           time.Time  `json:"createdAt"`
		UpdatedAt           *time.Time `json:"updatedAt"`
		DeletedAt           *time.Time `json:"deletedAt"`
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
