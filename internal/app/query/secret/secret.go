package secret

import "context"

type SecretHandler struct {
	readModel SecretModel
}

func NewSecretHandler(readModel SecretModel) SecretHandler {
	if readModel == nil {
		panic("nil readModel")
	}

	return SecretHandler{readModel: readModel}
}

type SecretModel interface {
	GetSecretQueryModel(ctx context.Context, id string) (*Secret, error)
}

func (h SecretHandler) Handle(ctx context.Context, id string) (*Secret, error) {
	return h.readModel.GetSecretQueryModel(ctx, id)
}
