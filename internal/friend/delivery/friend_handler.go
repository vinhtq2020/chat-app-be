package delivery

import (
	"encoding/json"
	friend_domain "go-service/internal/friend/domain"
	"go-service/pkg/logger"
	"go-service/pkg/response"
	"net/http"
	"time"
)

type FriendHandler struct {
	friendService friend_domain.FriendService
	logger        *logger.Logger
}

func NewFriendHandler(service friend_domain.FriendService, logger *logger.Logger) *FriendHandler {
	return &FriendHandler{friendService: service, logger: logger}
}

func (h *FriendHandler) Create(w http.ResponseWriter, r *http.Request) {
	var friendRequest friend_domain.FriendRequest

	userId, ok := r.Context().Value("userId").(string)
	if !ok {
		response.Response(w, http.StatusBadRequest, nil)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&friendRequest)
	if err != nil {
		h.logger.LogError(err.Error(), nil)
		response.Response(w, http.StatusBadRequest, nil)
		return
	}

	if len(friendRequest.RequesteeId) > 0 && len(userId) > 0 {
		friendRequest.CreatedAt = time.Now()
		friendRequest.CreatedBy = userId
		friendRequest.UpdatedAt = time.Now()
		friendRequest.UpdatedBy = userId
		friendRequest.RequesterId = userId

		res, err := h.friendService.SendFriendRequest(r.Context(), friendRequest)
		handleResponse(w, res, err)
	} else {
		response.Response(w, http.StatusBadRequest, nil)
	}

}

func (h *FriendHandler) UpdateRequestStatus(w http.ResponseWriter, r *http.Request) {
	requestId := r.PathValue("id")
	action := r.PathValue("action")

	if len(requestId) == 0 && len(action) == 0 && action != "accept" && action != "reject" {
		response.Response(w, http.StatusBadRequest, nil)
		return
	}

	// friendRequest["id"] = requestId
	// friendRequest["requesteeId"] = userId
	// friendRequest["action"] = action
	// friendRequest["updatedAt"] = time.Now()
	// friendRequest["updatedBy"] = userId

	res, err := h.friendService.UpdateFriendRequest(r.Context(), requestId, action)
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
