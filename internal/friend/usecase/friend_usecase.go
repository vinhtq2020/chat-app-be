package usecase

import (
	"context"
	"fmt"
	"go-service/internal/friend/domain"
	friend_domain "go-service/internal/friend/domain"
	notification_domain "go-service/internal/notification/domain"
	user_domain "go-service/internal/user/domain"
	"go-service/pkg/logger"

	"gorm.io/gorm"
)

const (
	addFriendTitle   = "%s đã gửi cho bạn lời mời kết bạn."
	addFriendContent = "%s đã gửi cho bạn lời mời kết bạn."
)

type FriendUsecase struct {
	friendRqRepository  friend_domain.FriendRequestRepository
	userInfoRepository  user_domain.UserRepository
	notificationService notification_domain.NotificationService
	logger              *logger.Logger
}

func NewFriendService(friendrqRepository friend_domain.FriendRequestRepository, userInfoRepository user_domain.UserRepository, notificationService notification_domain.NotificationService, logger *logger.Logger) *FriendUsecase {
	return &FriendUsecase{
		friendRqRepository:  friendrqRepository,
		userInfoRepository:  userInfoRepository,
		logger:              logger,
		notificationService: notificationService,
	}
}

func (u *FriendUsecase) SendFriendRequest(ctx context.Context, friendRequest domain.FriendRequest) (int64, error) {
	var err error
	defer func() {
		if err != nil {
			u.logger.LogError(err.Error(), nil)
		}
	}()

	// check friend Id exist
	exist, err := u.userInfoRepository.Exist(ctx, friendRequest.RequesterId)
	if err != nil || !exist {
		return 0, err
	}

	// check friend request exist
	exist, err = u.friendRqRepository.Exist(ctx, friendRequest.RequesterId, friendRequest.RequesteeId)
	if err != nil || exist {
		return -1, err
	}

	friendRequest.Status = friend_domain.StatusPending

	id, err := u.generateId(ctx)
	if err != nil {
		return -1, err
	}

	friendRequest.Id = id

	res, err := u.friendRqRepository.InTransaction(ctx, func(ctx context.Context, db *gorm.DB) (int64, error) {
		res, err := u.friendRqRepository.Create(ctx, friendRequest)
		if err != nil {
			return res, err
		}

		// Notify
		requesterInfo, err := u.userInfoRepository.Load(ctx, friendRequest.RequesterId)
		if err != nil {
			return res, nil
		}

		requester := notification_domain.Requester{
			Id:        friendRequest.RequesterId,
			Name:      *requesterInfo.UserName,
			AvatarURL: requesterInfo.AvatarURL,
		}

		title := fmt.Sprintf(addFriendTitle, *requesterInfo.UserName)
		content := fmt.Sprintf(addFriendContent, *requesterInfo.UserName)

		res, err = u.notificationService.Notify(ctx, func() string {
			return id
		}, requester, title, content, []string{friendRequest.RequesteeId}, nil, notification_domain.NotificationInvite)
		if err != nil {
			return res, err
		}
		return res, nil
	})

	if err != nil {
		return res, err
	}
	return res, nil
}

func (u *FriendUsecase) generateId(ctx context.Context) (string, error) {
	total, err := u.friendRqRepository.Total(ctx)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("F-%v", total), nil
}

func (u *FriendUsecase) Patch(ctx context.Context, friendRequest map[string]interface{}) (int64, error) {
	panic("")
}
func (u *FriendUsecase) Delete(ctx context.Context, friendId string) (int64, error) {
	panic("")
}
