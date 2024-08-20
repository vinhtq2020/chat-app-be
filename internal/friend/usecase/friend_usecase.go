package usecase

import (
	"context"
	"fmt"
	"go-service/internal/friend/domain"
	friend_domain "go-service/internal/friend/domain"
	notification_domain "go-service/internal/notification/domain"
	user_domain "go-service/internal/user/domain"
	"go-service/pkg/convert"
	"go-service/pkg/logger"
	"time"

	"gorm.io/gorm"
)

const (
	addFriend        = "add_friend"
	yourAcceptFriend = "your_accept_friend"
	yourRejectFriend = "your_reject_friend"
	acceptFriend     = "accept_friend"
	rejectFriend     = "reject_friend"
)

type FriendUsecase struct {
	friendRqRepository  friend_domain.FriendRequestRepository
	userInfoRepository  user_domain.UserRepository
	friendRepository    friend_domain.FriendRepository
	notificationService notification_domain.NotificationService
	logger              *logger.Logger
}

func NewFriendService(friendRepository friend_domain.FriendRepository, friendrqRepository friend_domain.FriendRequestRepository, userInfoRepository user_domain.UserRepository, notificationService notification_domain.NotificationService, logger *logger.Logger) *FriendUsecase {
	return &FriendUsecase{
		friendRqRepository:  friendrqRepository,
		userInfoRepository:  userInfoRepository,
		friendRepository:    friendRepository,
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
	exist, err = u.friendRqRepository.Exist(ctx, friendRequest.RequesterId, friendRequest.RequesteeId, friend_domain.StatusPending)
	if err != nil || exist {
		return -1, err
	}

	// check already friended
	res, err := u.friendRepository.AlreadyFriended(ctx, friendRequest.RequesteeId, friendRequest.RequesterId)
	if err != nil || res > 0 {
		return -1, err
	}

	friendRequest.Status = friend_domain.StatusPending

	id, err := u.generateId(ctx)
	notificationId := "N-" + id
	if err != nil {
		return -1, err
	}

	friendRequest.Id = id
	friendRequest.NotificationId = notificationId

	res, err = u.friendRqRepository.InTransaction(ctx, func(ctx context.Context, db *gorm.DB) (int64, error) {
		res, err := u.friendRqRepository.Create(ctx, friendRequest)
		if err != nil {
			return res, err
		}
		res, err = u.friendRepository.Upsert(ctx, friend_domain.Friend{
			UserId1: friendRequest.RequesterId,
			UserId2: friendRequest.RequesteeId,
			Status:  domain.StatusPending.Value(),
		})
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
			RequestId: &friendRequest.Id,
			Name:      *requesterInfo.UserName,
			AvatarURL: requesterInfo.AvatarURL,
		}

		title := addFriend
		content := addFriend

		res, err = u.notificationService.Notify(ctx, func() string {
			return notificationId
		}, requester, title, content, []string{friendRequest.RequesteeId}, nil, notification_domain.NotificationAddFriend)
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

// use for reply request from requestee to requester
func (u *FriendUsecase) UpdateFriendRequest(ctx context.Context, requestId string, action string) (int64, error) {
	var receiverTitle, receiverContent, senderTitle, senderContent string

	// Load add friend request
	request, err := u.friendRqRepository.Load(ctx, requestId)
	if err != nil || request == nil {
		return 0, err
	}

	// check user is authorized to update status request
	userId := ctx.Value("userId").(string)
	if userId != request.RequesteeId {
		return -2, nil
	}

	res, err := u.friendRqRepository.InTransaction(ctx, func(ctx context.Context, db *gorm.DB) (int64, error) {
		if action == domain.AcceptAction {
			request.Status = friend_domain.StatusAccept
			receiverTitle = acceptFriend
			receiverContent = acceptFriend
			senderTitle = yourAcceptFriend
			senderContent = yourAcceptFriend
		} else {
			request.Status = friend_domain.StatusReject
			receiverTitle = rejectFriend
			receiverContent = rejectFriend
			senderTitle = yourRejectFriend
			senderContent = yourRejectFriend
		}

		request.UpdatedBy = request.RequesteeId
		request.UpdatedAt = time.Now()

		requestMap := convert.ToMapOmitEmpty(request)
		res, err := u.friendRqRepository.Patch(ctx, requestMap)
		if err != nil {
			return res, err
		}

		// check is already friend
		res, err = u.friendRepository.AlreadyFriended(ctx, request.RequesterId, request.RequesteeId)
		if err != nil {
			return res, err
		}

		// add friend relation to table friend
		res, err = u.friendRepository.Upsert(ctx, friend_domain.Friend{
			UserId1: request.RequesterId,
			UserId2: request.RequesteeId,
			Status:  request.Status.Value(),
		})
		if err != nil {
			return res, err
		}

		// update invite notification to success and notify this to sender(requestee)
		res, err = u.notificationService.UpdateNotify(ctx, request.NotificationId, senderTitle, senderContent, request.RequesteeId, notification_domain.NotificationInform, nil, false)
		if err != nil {
			return res, err
		}

		// notify to requester that requestee accepted
		userInfo, err := u.userInfoRepository.Load(ctx, request.RequesteeId)

		if err != nil {
			return res, err
		}

		requester := notification_domain.Requester{
			Id:        userInfo.Id,
			Name:      *userInfo.UserName,
			AvatarURL: userInfo.AvatarURL,
		}

		notificationId, err := u.generateId(ctx)
		if err != nil {
			return -1, err
		}

		res, err = u.notificationService.Notify(ctx, func() string {
			return notificationId
		}, requester, receiverTitle, receiverContent, []string{request.RequesterId}, nil, notification_domain.NotificationInform)

		return res, err
	})
	return res, err

}
func (u *FriendUsecase) Delete(ctx context.Context, friendId string) (int64, error) {
	panic("")
}
