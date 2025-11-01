package domain

import "github.com/golang-jwt/jwt/v5"

type RegisterRequest struct {
	Nama         string `json:"nama" binding:"required"`
	KataSandi    string `json:"kata_sandi" binding:"required,min=6"`
	NoTelp       string `json:"no_telp" binding:"required"`
	TanggalLahir string `json:"tanggal_Lahir"`
	Pekerjaan    string `json:"pekerjaan"`    
	Email        string `json:"email" binding:"required,email"`
	IdProvinsi   string `json:"id_provinsi"` 
	IdKota       string `json:"id_kota"`     
	JenisKelamin string `json:"jenis_kelamin"` 
	Tentang      string `json:"tentang"`      
}


type LoginRequest struct {
	NoTelp    string `json:"no_telp" binding:"required"` 
	KataSandi string `json:"kata_sandi" binding:"required"` 
}

type LoginResponse struct {
	Nama         string      `json:"nama"`
	NoTelp       string      `json:"no_telp"`
	TanggalLahir string      `json:"tanggal_Lahir"`
	Tentang      string      `json:"tentang"`
	Pekerjaan    string      `json:"pekerjaan"`
	Email        string      `json:"email"`
	IdProvinsi   interface{} `json:"id_provinsi"` 
	IdKota       interface{} `json:"id_kota"`   
	Token        string      `json:"token"`
}

type JWTClaims struct {
	UserID  uint   `json:"user_id"`
	Email   string `json:"email"`
	IsAdmin bool   `json:"is_admin"`
	jwt.RegisteredClaims
}

type AuthUsecase interface {
	Register(req *RegisterRequest) (*User, error)
	Login(req *LoginRequest) (string, *User, error)
}