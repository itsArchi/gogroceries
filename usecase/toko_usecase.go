package usecase

import (
	"errors"
	"gogroceries/domain"
	"math"

	"gorm.io/gorm"
)

type tokoUsecase struct {
	tokoRepo domain.TokoRepository
}

func NewTokoUsecase(sr domain.TokoRepository) domain.TokoUsecase {
	return &tokoUsecase{
		tokoRepo: sr,
	}
}

func (uc *tokoUsecase) GetMyToko(userID uint) (*domain.Toko, error) {
	toko, err := uc.tokoRepo.FindByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("toko tidak ditemukan untuk user ini")
		}
		return nil, err
	}
	return toko, nil
}

func (uc *tokoUsecase) UpdateToko(id uint, req *domain.UpdateTokoRequest, userID uint) (*domain.Toko, error) {
	toko, err := uc.tokoRepo.FindByUserID(userID) 
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("toko tidak ditemukan")
		}
		return nil, err
	}

	if toko.ID != id {
		return nil, errors.New("forbidden: anda tidak bisa mengubah toko ini")
	}

	if req.NamaToko != "" {
		toko.NamaToko = req.NamaToko
	}
	if req.UrlFoto != "" { 
		toko.UrlFoto = req.UrlFoto
	}

	err = uc.tokoRepo.Update(toko)
	if err != nil {
		return nil, err
	}
	return toko, nil
}

func (uc *tokoUsecase) GetAllTokos(filter domain.TokoFilter, page, limit int) ([]domain.Toko, *domain.PaginationResponse, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10 
	}
	offset := (page - 1) * limit
	
	tokos, totalData, err := uc.tokoRepo.FindAll(filter, offset, limit)
	if err != nil {
		return nil, nil, err
	}

	totalPages := 0
	if totalData > 0 {
		totalPages = int(math.Ceil(float64(totalData) / float64(limit)))
	}

	pagination := &domain.PaginationResponse{
		Page:      page,
		Limit:     limit,
		TotalData: int(totalData),
		TotalPage: totalPages,
	}

	return tokos, pagination, err
}

func (uc *tokoUsecase) GetTokoByID(id uint) (*domain.Toko, error) {
    toko, err := uc.tokoRepo.FindByID(id) 
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("toko tidak ditemukan")
			}
			return nil, err
		}
		return toko, nil
}