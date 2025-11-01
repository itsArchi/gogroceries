package domain

import (
	"time"

	"gorm.io/gorm"
)

type Toko struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	IdUser    uint           `gorm:"not null;unique" json:"user_id"`
	NamaToko  string         `gorm:"size:255;not null" json:"nama_toko"`
	UrlFoto   string         `gorm:"size:255" json:"url_foto"`
	User      *User          `gorm:"foreignKey:IdUser;references:ID" json:"-"`
	Produk    []Produk       `gorm:"foreignKey:IdToko" json:"produk,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type TokoRepository interface {
	Create(toko *Toko) error
	FindByUserID(userID uint) (*Toko, error)
	FindByID(id uint) (*Toko, error)
	Update(toko *Toko) error
	Delete(toko *Toko) error
	FindAll(filter TokoFilter, offset, limit int) ([]Toko, int64, error)
}

type TokoUsecase interface { 
	GetMyToko(userID uint) (*Toko, error)
	UpdateToko(id uint, req *UpdateTokoRequest, userID uint) (*Toko, error)
	GetAllTokos(filter TokoFilter, page, limit int) ([]Toko, *PaginationResponse, error)
	GetTokoByID(id uint) (*Toko, error)
}

type UpdateTokoRequest struct {
	NamaToko string `form:"nama_toko"`
	UrlFoto  string `form:"-"` 
}

type TokoFilter struct {
	NamaToko string
}