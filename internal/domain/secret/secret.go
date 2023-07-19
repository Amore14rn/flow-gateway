package secret

import (
	"github.com/google/uuid"
	"gitlab.com/mildd/flow-gateway/pkg/client/struct/changed"
)

type Secret struct {
	changed.Changed
	UUID   uuid.UUID `db:"uuid"`
	Secret string    `db:"secret"`
}

func NewSecret(opts ...Option) (*Secret, error) {
	user := &Secret{UUID: uuid.New(), Changed: map[string]interface{}{}}
	for _, opt := range opts {
		err := opt(user)
		if err != nil {
			return nil, err
		}
	}
	return user, nil
}

func NewEmptySecret() *Secret {
	return &Secret{Changed: map[string]interface{}{}}
}
