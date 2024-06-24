package delivery

import (
	"go-service/internal/friend/domain"
	"go-service/pkg/logger"
	"go-service/pkg/response"
	"net/http"
)

type FriendHandler struct {
	friendService domain.FriendService
	logger        *logger.Logger
}

func NewFriendHandler(service domain.FriendService) *FriendHandler {
	return &FriendHandler{friendService: service}
}

func (h *FriendHandler) Create(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userId").(string)
	friendId := r.PathValue("friendId")
	if len(userId) > 0 && len(friendId) > 0 {
		res, err := h.friendService.SendFriendRequest(r.Context(), userId, friendId)
		handleError(w, res, err)
	}
}

func (h *FriendHandler) Patch(w http.ResponseWriter, r *http.Request) {

}

func (h *FriendHandler) Delete(w http.ResponseWriter, r *http.Request) {

}

func handleError(w http.ResponseWriter, res int64, err error) {
	if err != nil {
		http.Error(w, "Internal Sever Error", http.StatusInternalServerError)
	} else if res > 0 {
		response.Response(w, http.StatusOK, res)
	} else if res == 0 {
		http.Error(w, "Not Found", http.StatusNotFound)
	} else if res == -1 {
		response.Response(w, http.StatusConflict, nil)
	} else if res == -2 {
		response.Response(w, http.StatusUnauthorized, nil)
	}
}
