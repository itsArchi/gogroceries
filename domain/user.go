package domain

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Nama         string         `gorm:"size:255;not null" json:"nama"`
	KataSandi    string         `gorm:"size:255;not null" json:"-"` 
	NoTelp       string         `gorm:"size:255;unique;not null" json:"no_telp"`
	TanggalLahir string         `gorm:"size:255" json:"tanggal_lahir"` 
	JenisKelamin string         `gorm:"size:255" json:"jenis_kelamin"` 
	Tentang      string         `gorm:"type:text" json:"tentang"`     
	Pekerjaan    string         `gorm:"size:255" json:"pekerjaan"`    
	Email        string         `gorm:"size:255;unique;not null" json:"email"`
	IdProvinsi   string         `gorm:"size:255" json:"id_provinsi"` 
	IdKota       string         `gorm:"size:255" json:"id_kota"`     
	IsAdmin      bool           `gorm:"default:false" json:"is_admin"`

	Toko        *Toko           `gorm:"foreignKey:IdUser" json:"toko,omitempty"`    
	Alamat      []Alamat        `gorm:"foreignKey:IdUser" json:"alamat,omitempty"`    
	Trx         []Trx           `gorm:"foreignKey:IdUser" json:"trx,omitempty"` 

	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type UserRepository interface {
	Create(user *User) error
	Update(user *User) error
	Delete(user *User) error
	FindAll() ([]User, error)
	FindById(id uint) (*User, error)
	FindByEmail(email string) (*User, error)
	FindByNoTelp(noTelp string) (*User, error)
}

type UserUsecase interface {
	GetProfileById(id uint) (*User, error)
	UpdateProfile(id uint, req *UpdateProfileRequest) (*User, error)
	DeleteProfile(id uint) error
}

type UpdateProfileRequest struct {
    Nama         *string `json:"nama"` 
    KataSandi    *string `json:"kata_sandi"` 
    NoTelp       *string `json:"no_telp"`
    TanggalLahir *string `json:"tanggal_Lahir"`
    JenisKelamin *string `json:"jenis_kelamin"`
    Tentang      *string `json:"tentang"`
    Pekerjaan    *string `json:"pekerjaan"`
    Email        *string `json:"email"` 
    IdProvinsi   *string `json:"id_provinsi"`
    IdKota       *string `json:"id_kota"`
}