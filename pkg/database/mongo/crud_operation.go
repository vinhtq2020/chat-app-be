package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

func InTransaction(ctx context.Context, client *mongo.Client, callback func(sessCtx mongo.SessionContext) (interface{}, error)) (interface{}, error) {
	session, err := client.StartSession()
	if err != nil {
		return nil, err
	}

	defer session.EndSession(ctx)
	res, err := session.WithTransaction(ctx, callback)
	if err != nil {
		return nil, err
	}
	return res, nil
}
