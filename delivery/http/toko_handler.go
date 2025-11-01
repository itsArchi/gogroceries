package http

import (
	"net/http"
	"strconv"

	"gogroceries/domain"
	"gogroceries/internal/helper"

	"github.com/gin-gonic/gin"
)

type TokoHandler struct {
	tokoUC  domain.TokoUsecase
	jwtAuth helper.JWTInterface
}

func NewTokoHandler(tokoUC domain.TokoUsecase, jwtAuth helper.JWTInterface) *TokoHandler {
	return &TokoHandler{
		tokoUC:  tokoUC,
		jwtAuth: jwtAuth,
	}
}

func (h *TokoHandler) GetMyToko(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	toko, err := h.tokoUC.GetMyToko(userID)
	if err != nil {
		helper.SendError(c, http.StatusNotFound, err.Error(), nil)
		return
	}

	helper.SendSuccess(c, "Success get my toko", toko)
}

func (h *TokoHandler) UpdateToko(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	id, err := strconv.ParseUint(c.Param("id_toko"), 10, 32)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid toko ID", nil)
		return
	}

	var req domain.UpdateTokoRequest
	if err := c.ShouldBind(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid request payload", nil)
		return
	}

	file, _ := c.FormFile("url_foto")
	if file != nil {
		req.UrlFoto = "path/to/uploaded/" + file.Filename
	}

	toko, err := h.tokoUC.UpdateToko(uint(id), &req, userID)
	if err != nil {
		helper.SendError(c, http.StatusNotFound, err.Error(), nil)
		return
	}

	helper.SendSuccess(c, "Toko updated successfully", toko)
}

func (h *TokoHandler) GetAllToko(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	filter := domain.TokoFilter{
		NamaToko: c.Query("nama_toko"),
	}

	tokos, pagination, err := h.tokoUC.GetAllTokos(filter, page, limit)
	if err != nil {
		helper.SendError(c, http.StatusNotFound, err.Error(), nil)
		return
	}

	response := map[string]interface{}{
		"data":       tokos,
		"pagination": pagination,
	}

	helper.SendSuccess(c, "Success get all tokos", response)
}

func (h *TokoHandler) GetTokoByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid store ID", nil)
		return
	}

	store, err := h.tokoUC.GetTokoByID(uint(id))
	if err != nil {
		helper.SendError(c, http.StatusNotFound, err.Error(), nil)
		return
	}

	helper.SendSuccess(c, "Success get store", store)
}