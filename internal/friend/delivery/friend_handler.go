package delivery

import (
	friend_domain "go-service/internal/friend/domain"
	"go-service/pkg/logger"
	"go-service/pkg/response"
	"net/http"
)

type FriendHandler struct {
	friendService friend_domain.FriendService
	logger        *logger.Logger
}

func NewFriendHandler(service friend_domain.FriendService, logger *logger.Logger) *FriendHandler {
	return &FriendHandler{friendService: service, logger: logger}
}

func (h *FriendHandler) Create(w http.ResponseWriter, r *http.Request) {
	friendId := r.PathValue("friendId")
	if len(friendId) == 0 {
		response.Response(w, http.StatusBadRequest, nil)
		return
	}
	userId, ok := r.Context().Value("userId").(string)
	if !ok {
		response.Response(w, http.StatusUnauthorized, nil)
		return
	}

	if len(friendId) > 0 {
		res, err := h.friendService.Create(r.Context(), userId, friendId, friend_domain.FriendRelation.Value())
		handleResponse(w, res, err)
	} else {
		response.Response(w, http.StatusBadRequest, nil)
	}

}

func (h *FriendHandler) Patch(w http.ResponseWriter, r *http.Request) {
	var res int64
	var err error
	friendId := r.PathValue("friendId")
	userId, ok := r.Context().Value("userId").(string)
	if !ok {
		response.Response(w, http.StatusUnauthorized, nil)
		return
	}

	action := r.PathValue("action")

	if len(friendId) == 0 && len(action) == 0 && action != "accept" && action != "reject" && action != "cancel" {
		response.Response(w, http.StatusBadRequest, nil)
		return
	}
	switch action {
	case "accept", "reject":
		res, err = h.friendService.Response(r.Context(), userId, friendId, action)
	case "cancel":
		res, err = h.friendService.Cancel(r.Context(), userId, friendId)
	case "unfriend":
		res, err = h.friendService.Unfriend(r.Context(), userId, friendId)
	default:
		response.Response(w, http.StatusBadRequest, nil)
		return
	}

	handleResponse(w, res, err)
}

func handleResponse(w http.ResponseWriter, res int64, err error) {
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
