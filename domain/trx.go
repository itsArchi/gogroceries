package domain

import (
	"time"

	"gorm.io/gorm"
)

type Trx struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	IdUser        uint           `gorm:"not null;index" json:"-"`
	IdAlamatKirim uint           `gorm:"not null" json:"-"`
	HargaTotal    int            `gorm:"not null" json:"harga_total"` 
	KodeInvoice   string         `gorm:"size:255;uniqueIndex" json:"kode_invoice"`
	MethodBayar   string         `gorm:"size:255;not null" json:"method_bayar"`
	Status        string         `gorm:"size:50;default:'pending';index" json:"status"`

	User          *User          `gorm:"foreignKey:IdUser;references:ID" json:"-"`
	AlamatKirim   *Alamat        `gorm:"foreignKey:IdAlamatKirim;references:ID" json:"alamat_kirim"`
	DetailTrx     []DetailTrx    `gorm:"foreignKey:IdTrx;constraint:OnDelete:CASCADE;" json:"detail_trx,omitempty"`

	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

type TrxRepository interface {
	Create(trx *Trx, details []DetailTrx, logs []LogProduk) (*Trx, error)
	FindByID(id, userID uint) (*Trx, error)
	FindAllByUserID(userID uint, filter TrxFilter, limit, offset int) ([]Trx, int64, error) 
	FindByIDAndUserID(id uint, userID uint) (*Trx, error) 
}

type TrxUsecase interface {
	CreateTransaksi(req *CreateTransaksiRequest, userID uint) (*Trx, error)
	GetAllTransaksiUser(userID uint, filter TrxFilter, page, limit int) ([]Trx, *PaginationResponse, error) 
	GetTransaksiByID(id uint, userID uint) (*Trx, error)
}

type CreateTransaksiRequest struct {
	MethodBayar   string                    `json:"method_bayar" binding:"required"`
	IdAlamatKirim uint                      `json:"alamat_kirim" binding:"required"`
	DetailTrx     []CreateDetailTrxRequest `json:"detail_trx" binding:"required,min=1,dive"`
}

type TrxFilter struct {
    KodeInvoice string
    Status      string
}