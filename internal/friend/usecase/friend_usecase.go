package usecase

import (
	"context"
	"fmt"
	friend_domain "go-service/internal/friend/domain"
	domain_notification "go-service/internal/notification/notification_domain"
	"go-service/pkg/logger"
	"math/rand"
	"time"
)

const (
	addFriendContent = "%s đã gửi cho bạn lời mời kết bạn."
)

type FriendUsecase struct {
	friendRqRepository  friend_domain.FriendRequestRepository
	notificationService domain_notification.NotificationService
	logger              *logger.Logger
}

func NewFriendService(friendrqRepository friend_domain.FriendRequestRepository, notificationService domain_notification.NotificationService, logger *logger.Logger) *FriendUsecase {
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

	addFriendRequest := friend_domain.FriendRequest{
		Uid1:      userId,
		Uid2:      friendId,
		CreatedAt: time.Now(),
		CreatedBy: userId,
		UpdatedAt: time.Now(),
		UpdatedBy: userId,
	}

	res, err := u.friendRqRepository.Create(ctx, addFriendRequest)
	if err != nil {
		return res, err
	}
	_, err = u.notificationService.Notify(ctx, u.generateId, userId, addFriendContent, []string{friendId})
	return res, nil
}

func (u *FriendUsecase) generateId() string {
	return fmt.Sprintf("F-%v", rand.Intn(6))
}

func (u *FriendUsecase) Patch(ctx context.Context, friendId string, status friend_domain.FriendRequestStatus) (int64, error) {
	panic("")
}
func (u *FriendUsecase) Delete(ctx context.Context, friendId string) (int64, error) {
	panic("")
}
