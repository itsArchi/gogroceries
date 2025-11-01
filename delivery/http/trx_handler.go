package http

import (
	"errors"
	"gogroceries/domain"
	"gogroceries/internal/helper"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TrxHandler struct {
	trxUC   domain.TrxUsecase
	jwtAuth helper.JWTInterface
}

func NewTrxHandler(trxUC domain.TrxUsecase, jwtAuth helper.JWTInterface) *TrxHandler {
	return &TrxHandler{
		trxUC:   trxUC,
		jwtAuth: jwtAuth,
	}
}

func (h *TrxHandler) CreateTransaksi(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		helper.SendError(c, http.StatusUnauthorized, "User ID tidak ditemukan di context", nil)
		return
	}
	userIDUint, ok := userID.(uint)
	if !ok {
		helper.SendError(c, http.StatusInternalServerError, "User ID di context bukan uint", nil)
		return
	}

	var req domain.CreateTransaksiRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, "Input JSON tidak valid", err.Error())
		return
	}

	newTrx, err := h.trxUC.CreateTransaksi(&req, userIDUint)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "tidak ditemukan") {
			helper.SendError(c, http.StatusNotFound, "Gagal membuat transaksi", err.Error())
		} else if strings.Contains(strings.ToLower(err.Error()), "stok") {
			helper.SendError(c, http.StatusBadRequest, "Gagal membuat transaksi", err.Error())
		} else if strings.Contains(strings.ToLower(err.Error()), "alamat") {
			helper.SendError(c, http.StatusBadRequest, "Gagal membuat transaksi", err.Error())
		} else {
			helper.SendError(c, http.StatusInternalServerError, "Gagal membuat transaksi", err.Error())
		}
		return
	}

	helper.SendSuccess(c, "Transaksi berhasil dibuat", newTrx)
	c.Status(http.StatusCreated)
}

func (h *TrxHandler) GetAllTransaksiUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		helper.SendError(c, http.StatusUnauthorized, "User ID tidak ditemukan di context", nil)
		return
	}
	userIDUint, ok := userID.(uint)
	if !ok {
		helper.SendError(c, http.StatusInternalServerError, "User ID di context bukan uint", nil)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 { page = 1 }
	if limit < 1 { limit = 10 }

	filter := domain.TrxFilter{
		KodeInvoice: c.Query("kode_invoice"),
		Status:      c.Query("status"),
	}

	trxs, paginationInfo, err := h.trxUC.GetAllTransaksiUser(userIDUint, filter, page, limit)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Gagal mengambil daftar transaksi", err.Error())
		return
	}

	if paginationInfo != nil {
		paginationInfo.Data = trxs
		helper.SendSuccess(c, "Berhasil mengambil daftar transaksi", paginationInfo)
	} else {
		helper.SendSuccess(c, "Berhasil mengambil daftar transaksi", trxs)
	}
}

func (h *TrxHandler) GetTransaksiByID(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		helper.SendError(c, http.StatusUnauthorized, "User ID tidak ditemukan di context", nil)
		return
	}
	userIDUint, ok := userID.(uint)
	if !ok {
		helper.SendError(c, http.StatusInternalServerError, "User ID di context bukan uint", nil)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, "ID transaksi tidak valid: "+idStr, nil)
		return
	}

	trx, err := h.trxUC.GetTransaksiByID(uint(id), userIDUint)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || strings.Contains(strings.ToLower(err.Error()), "tidak ditemukan") {
			helper.SendError(c, http.StatusNotFound, "Transaksi tidak ditemukan", nil)
		} else if strings.Contains(strings.ToLower(err.Error()), "bukan milik anda") {
			helper.SendError(c, http.StatusForbidden, "Anda tidak punya akses ke transaksi ini", nil)
		} else {
			helper.SendError(c, http.StatusInternalServerError, "Gagal mengambil detail transaksi", err.Error())
		}
		return
	}

	helper.SendSuccess(c, "Berhasil mengambil detail transaksi", trx)
}