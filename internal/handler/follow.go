package handler

import (
	"errors"
	"net/http"

	"github.com/Infamous003/follow-service/internal/domain"
	"github.com/Infamous003/follow-service/internal/service"
)

type FollowHandler struct {
	followService *service.FollowService
}

func NewFollowHandler(followService *service.FollowService) *FollowHandler {
	return &FollowHandler{followService: followService}
}

func (h *FollowHandler) FollowUser(w http.ResponseWriter, r *http.Request) {
	var input struct {
		FollowerID int64 `json:"follower_id"`
		FolloweeID int64 `json:"followee_id"`
	}

	if err := readJSON(w, r, &input); err != nil {
		badRequestResponse(w, r, err)
		return
	}

	err := h.followService.FollowUser(input.FollowerID, input.FolloweeID)
	if err != nil {
		switch err {
		case domain.ErrCannotFollowSelf:
			badRequestResponse(w, r, err)
		case domain.ErrUserNotFound:
			notfoundResponse(w, r)
		case domain.ErrAlreadyFollowing:
			conflictResponse(w, r, "already following this user")
		default:
			serverErrorResponse(w, r, err)
		}
		return
	}

	err = writeJSON(w, http.StatusOK, envelope{"message": "followed successfully"}, nil)
	if err != nil {
		serverErrorResponse(w, r, err)
	}
}

func (h *FollowHandler) UnfollowUser(w http.ResponseWriter, r *http.Request) {
	var input struct {
		FollowerID int64 `json:"follower_id"`
		FolloweeID int64 `json:"followee_id"`
	}

	if err := readJSON(w, r, &input); err != nil {
		badRequestResponse(w, r, err)
		return
	}

	err := h.followService.UnfollowUser(input.FollowerID, input.FolloweeID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			notfoundResponse(w, r)
		default:
			serverErrorResponse(w, r, err)
		}
		return
	}

	err = writeJSON(w, http.StatusOK, envelope{"message": "unfollowed successfully"}, nil)
	if err != nil {
		serverErrorResponse(w, r, err)
	}
}

func (h *FollowHandler) ListFollowers(w http.ResponseWriter, r *http.Request) {
	id, err := readIDParam(r)
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	followers, err := h.followService.ListFollowers(id)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			notfoundResponse(w, r)
		default:
			serverErrorResponse(w, r, err)
		}
		return
	}

	err = writeJSON(w, http.StatusOK, envelope{"followers": followers}, nil)
	if err != nil {
		serverErrorResponse(w, r, err)
	}
}

func (h *FollowHandler) ListFollowing(w http.ResponseWriter, r *http.Request) {
	id, err := readIDParam(r)
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	following, err := h.followService.ListFollowing(id)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			notfoundResponse(w, r)
		default:
			serverErrorResponse(w, r, err)
		}
		return
	}

	err = writeJSON(w, http.StatusOK, envelope{"following": following}, nil)
	if err != nil {
		serverErrorResponse(w, r, err)
	}
}
