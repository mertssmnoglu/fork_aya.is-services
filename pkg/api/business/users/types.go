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
	CreatedAt    time.Time `json:"created_at"`
	Email        *string   `json:"email"`
	Phone        *string   `json:"phone"`
	GithubHandle *string   `json:"github_handle"`
	// GithubRemoteId      *string    `json:"github_remote_id"`
	BskyHandle *string `json:"bsky_handle"`
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
