package usecase

import (
	"context"
	"fmt"
	"go-service/internal/friend/domain"
	friend_domain "go-service/internal/friend/domain"
	notification_domain "go-service/internal/notification/domain"
	sequence_domain "go-service/internal/sequence/domain"
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
	userInfoRepository  user_domain.UserRepository
	friendRepository    friend_domain.FriendRepository
	notificationService notification_domain.NotificationService
	sequence            sequence_domain.SequenceService
	logger              *logger.Logger
}

func NewFriendService(friendRepository friend_domain.FriendRepository, userInfoRepository user_domain.UserRepository, notificationService notification_domain.NotificationService, sequence sequence_domain.SequenceService, logger *logger.Logger) *FriendUsecase {
	return &FriendUsecase{
		userInfoRepository:  userInfoRepository,
		friendRepository:    friendRepository,
		logger:              logger,
		notificationService: notificationService,
		sequence:            sequence,
	}
}

func (u *FriendUsecase) Create(ctx context.Context, userId string, friendId string, relation string) (int64, error) {
	var err error
	defer func() {
		if err != nil {
			u.logger.LogError(err.Error(), nil)
		}
	}()

	// check not friended
	friendRelation, err := u.friendRepository.Load(ctx, userId, friendId, domain.FriendRelation)
	if err != nil {
		return -1, err
	}

	if friendRelation.Status == acceptFriend {
		return -1, nil
	}

	id, err := u.generateId(ctx)
	notificationId := "N-" + id
	if err != nil {
		return -1, err
	}

	res, err := u.friendRepository.InTransaction(ctx, func(ctx context.Context, db *gorm.DB) (int64, error) {
		// add friend request
		relation := friend_domain.NewRelation(userId, friendId, userId, domain.FriendRelation, domain.StatusPending)
		res, err := u.friendRepository.Upsert(ctx, relation)
		if err != nil {
			return res, err
		}

		// Notify
		requesterInfo, err := u.userInfoRepository.Load(ctx, userId)
		if err != nil {
			return res, nil
		}

		requester := notification_domain.Requester{
			Id:        userId,
			Name:      *requesterInfo.UserName,
			AvatarURL: requesterInfo.AvatarURL,
		}

		title := addFriend
		content := addFriend

		res, err = u.notificationService.Notify(ctx, func() string {
			return notificationId
		}, requester, title, content, []string{friendId}, nil, notification_domain.NotificationAddFriend)
		if err != nil {
			return res, err
		}
		return res, nil
	})

	if err != nil {
		return -1, err
	}
	return res, nil
}

func (u *FriendUsecase) generateId(ctx context.Context) (string, error) {
	count, err := u.sequence.Next(ctx, "relation-request")
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("F-%v", count), nil
}

// need inform to user that request cancel or not ?
func (u *FriendUsecase) Cancel(ctx context.Context, userId string, friendId string) (int64, error) {
	friendRelation, err := u.friendRepository.Load(ctx, userId, friendId, domain.FriendRelation)
	if friendRelation == nil || err != nil {
		return 0, err
	}

	// check user is friend request sender
	if userId != friendRelation.UpdatedBy {
		return -2, nil
	}

	// check is friend request pending
	if friendRelation.Status != domain.StatusPending {
		return -1, nil
	}

	res, err := u.friendRepository.InTransaction(ctx, func(ctx context.Context, db *gorm.DB) (int64, error) {
		// cancel request
		friendRelation.Status = friend_domain.StatusCancel
		friendRelation.UpdatedBy = userId
		friendRelation.UpdatedAt = time.Now()

		requestMap := convert.ToMapOmitEmpty(friendRelation)
		res, err := u.friendRepository.Patch(ctx, requestMap)
		if err != nil {
			return res, err
		}

		// delete notification invite in server
		notifications, err := u.notificationService.Search(ctx, notification_domain.NotificationFilter{
			SubscriberId: &friendId,
		})
		if err != nil || len(notifications) == 0 {
			return 0, err
		}
		res, err = u.notificationService.Delete(ctx, notifications[0].Id)
		if err != nil {
			return res, err
		}

		return 1, nil
	})

	return res, err
}

// use for reply request from requestee to requester. Choose accept or reject
func (u *FriendUsecase) Response(ctx context.Context, userId string, friendId string, action string) (int64, error) {
	var receiverTitle, receiverContent, senderTitle, senderContent string

	// check friend request is existed
	friendRelation, err := u.friendRepository.Load(ctx, userId, friendId, domain.FriendRelation)
	if err != nil || friendRelation == nil {
		return 0, err
	}

	// check is pending
	if friendRelation.Status != domain.StatusPending {
		return -1, nil
	}

	// check user is authorized to update status request. If is Sender, updatedBy now not replier (current is userId)
	if userId != friendRelation.UpdatedBy {
		return -2, nil
	}

	res, err := u.friendRepository.InTransaction(ctx, func(ctx context.Context, db *gorm.DB) (int64, error) {
		if action == domain.AcceptAction {
			friendRelation.Status = friend_domain.StatusAccept
			receiverTitle = acceptFriend
			receiverContent = acceptFriend
			senderTitle = yourAcceptFriend
			senderContent = yourAcceptFriend
		} else if action == domain.RejectAction {
			friendRelation.Status = friend_domain.StatusReject
			receiverTitle = rejectFriend
			receiverContent = rejectFriend
			senderTitle = yourRejectFriend
			senderContent = yourRejectFriend
		}

		friendRelation.UpdatedBy = userId
		friendRelation.UpdatedAt = time.Now()

		requestMap := convert.ToMapOmitEmpty(friendRelation)
		res, err := u.friendRepository.Patch(ctx, requestMap)
		if err != nil {
			return res, err
		}

		// update invite notification to success and notify this to sender(requestee)
		notifications, err := u.notificationService.Search(ctx, notification_domain.NotificationFilter{
			SubscriberId: &friendId,
		})
		if err != nil || len(notifications) == 0 {
			return 0, err
		}

		res, err = u.notificationService.UpdateNotify(ctx, notifications[0].Id, senderTitle, senderContent, userId, notification_domain.NotificationInform, nil, false)
		if err != nil {
			return res, err
		}

		// notify to requester that requestee accepted
		userInfo, err := u.userInfoRepository.Load(ctx, userId)

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
		}, requester, receiverTitle, receiverContent, []string{userId}, nil, notification_domain.NotificationInform)

		return res, err
	})
	return res, err

}
func (u *FriendUsecase) Unfriend(ctx context.Context, userId string, friendId string) (int64, error) {

	// check is friend
	friendRelation, err := u.friendRepository.Load(ctx, friendId, userId, domain.FriendRelation)
	if err != nil || friendRelation == nil {
		return 0, err
	}
	if friendRelation.Status != domain.AcceptAction {
		return -1, nil
	}

	friendRelation.UpdatedBy = userId
	friendRelation.UpdatedAt = time.Now()
	friendRelation.Status = domain.StatusUnfriend
	friendRelationMap := convert.ToMapOmitEmpty(friendRelation)

	res, err := u.friendRepository.Patch(ctx, friendRelationMap)
	if err != nil {
		return res, err
	}

	return res, nil
}
