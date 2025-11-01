package http

import (
	"net/http"
	"strconv"

	"gogroceries/domain"
	"gogroceries/internal/helper"

	"github.com/gin-gonic/gin"
)

type AlamatHandler struct {
	alamatUC domain.AlamatUsecase
	jwtAuth  helper.JWTInterface
}

func NewAlamatHandler(alamatUC domain.AlamatUsecase, jwtAuth helper.JWTInterface) *AlamatHandler {
	return &AlamatHandler{
		alamatUC: alamatUC,
		jwtAuth:  jwtAuth,
	}
}

func (h *AlamatHandler) CreateAlamat(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	var req domain.CreateAlamatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid request payload", nil)
		return
	}

	alamat, err := h.alamatUC.CreateAlamat(&req, userID)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	helper.SendSuccess(c, "Alamat created successfully", alamat)
}

func (h *AlamatHandler) GetAllAlamatUser(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))	
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	filter := domain.AlamatFilter{
		JudulAlamat: c.Query("judul_alamat"),
	}

	alamats, pagination, err := h.alamatUC.GetAllAlamatUser(userID, filter, page, limit)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response := map[string]interface{}{
		"data":       alamats,
		"pagination": pagination,
	}

	helper.SendSuccess(c, "Success get all alamat", response)
}

func (h *AlamatHandler) GetAlamatByID(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid alamat ID", nil)
		return
	}

	alamat, err := h.alamatUC.GetAlamatByID(uint(id), userID)
	if err != nil {
		if err.Error() == "alamat tidak ditemukan" {
			helper.SendError(c, http.StatusNotFound, err.Error(), nil)
			return
		}

		helper.SendError(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	helper.SendSuccess(c, "Success get alamat", alamat)
}

func (h *AlamatHandler) UpdateAlamat(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid alamat ID", nil)
		return
	}

	var req domain.UpdateAlamatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid request payload", nil)
		return
	}

	alamat, err := h.alamatUC.UpdateAlamat(uint(id), &req, userID)
	if err != nil {
		if err.Error() == "alamat tidak ditemukan" {
			helper.SendError(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		
		helper.SendError(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	helper.SendSuccess(c, "Alamat updated successfully", alamat)
}

func (h *AlamatHandler) DeleteAlamat(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid alamat ID", nil)
		return
	}

	if err := h.alamatUC.DeleteAlamat(uint(id), userID); err != nil {
		if err.Error() == "alamat tidak ditemukan" {
			helper.SendError(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		helper.SendError(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	helper.SendSuccess(c, "Alamat deleted successfully", nil)
}