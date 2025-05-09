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
		CustomDomain      *string `json:"customDomain"`
		ProfilePictureUri *string `json:"profilePictureUri"`
		Pronouns          *string `json:"pronouns"`
		Title             string  `json:"title"`
		Description       string  `json:"description"`
		// Properties        *string    `json:"properties"`
		CreatedAt time.Time  `json:"createdAt"`
		UpdatedAt *time.Time `json:"updatedAt"`
		DeletedAt *time.Time `json:"deletedAt"`
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
