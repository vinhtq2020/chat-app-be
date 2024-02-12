package usecase

import "go-service/internal/friend/domain"

type FriendService struct {
	friendRepository domain.FriendRepository
}

func NewFriendService(friendRepository domain.FriendRepository) *FriendService {
	return &FriendService{
		friendRepository: friendRepository,
	}
}
