package profiles

import (
	"encoding/json"
	"time"

	"github.com/eser/aya.is-services/pkg/lib/vars"
)

func (s *Profile) MarshalJSON() ([]byte, error) {
	formatted := &struct {
		Id                string  `json:"id"`
		Slug              string  `json:"slug"`
		Kind              string  `json:"kind"`
		CustomDomain      *string `json:"custom_domain"`
		ProfilePictureUri *string `json:"profile_picture_uri"`
		Pronouns          *string `json:"pronouns"`
		Title             string  `json:"title"`
		Description       string  `json:"description"`
		// Properties        *string    `json:"properties"`
		CreatedAt time.Time  `json:"created_at"`
		UpdatedAt *time.Time `json:"updated_at"`
		DeletedAt *time.Time `json:"deleted_at"`
	}{
		Id:                s.Id,
		Slug:              s.Slug,
		Kind:              s.Kind,
		CustomDomain:      vars.ToStringPtr(s.CustomDomain),
		ProfilePictureUri: vars.ToStringPtr(s.ProfilePictureUri),
		Pronouns:          vars.ToStringPtr(s.Pronouns),
		Title:             s.Title,
		Description:       s.Description,
		// Properties:        vars.ToStringPtr(s.Properties),
		CreatedAt: s.CreatedAt,
		UpdatedAt: vars.ToTimePtr(s.UpdatedAt),
		DeletedAt: vars.ToTimePtr(s.DeletedAt),
	}

	return json.Marshal(formatted) //nolint:wrapcheck
}
