package http

import (
	"gogroceries/domain"
	"gogroceries/internal/helper"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUsecase domain.UserUsecase
	jwtHelper   helper.JWTInterface
}

func NewUserHandler(uc domain.UserUsecase, jwtHelper helper.JWTInterface) *UserHandler {
	return &UserHandler{
		userUsecase: uc,
		jwtHelper:   jwtHelper,
	}
}

func (h *UserHandler) GetMyProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		helper.SendError(c, http.StatusUnauthorized, "Unauthorized: User ID not found in token", nil)
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		helper.SendError(c, http.StatusInternalServerError, "Invalid user ID format", nil)
		return
	}

	user, err := h.userUsecase.GetProfileById(userIDUint)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Internal server error", err.Error())
		return
	}

	helper.SendSuccess(c, "Success get user profile", user)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		helper.SendError(c, http.StatusUnauthorized, "Unauthorized: User ID not found in token", nil)
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		helper.SendError(c, http.StatusInternalServerError, "Invalid user ID format", nil)
		return
	}

	var req domain.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	updatedUser, err := h.userUsecase.UpdateProfile(userIDUint, &req)
	if err != nil {
		if strings.Contains(err.Error(), "sudah terdaftar") {
			helper.SendError(c, http.StatusBadRequest, "Update failed", err.Error())
		} else {
			helper.SendError(c, http.StatusInternalServerError, "Update failed", err.Error())
		}
		return
	}

	helper.SendSuccess(c, "Profile updated successfully", updatedUser)
}