package usecase

import (
	"context"

	"hypefast.io/services/config"
	shorterModel "hypefast.io/services/domain/shorter/model"
)

type ShorterUseCase interface {
	GetURL(key string) (string, error)
	SetURL(payload shorterModel.Payload) (interface{}, error)
}
type ShorterUseCaseImpl struct {
	Ctx  context.Context
	Conf config.Configuration
}

func NewShorterUseCase(ctx context.Context, conf config.Configuration) ShorterUseCase {
	return &ShorterUseCaseImpl{
		Ctx:  ctx,
		Conf: conf,
	}
}
