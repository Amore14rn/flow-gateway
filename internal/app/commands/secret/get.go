package secret

import (
	"context"

	"github.com/davecgh/go-spew/spew"
	"gitlab.com/mildd/flow-gateway/internal/domain/secret"
)

type GetHandler struct {
	secretRepo secret.Repository
}

func NewGetHandler(secretRepo secret.Repository) GetHandler {
	return GetHandler{
		secretRepo: secretRepo,
	}
}

func (h GetHandler) GetHandler(ctx context.Context, id string) (err error) {

	userObj, err := h.secretRepo.GetSecretById(ctx, id)
	spew.Dump(userObj)
	if err != nil {
		return
	}
	return nil
}
