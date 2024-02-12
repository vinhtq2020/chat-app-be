package delivery

import "go-service/internal/friend/domain"

type FriendHandler struct {
	friendService domain.FriendService
}

func NewFriendHandler(service domain.FriendService) *FriendHandler {
	return &FriendHandler{friendService: service}
}
