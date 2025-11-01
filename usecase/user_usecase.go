package usecase

import (
	"errors"
	"fmt"
	"gogroceries/domain"
	"gogroceries/internal/helper"
	"gorm.io/gorm"
)

type userUsecase struct {
	userRepo domain.UserRepository
}

func NewUserUsecase(userRepo domain.UserRepository) domain.UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
	}
}

func (uc *userUsecase) GetProfileById(id uint) (*domain.User, error) {
	user, err := uc.userRepo.FindById(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (uc *userUsecase) UpdateProfile(id uint, req *domain.UpdateProfileRequest) (*domain.User, error) {
	existingUser, err := uc.userRepo.FindById(id)
	if err != nil {
		return nil, err
	}

	if req.Nama != nil {
		existingUser.Nama = *req.Nama
	}
	if req.TanggalLahir != nil {
		existingUser.TanggalLahir = *req.TanggalLahir
	}
	if req.JenisKelamin != nil {
		existingUser.JenisKelamin = *req.JenisKelamin
	}
	if req.Tentang != nil {
		existingUser.Tentang = *req.Tentang
	}
	if req.Pekerjaan != nil {
		existingUser.Pekerjaan = *req.Pekerjaan
	}
	if req.IdProvinsi != nil {
		existingUser.IdProvinsi = *req.IdProvinsi
	}
	if req.IdKota != nil {
		existingUser.IdKota = *req.IdKota
	}

	if req.Email != nil && *req.Email != existingUser.Email {
		_, err := uc.userRepo.FindByEmail(*req.Email)
		if err == nil {
			return nil, errors.New("email sudah terdaftar")
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("gagal cek email: %w", err)
		}
		existingUser.Email = *req.Email
	}

	if req.NoTelp != nil && *req.NoTelp != existingUser.NoTelp {
		_, err := uc.userRepo.FindByNoTelp(*req.NoTelp)
		if err == nil {
			return nil, errors.New("nomor telepon sudah terdaftar")
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("gagal cek no telp: %w", err)
		}
		existingUser.NoTelp = *req.NoTelp
	}

	if req.KataSandi != nil && *req.KataSandi != "" {
		hashedPassword, err := helper.HashPassword(*req.KataSandi)
		if err != nil {
			return nil, errors.New("gagal hash password baru")
		}
		existingUser.KataSandi = hashedPassword
	}

	if err := uc.userRepo.Update(existingUser); err != nil {
		return nil, err
	}

	return existingUser, nil
}

func (uc *userUsecase) DeleteProfile(id uint) error {
	user, err := uc.userRepo.FindById(id)
	if err != nil {
		return err
	}
	return uc.userRepo.Delete(user)
}