package usecase

import (
	"context"
	"go-service/internal/friend/domain"
	domain_user "go-service/internal/user/domain"
	"time"
)

type FriendUsecase struct {
	friendRqRepository domain.FriendRequestRepository
	userRepository     domain_user.UserRepository
}

func NewFriendService(friendrqRepository domain.FriendRequestRepository) *FriendUsecase {
	return &FriendUsecase{
		friendRqRepository: friendrqRepository,
	}
}

func (u *FriendUsecase) SendFriendRequest(ctx context.Context, userId string, friendId string) (int64, error) {
	exist, err := u.userRepository.Exist(ctx, friendId)
	if err != nil || !exist {
		return 0, err
	}

	exist, err = u.friendRqRepository.Exist(ctx, userId, friendId)
	if err != nil || exist {
		return -1, err
	}

	addFriendRequest := domain.FriendRequest{
		Uid1:      userId,
		Uid2:      friendId,
		CreatedAt: time.Now(),
		CreatedBy: userId,
		UpdatedAt: time.Now(),
		UpdatedBy: userId,
	}

	res, err := u.friendRqRepository.Create(ctx, addFriendRequest)
	if err != nil {
		return -1, err
	}
	return res, nil
}

func (u *FriendUsecase) Patch(ctx context.Context, friendId string, status domain.FriendRequestStatus) (int64, error) {
	panic("")
}
func (u *FriendUsecase) Delete(ctx context.Context, friendId string) (int64, error) {
	panic("")
}
