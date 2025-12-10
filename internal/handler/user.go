package handler

import (
	"errors"
	"net/http"

	"github.com/Infamous003/follow-service/internal/domain"
	"github.com/Infamous003/follow-service/internal/service"
	"github.com/Infamous003/follow-service/internal/validator"
)

type UserHandler struct {
	UserService *service.UserService
}

func NewUserHandler(us *service.UserService) *UserHandler {
	return &UserHandler{
		UserService: us,
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
	}

	if err := readJSON(w, r, &input); err != nil {
		badRequestResponse(w, r, err)
		return
	}

	user, err := h.UserService.CreateUser(input.Username)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			notfoundResponse(w, r)
		case errors.Is(err, domain.ErrUsernameTaken):
			conflictResponse(w, r, "username already taken")
		default:
			if vErr, ok := err.(*validator.Validator); ok {
				failedValidationResponse(w, r, vErr.Errors)
				return
			}
			serverErrorResponse(w, r, err)
		}
		return
	}

	err = writeJSON(w, http.StatusCreated, envelope{"user": user}, nil)
	if err != nil {
		serverErrorResponse(w, r, err)
	}
}
