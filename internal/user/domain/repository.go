package domain

import "context"

type UserRepository interface {
	Exist(ctx context.Context, id string) (bool, error)
	Create(ctx context.Context, user User) (int64, error)
}
