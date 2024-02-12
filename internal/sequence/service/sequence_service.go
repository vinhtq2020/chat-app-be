package service

import (
	"go-service/internal/sequence/domain"

	"github.com/gin-gonic/gin"
)

type SequenceUsecase struct {
	repository domain.SequenceRepository
}

func NewSequenceUsecase(repository domain.SequenceRepository) *SequenceUsecase {
	return &SequenceUsecase{
		repository: repository,
	}
}

func (u *SequenceUsecase) Next(c *gin.Context, module string) (int64, error) {
	res, err := u.repository.Next(c, module)
	if err != nil || res <= 0 {
		return res, err
	}
	return u.repository.GetSequence(c, module)
}
