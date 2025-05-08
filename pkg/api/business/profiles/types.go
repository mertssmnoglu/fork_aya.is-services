package profiles

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type RecordID string

type RecordIDGenerator func() RecordID

func DefaultIDGenerator() RecordID {
	return RecordID(ulid.Make().String())
}

type ExternalPost struct {
	Id        string     `json:"id"`
	Content   string     `json:"content"`
	Permalink string     `json:"permalink"`
	CreatedAt *time.Time `json:"created_at"` //nolint:tagliatelle
}
