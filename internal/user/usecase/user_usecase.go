package usecase

import user_domain "go-service/internal/user/domain"

type UserUsecase struct {
	repo user_domain.UserRepository
}

func NewUserUsecase(repo user_domain.UserRepository) *UserUsecase {
	return &UserUsecase{repo: repo}
}
