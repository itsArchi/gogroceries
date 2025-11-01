package domain

import (
	"time"

	"gorm.io/gorm"
)

type Produk struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	IdToko        uint           `gorm:"not null;index" json:"-"` 
	IdCategory    uint           `gorm:"not null;index" json:"-"` 
	NamaProduk    string         `gorm:"size:255;not null" json:"nama_produk"`
	Slug          string         `gorm:"size:255;uniqueIndex" json:"slug"`
	HargaReseller int            `json:"harga_reseller"`
	HargaKonsumen int            `json:"harga_konsumen"` 
	Stok          int            `gorm:"not null;default:0" json:"stok"`
	Deskripsi     string         `gorm:"type:text" json:"deskripsi"`

	Toko          *Toko          `gorm:"foreignKey:IdToko;references:ID" json:"toko"`        
	Category      *Category      `gorm:"foreignKey:IdCategory;references:ID" json:"category"` 
	FotoProduk    []FotoProduk   `gorm:"foreignKey:IdProduk;constraint:OnDelete:CASCADE;" json:"photos,omitempty"` 

	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

type ProdukRepository interface {
	Create(produk *Produk, fotoUrls []string) (*Produk, error)
	FindAll(filter ProdukFilter, limit, offset int) ([]Produk, int64, error)
	FindByID(id uint) (*Produk, error)
	FindBySlug(slug string) (*Produk, error)
	FindByIDs(ids []uint) ([]Produk, error)
	Update(produk *Produk) error
	Delete(id uint) error
	FindFotoByProdukID(produkID uint) ([]FotoProduk, error)
	UpdateStok(tx *gorm.DB, produkID uint, kuantitas int) error 
}

type ProdukUsecase interface {
	CreateProduk(req *CreateProdukRequest, userID uint) (*Produk, error)
	GetAllProduk(filter ProdukFilter, page, limit int) ([]Produk, *PaginationResponse, error)
	GetProdukByID(id uint) (*Produk, error)
	UpdateProduk(id uint, req *UpdateProdukRequest, userID uint) (*Produk, error)
	DeleteProduk(id uint, userID uint) error
}

type LogProduk struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	IdDetailTrx   uint           `gorm:"not null;uniqueIndex" json:"-"`
	NamaProduk    string         `gorm:"size:255" json:"nama_produk"`
	Slug          string         `gorm:"size:255" json:"slug"`
	HargaReseller int            `json:"harga_reseller"`
	HargaKonsumen int            `json:"harga_konsumen"`
	Deskripsi     string         `gorm:"type:text" json:"deskripsi"`
	IdCategory    uint           `json:"-"`
	NamaCategory  string         `gorm:"size:255" json:"nama_category"`
	IdToko        uint           `json:"-"`
	NamaToko      string         `gorm:"size:255" json:"nama_toko"`
	UrlFotoToko   string         `gorm:"size:255" json:"url_foto_toko"` 
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

type FotoProduk struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	IdProduk  uint           `gorm:"not null;index" json:"product_id"`
	Url       string         `gorm:"size:255;not null" json:"url"`
	Produk    *Produk        `gorm:"foreignKey:IdProduk;references:ID" json:"-"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type CreateProdukRequest struct {
	NamaProduk    string   `form:"nama_produk" binding:"required"`
	IdCategory    uint     `form:"category_id" binding:"required"`
	HargaReseller int      `form:"harga_reseller" binding:"required"`
	HargaKonsumen int      `form:"harga_konsumen" binding:"required"`
	Stok          int      `form:"stok" binding:"required"`
	Deskripsi     string   `form:"deskripsi"`
	Photos        []string `form:"-"` 
}

type UpdateProdukRequest struct {
	NamaProduk    string `form:"nama_produk"`
	IdCategory    uint   `form:"category_id"`
	HargaReseller int    `form:"harga_reseller"`
	HargaKonsumen int    `form:"harga_konsumen"`
	Stok          int    `form:"stok"`
	Deskripsi     string `form:"deskripsi"`
}

type ProdukFilter struct {
	NamaProduk   string
	CategoryID   uint
	TokoID       uint
	MinHarga     int
	MaxHarga     int
}