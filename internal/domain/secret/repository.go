package secret

import (
	"context"
	"errors"
)

type Repository interface {
	GetSecretById(ctx context.Context, id string) (*Secret, error)
}

var (
	ErrorNotFound = errors.New("secret wasn't found")
)
