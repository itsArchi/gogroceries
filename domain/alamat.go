package domain

import (
	"time"

	"gorm.io/gorm"
)

type Alamat struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	IdUser    uint           `gorm:"not null" json:"-"`
	JudulAlamat  string         `gorm:"size:255" json:"judul_alamat"`
	NamaPenerima string         `gorm:"size:255" json:"nama_penerima"`
	NoTelp       string         `gorm:"size:255" json:"no_telp"`
	DetailAlamat string         `gorm:"type:text" json:"detail_alamat"`
	User      *User          `gorm:"foreignKey:IdUser;references:ID" json:"-"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type AlamatRepository interface { 
	FindByIDAndUserID(id uint, userID uint) (*Alamat, error)
	FindAllByUserID(userID uint, filter AlamatFilter, limit, offset int) ([]Alamat, int64, error) 
	Create(alamat *Alamat) error 
	Update(alamat *Alamat) error 
	Delete(id uint, userID uint) error
	FindByID(id uint) (*Alamat, error)
}

type AlamatUsecase interface { 
	CreateAlamat(req *CreateAlamatRequest, userID uint) (*Alamat, error)
	GetAllAlamatUser(userID uint, filter AlamatFilter, page, limit int) ([]Alamat, *PaginationResponse, error)
	GetAlamatByID(id uint, userID uint) (*Alamat, error)
	UpdateAlamat(id uint, req *UpdateAlamatRequest, userID uint) (*Alamat, error)
	DeleteAlamat(id uint, userID uint) error
}

type CreateAlamatRequest struct {
	JudulAlamat  string `json:"judul_alamat" binding:"required"`
	NamaPenerima string `json:"nama_penerima" binding:"required"`
	NoTelp       string `json:"no_telp" binding:"required"`
	DetailAlamat string `json:"detail_alamat" binding:"required"`
}
type UpdateAlamatRequest struct {
	JudulAlamat  *string `json:"judul_alamat"`
	NamaPenerima *string `json:"nama_penerima"`
	NoTelp       *string `json:"no_telp"`
	DetailAlamat *string `json:"detail_alamat"`
}
type AlamatFilter struct { 
	JudulAlamat string
}