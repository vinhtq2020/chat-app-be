package usecase

import "go-service/internal/user/user_domain"

type UserUsecase struct {
	repo user_domain.UserRepository
}

func NewUserUsecase(repo user_domain.UserRepository) *UserUsecase {
	return &UserUsecase{repo: repo}
}
