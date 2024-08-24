package service

import (
	"context"
	"go-service/internal/sequence/domain"

	"gorm.io/gorm"
)

type SequenceUsecase struct {
	repository domain.SequenceRepository
}

func NewSequenceUsecase(repository domain.SequenceRepository) *SequenceUsecase {
	return &SequenceUsecase{
		repository: repository,
	}
}

func (u *SequenceUsecase) Next(ctx context.Context, module string) (int64, error) {
	return u.repository.InTransaction(ctx, func(ctx context.Context, tx *gorm.DB) (int64, error) {
		res, err := u.repository.Next(ctx, module)
		if err != nil || res <= 0 {
			return res, err
		}
		return u.repository.GetSequence(ctx, module)
	})
}
