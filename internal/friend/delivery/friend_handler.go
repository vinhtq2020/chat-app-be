package delivery

import (
	"encoding/json"
	"go-service/internal/friend/domain"
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
	var friendRequest domain.FriendRequest

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

func (h *FriendHandler) Patch(w http.ResponseWriter, r *http.Request) {
	friendRequest := map[string]interface{}{}
	userId, ok := r.Context().Value("userId").(*string)
	if !ok {
		response.Response(w, http.StatusBadRequest, nil)
	}

	err := json.NewDecoder(r.Body).Decode(&friendRequest)
	if err != nil {
		response.Response(w, http.StatusBadRequest, nil)
		return
	}

	id := r.PathValue("id")
	if len(id) == 0 {
		response.Response(w, http.StatusBadRequest, nil)
		return
	}

	if _, exits := friendRequest["action"]; !exits {
		response.Response(w, http.StatusBadRequest, nil)
		return
	}

	friendRequest["id"] = id
	friendRequest["requesterId"] = userId
	friendRequest["updatedAt"] = time.Now()
	friendRequest["updatedBy"] = userId

	res, err := h.friendService.Patch(r.Context(), friendRequest)
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
