package http

import (
	"net/http"
	"strconv"

	"gogroceries/domain"
	"gogroceries/internal/helper"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	categoryUsecase domain.CategoryUsecase
}

func NewCategoryHandler(categoryUC domain.CategoryUsecase) *CategoryHandler {
	return &CategoryHandler{
		categoryUsecase: categoryUC,
	}
}

func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req domain.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	category, err := h.categoryUsecase.CreateCategory(&req)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to create category", nil)
		return
	}

	helper.SendSuccess(c, "Category created successfully", category)
}

func (h *CategoryHandler) GetAllCategories(c *gin.Context) {
	categories, err := h.categoryUsecase.GetAllCategories()
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to fetch categories", nil)
		return
	}

	helper.SendSuccess(c, "Categories retrieved successfully", categories)
}

func (h *CategoryHandler) GetCategoryByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid category ID", nil)
		return
	}

	category, err := h.categoryUsecase.GetCategoryByID(uint(id))
    if err != nil {
        if err.Error() == "category tidak ditemukan" {
            helper.SendError(c, http.StatusNotFound, err.Error(), nil)
        } else {
            helper.SendError(c, http.StatusInternalServerError, "Internal server error", nil)
        }
        return
    }

	helper.SendSuccess(c, "Category retrieved successfully", category)
}

func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid category ID", nil)
		return
	}

	var req domain.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	category, err := h.categoryUsecase.UpdateCategory(uint(id), &req)
	if err != nil {
		if err.Error() == "category tidak ditemukan" {
			helper.SendError(c, http.StatusNotFound, err.Error(), nil)
		} else {
			helper.SendError(c, http.StatusInternalServerError, "Internal server error", nil)
		}
		return
	}

	helper.SendSuccess(c, "Category updated successfully", category)
}

func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid category ID", nil)
		return
	}

	err = h.categoryUsecase.DeleteCategory(uint(id))
	if err != nil {
		if err.Error() == "category tidak ditemukan" {
			helper.SendError(c, http.StatusNotFound, err.Error(), nil)
		} else {
			helper.SendError(c, http.StatusInternalServerError, "Internal server error", nil)
		}
		return
	}

	helper.SendSuccess(c, "Category deleted successfully", nil)
}