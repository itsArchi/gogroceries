package usecase

import (
	"errors"
	"gogroceries/domain"

	"gorm.io/gorm"
)

type alamatUsecase struct {
	alamatRepo domain.AlamatRepository
}

func NewAlamatUsecase(ar domain.AlamatRepository) domain.AlamatUsecase {
	return &alamatUsecase{
		alamatRepo: ar,
	}
}

func (uc *alamatUsecase) CreateAlamat(req *domain.CreateAlamatRequest, userID uint) (*domain.Alamat, error) {

	newAlamat := &domain.Alamat{
		IdUser:       userID,
		JudulAlamat:  req.JudulAlamat,
		NamaPenerima: req.NamaPenerima,
		NoTelp:       req.NoTelp,
		DetailAlamat: req.DetailAlamat,
	}

	err := uc.alamatRepo.Create(newAlamat)
	if err != nil {
		return nil, err
	}

	return newAlamat, nil
}

func (uc *alamatUsecase) GetAllAlamatUser(userID uint, filter domain.AlamatFilter, page, limit int) ([]domain.Alamat, *domain.PaginationResponse, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	offset := (page - 1) * limit

	alamats, totalData, err := uc.alamatRepo.FindAllByUserID(userID, filter, offset, limit)
	if err != nil {
		return nil, nil, err
	}

	totalPages := (int(totalData) + limit - 1) / limit

	pagination := &domain.PaginationResponse{
		Page:      page,
		Limit:     limit,
		TotalData: int(totalData),
		TotalPage: totalPages,
	}

	return alamats, pagination, nil
}

func (uc *alamatUsecase) GetAlamatByID(id uint, userID uint) (*domain.Alamat, error) {	
	alamat, err := uc.alamatRepo.FindByIDAndUserID(id, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("alamat tidak ditemukan")
		}
		return nil, err
	}

	return alamat, nil
}

func (uc *alamatUsecase) UpdateAlamat(id uint, req *domain.UpdateAlamatRequest, userID uint) (*domain.Alamat, error) {	
	alamat, err := uc.alamatRepo.FindByIDAndUserID(id, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("alamat tidak ditemukan")
		}
		return nil, err
	}

	if req.JudulAlamat != nil {
		alamat.JudulAlamat = *req.JudulAlamat
	}
	if req.NamaPenerima != nil {
		alamat.NamaPenerima = *req.NamaPenerima
	}
	if req.NoTelp != nil {
		alamat.NoTelp = *req.NoTelp
	}
	if req.DetailAlamat != nil {
		alamat.DetailAlamat = *req.DetailAlamat
	}

	err = uc.alamatRepo.Update(alamat)
	if err != nil {
		return nil, err
	}

	return alamat, nil
}

func (uc *alamatUsecase) DeleteAlamat(id uint, userID uint) error {	
	_, err := uc.alamatRepo.FindByIDAndUserID(id, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("alamat tidak ditemukan")
		}
		return err
	}

	return uc.alamatRepo.Delete(id, userID)
}