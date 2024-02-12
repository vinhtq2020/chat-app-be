package sequence

import (
	"go-service/internal/sequence/domain"
	"go-service/internal/sequence/repository"
	"go-service/internal/sequence/service"

	"gorm.io/gorm"
)

func NewSequenceService(db *gorm.DB) domain.SequenceService {
	repository := repository.NewSequenceRepository(db, "sequences", "sequence_no", "name")
	return service.NewSequenceUsecase(repository)
}
