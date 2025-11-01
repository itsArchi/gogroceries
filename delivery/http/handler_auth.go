package http

import (
	"gogroceries/domain"
	"gogroceries/internal/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authUsecase domain.AuthUsecase
}

func NewAuthHandler(router *gin.RouterGroup, uc domain.AuthUsecase) {
	handler := &AuthHandler{
		authUsecase: uc,
	}

	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/register", handler.Register)
		authRoutes.POST("/login", handler.Login)
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req domain.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, "Input is not valid", err.Error())
		return
	}

	newUser, err := h.authUsecase.Register(&req)
	if err != nil {
		if err.Error() == "email already registered" || err.Error() == "phone number already registered" {
			helper.SendError(c, http.StatusBadRequest, "Registration failed", err.Error())
		} else {
			helper.SendError(c, http.StatusInternalServerError, "Failed to register", err.Error())
		}
		return
	}

	helper.SendSuccess(c, "Registration success", newUser)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req domain.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, "Input is not valid", err.Error())
		return
	}

	token, user, err := h.authUsecase.Login(&req)
	if err != nil {
		if err.Error() == "phone number or password is incorrect" {
			helper.SendError(c, http.StatusUnauthorized, "Login failed", err.Error())
		} else {
			helper.SendError(c, http.StatusInternalServerError, "Failed to login", err.Error())
		}
		return
	}

	loginResp := domain.LoginResponse{
		Nama:         user.Nama,
		NoTelp:       user.NoTelp,
		TanggalLahir:       user.TanggalLahir,
		Email:        user.Email,
		Token:      token,
		Tentang:      user.Tentang,
		Pekerjaan:    user.Pekerjaan,
		IdProvinsi:   user.IdProvinsi,
		IdKota:       user.IdKota,
	}

	helper.SendSuccess(c, "Login success", loginResp)
}