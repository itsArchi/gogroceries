package http

import (
	"errors"
	"fmt"
	"gogroceries/domain"
	"gogroceries/internal/helper"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm" 
)

const produkUploadDir = "./uploads/produk" 

type ProdukHandler struct {
	produkUsecase domain.ProdukUsecase
	jwtAuth       helper.JWTInterface
}
func NewProdukHandler(uc domain.ProdukUsecase, jwtAuth helper.JWTInterface) *ProdukHandler {
	err := os.MkdirAll(produkUploadDir, os.ModePerm)
	if err != nil {
		panic(fmt.Sprintf("Gagal membuat direktori upload produk: %v", err))
	}

	return &ProdukHandler{
		produkUsecase: uc,
		jwtAuth:       jwtAuth,
	}
}

func (h *ProdukHandler) CreateProduk(c *gin.Context) {
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

	if err := c.Request.ParseMultipartForm(10 << 20); err != nil { 
		helper.SendError(c, http.StatusBadRequest, "Gagal parse form data", err.Error())
		return
	}

	var req domain.CreateProdukRequest
	if err := c.ShouldBind(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, "Input form tidak valid", err.Error())
		return
	}

	form, _ := c.MultipartForm()
	files := form.File["photos"]
	photoFilenames := []string{}

	if len(files) == 0 {
		helper.SendError(c, http.StatusBadRequest, "Minimal 1 foto produk dibutuhkan", nil)
		return
	}

	for _, file := range files {
		ext := filepath.Ext(file.Filename)
		lowerExt := strings.ToLower(ext)
		if lowerExt != ".jpg" && lowerExt != ".png" && lowerExt != ".jpeg" {
			helper.SendError(c, http.StatusBadRequest, "Format file tidak didukung: "+ext, nil)
			return
		}

		newFileName := fmt.Sprintf("%d-%s%s", time.Now().UnixNano(), uuid.NewString(), ext)
		dst := filepath.Join(produkUploadDir, newFileName)

		if err := c.SaveUploadedFile(file, dst); err != nil {
			helper.SendError(c, http.StatusInternalServerError, "Gagal menyimpan file upload", err.Error())
			return
		}
		photoFilenames = append(photoFilenames, newFileName)
	}
	req.Photos = photoFilenames

	newProduk, err := h.produkUsecase.CreateProduk(&req, userIDUint)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Gagal membuat produk", err.Error())
		return
	}

	helper.SendSuccess(c, "Produk berhasil dibuat", newProduk)
	c.Status(http.StatusCreated)
}

func (h *ProdukHandler) GetAllProduk(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 { page = 1 }
	if limit < 1 { limit = 10 }

	categoryID, _ := strconv.ParseUint(c.Query("category_id"), 10, 32)
	tokoID, _ := strconv.ParseUint(c.Query("toko_id"), 10, 32)
	minHarga, _ := strconv.Atoi(c.Query("min_harga"))
	maxHarga, _ := strconv.Atoi(c.Query("max_harga"))

	filter := domain.ProdukFilter{
		NamaProduk: c.Query("nama_produk"),
		CategoryID: uint(categoryID),
		TokoID:     uint(tokoID),
		MinHarga:   minHarga,
		MaxHarga:   maxHarga,
	}

	produks, paginationInfo, err := h.produkUsecase.GetAllProduk(filter, page, limit)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Gagal mengambil daftar produk", err.Error())
		return
	}

	if paginationInfo != nil {
		paginationInfo.Data = produks
		helper.SendSuccess(c, "Berhasil mengambil daftar produk", paginationInfo)
	} else {
		helper.SendSuccess(c, "Berhasil mengambil daftar produk", produks)
	}

}

func (h *ProdukHandler) GetProdukByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, "ID produk tidak valid: "+idStr, nil)
		return
	}

	produk, err := h.produkUsecase.GetProdukByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || strings.Contains(strings.ToLower(err.Error()), "tidak ditemukan") {
			helper.SendError(c, http.StatusNotFound, "Produk tidak ditemukan", nil)
		} else {
			helper.SendError(c, http.StatusInternalServerError, "Gagal mengambil detail produk", err.Error())
		}
		return
	}

	helper.SendSuccess(c, "Berhasil mengambil detail produk", produk)
}

func (h *ProdukHandler) UpdateProduk(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, "ID produk tidak valid: "+idStr, nil)
		return
	}

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

	var req domain.UpdateProdukRequest
	if err := c.ShouldBind(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, "Input form tidak valid", err.Error())
		return
	}

	updatedProduk, err := h.produkUsecase.UpdateProduk(uint(id), &req, userIDUint)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || strings.Contains(strings.ToLower(err.Error()), "tidak ditemukan") {
			helper.SendError(c, http.StatusNotFound, "Produk tidak ditemukan", nil)
		} else if strings.Contains(strings.ToLower(err.Error()), "forbidden") || strings.Contains(strings.ToLower(err.Error()), "tidak bisa mengubah") {
			helper.SendError(c, http.StatusForbidden, "Anda tidak punya akses untuk mengubah produk ini", nil)
		} else {
			helper.SendError(c, http.StatusInternalServerError, "Gagal update produk", err.Error())
		}
		return
	}

	helper.SendSuccess(c, "Produk berhasil diupdate", updatedProduk)
}

func (h *ProdukHandler) DeleteProduk(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, "ID produk tidak valid: "+idStr, nil)
		return
	}

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

	err = h.produkUsecase.DeleteProduk(uint(id), userIDUint)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || strings.Contains(strings.ToLower(err.Error()), "tidak ditemukan") {
			helper.SendError(c, http.StatusNotFound, "Produk tidak ditemukan", nil)
		} else if strings.Contains(strings.ToLower(err.Error()), "forbidden") || strings.Contains(strings.ToLower(err.Error()), "tidak bisa menghapus") {
			helper.SendError(c, http.StatusForbidden, "Anda tidak punya akses untuk menghapus produk ini", nil)
		} else {
			helper.SendError(c, http.StatusInternalServerError, "Gagal menghapus produk", err.Error())
		}
		return
	}

	helper.SendSuccess(c, "Produk berhasil dihapus", fmt.Sprintf("Produk dengan ID %d telah dihapus", id))
}