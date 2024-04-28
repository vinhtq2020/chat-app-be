package usecase

import (
	"bufio"
	"context"
	"go-service/internal/autocomplete/aggregator/domain"
	"go-service/internal/autocomplete/model"
	"go-service/pkg/logger"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type LogAggregatorUsecase struct {
	logger     *logger.Logger
	repository domain.AggregatorRepository
}

func NewLogAggregatorService(repository domain.AggregatorRepository, logger *logger.Logger) *LogAggregatorUsecase {
	return &LogAggregatorUsecase{
		repository: repository,
		logger:     logger,
	}
}

func (l *LogAggregatorUsecase) AggregatedData(ctx context.Context) (int64, error) {

	// get root working directory
	root, err := os.Getwd()
	if err != nil {
		l.logger.LogError(err.Error(), nil)
		return -1, err
	}

	file, err := os.Open(root + "/logs/log_query.log")
	if err != nil {
		l.logger.LogError(err.Error(), nil)
		return -1, err
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)

	aggregatedDataMap := map[string]model.QueryCount{}
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		parts := strings.Split(line, " | ")
		query := parts[4]
		timeSearch := parts[0]
		timeParse, err := time.Parse(time.RFC3339, timeSearch)
		if err != nil {
			l.logger.LogError(err.Error(), nil)
			return -1, err
		}

		if v, exist := aggregatedDataMap[parts[4]]; exist {
			v.Prequency += 1
			aggregatedDataMap[parts[4]] = v
		} else {
			v = model.QueryCount{
				Query:     query,
				Time:      timeParse,
				Prequency: 1,
			}
			aggregatedDataMap[parts[4]] = v
		}

	}

	if err := scanner.Err(); err != nil {
		l.logger.LogError(err.Error(), nil)
		return -1, err
	}

	queryCount := []interface{}{}
	for _, v := range aggregatedDataMap {
		queryCount = append(queryCount, v)
	}

	if len(queryCount) == 0 {
		return 0, nil
	}

	// inserts aggregatedData bson
	res, err := l.repository.InsertMany(ctx, queryCount)
	if err != nil {
		l.logger.LogError(err.Error(), nil)
		return -1, err
	}

	// clear all data in log file.
	if err := os.Truncate(root+"/logs/log_query.log", 0); err != nil {
		l.logger.LogError(err.Error(), nil)
		return -1, err
	}
	return res, nil
}

func (l *LogAggregatorUsecase) All(ctx context.Context) ([]model.QueryCount, error) {
	filter := bson.M{
		"time": bson.M{
			"$gte": time.Now().Add(-1 * time.Hour).UTC(),
		},
	}
	return l.repository.All(ctx, filter)
}
