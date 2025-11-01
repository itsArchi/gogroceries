package domain

import (
	"time"

	"gorm.io/gorm"
)

type DetailTrx struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	IdTrx       uint           `gorm:"not null;index" json:"-"`
	IdProduk    uint           `gorm:"not null" json:"-"` 
	IdToko      uint           `gorm:"not null" json:"-"`
	Kuantitas   int            `gorm:"not null" json:"kuantitas"`   
	HargaTotal  int            `gorm:"not null" json:"harga_total"` 

	Trx         *Trx           `gorm:"foreignKey:IdTrx;references:ID" json:"-"`
	LogProduk   *LogProduk     `gorm:"foreignKey:IdDetailTrx;references:ID;constraint:OnDelete:CASCADE;" json:"product"` 
	Toko        *Toko          `gorm:"foreignKey:IdToko;references:ID" json:"toko"`                 

	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type CreateDetailTrxRequest struct {
	IdProduk  uint `json:"product_id" binding:"required"`
	Kuantitas int  `json:"kuantitas" binding:"required,gt=0"`
}