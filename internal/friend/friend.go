package friend

import (
	"go-service/internal/friend/delivery"
	"go-service/internal/friend/domain"
	"go-service/internal/friend/repository"
	"go-service/internal/friend/usecase"
)

func NewFriendHandler() *domain.FriendTransport {
	repo := repository.NewRequestFriendRepository()
	sv := usecase.NewFriendService()
	return &delivery.NewFriendHandler(sv)
}
