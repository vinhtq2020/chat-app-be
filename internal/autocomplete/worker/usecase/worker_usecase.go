package usecase

import (
	"context"
	"go-service/internal/autocomplete/model"
	"go-service/internal/autocomplete/worker/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type WorkerUsecase struct {
	repository domain.WorkerRepository
}

func NewWorkerService(repository domain.WorkerRepository) *WorkerUsecase {
	return &WorkerUsecase{repository: repository}
}

func (s *WorkerUsecase) CreateTries(ctx context.Context, queryCounts []model.QueryCount) (int64, error) {
	tries := model.NewTrie()
	for _, queryCount := range queryCounts {
		tries.Insert(queryCount.Query, queryCount.Prequency)
	}

	res, err := s.repository.InTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		filter := bson.M{}
		_, err := s.repository.Delete(ctx, filter)
		if err != nil {
			return -1, err
		}
		return s.repository.Insert(ctx, tries)

	})
	return res.(int64), err
}

func (s *WorkerUsecase) LoadTries(ctx context.Context) (*model.Trie, error) {
	return s.repository.Load(ctx, bson.M{})

}
