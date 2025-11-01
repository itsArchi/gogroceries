package usecase

import (
	"gogroceries/domain"
)

type TokoUsecase interface {
	GetMyToko(userID uint) (*domain.Toko, error)
	UpdateToko(id uint, req *UpdateTokoRequest, userID uint) (*domain.Toko, error)
	GetAllTokos(filter TokoFilter, page, limit int) ([]domain.Toko, *domain.PaginationResponse, error)
	GetTokoByID(id uint) (*domain.Toko, error)
}

type UpdateTokoRequest struct {
	NamaToko string `form:"nama_toko"`
	UrlFoto   string `form:"-"`
}
type TokoFilter struct {
	NamaToko string
}