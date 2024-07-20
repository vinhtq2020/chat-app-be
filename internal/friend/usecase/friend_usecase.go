package usecase

import (
	"context"
	"fmt"
	friend_domain "go-service/internal/friend/domain"
	notification_domain "go-service/internal/notification/domain"
	"go-service/pkg/logger"
)

const (
	addFriendContent = "%s đã gửi cho bạn lời mời kết bạn."
)

type FriendUsecase struct {
	friendRqRepository  friend_domain.FriendRequestRepository
	notificationService notification_domain.NotificationService
	logger              *logger.Logger
}

func NewFriendService(friendrqRepository friend_domain.FriendRequestRepository, notificationService notification_domain.NotificationService, logger *logger.Logger) *FriendUsecase {
	return &FriendUsecase{
		friendRqRepository:  friendrqRepository,
		logger:              logger,
		notificationService: notificationService,
	}
}

func (u *FriendUsecase) SendFriendRequest(ctx context.Context, userId string, friendId string) (int64, error) {
	var err error
	defer func() {
		if err != nil {
			u.logger.LogError(err.Error(), nil)
		}
	}()
	exist, err := u.friendRqRepository.UserExist(ctx, friendId)
	if err != nil || !exist {
		return 0, err
	}

	exist, err = u.friendRqRepository.Exist(ctx, userId, friendId)
	if err != nil || exist {
		return -1, err
	}

	// addFriendRequest := friend_domain.FriendRequest{
	// 	RequesterId: userId,
	// 	RequesteeId: friendId,
	// 	CreatedAt:   time.Now(),
	// 	CreatedBy:   userId,
	// 	Status:      &friend_domain.StatusPending,
	// 	UpdatedAt:   time.Now(),
	// 	UpdatedBy:   userId,
	// }

	// res, err := u.friendRqRepository.Create(ctx, addFriendRequest)
	// if err != nil {
	// 	return res, err
	// }
	id, err := u.generateId(ctx)
	if err != nil {
		return -1, err
	}
	_, err = u.notificationService.Notify(ctx, func() string {
		return id
	}, userId, addFriendContent, []string{friendId})
	return 1, nil
}

func (u *FriendUsecase) generateId(ctx context.Context) (string, error) {
	total, err := u.friendRqRepository.Total(ctx)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("F-%v", total), nil
}

func (u *FriendUsecase) Patch(ctx context.Context, friendId string, status friend_domain.FriendRequestStatus) (int64, error) {
	panic("")
}
func (u *FriendUsecase) Delete(ctx context.Context, friendId string) (int64, error) {
	panic("")
}
