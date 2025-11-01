package domain

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	NamaCategory string         `gorm:"size:255;not null;uniqueIndex" json:"nama_category"` 
	Produk       []Produk       `gorm:"foreignKey:IdCategory" json:"-"` 
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

type CategoryRepository interface {
	FindByID(id uint) (*Category, error)
	FindAll() ([]Category, error)
	Create(category *Category) error
	Update(category *Category) error
	Delete(id uint) error
}

type CategoryUsecase interface {
	CreateCategory(req *CreateCategoryRequest) (*Category, error)
	GetAllCategories() ([]Category, error)
	GetCategoryByID(id uint) (*Category, error)
	UpdateCategory(id uint, req *UpdateCategoryRequest) (*Category, error)
	DeleteCategory(id uint) error
}

type CreateCategoryRequest struct {
	NamaCategory string `json:"nama_category" binding:"required"` 
}

type UpdateCategoryRequest struct {
	NamaCategory string `json:"nama_category" binding:"required"` 
}