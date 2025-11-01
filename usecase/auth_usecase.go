package usecase

import (
	"errors"
	"fmt"
	"gogroceries/domain"
	"gogroceries/internal/helper"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type authUsecase struct {
	userRepo domain.UserRepository
	tokoRepo domain.TokoRepository
	jwtAuth  helper.JWTInterface
}

func NewAuthUsecase(ur domain.UserRepository, tr domain.TokoRepository, jwtAuth helper.JWTInterface) domain.AuthUsecase {
	return &authUsecase{
		userRepo: ur,
		tokoRepo: tr,
		jwtAuth:  jwtAuth,
	}
}	

func (uc *authUsecase) Register(req *domain.RegisterRequest) (*domain.User, error) {
	existingUser, err := uc.userRepo.FindByEmail(req.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("failed to check email")
	}
	if existingUser != nil {
		return nil, errors.New("email already registered")
	}

	existingUser, err = uc.userRepo.FindByNoTelp(req.NoTelp)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("failed to check phone number")
	}
	if existingUser != nil {
		return nil, errors.New("phone number already registered")
	}

	hashedPassword, err := helper.HashPassword(req.KataSandi)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	newUser := &domain.User{
		Nama:         req.Nama,
		KataSandi:    hashedPassword,
		NoTelp:       req.NoTelp,
		TanggalLahir: req.TanggalLahir,
		Pekerjaan:    req.Pekerjaan,    
		Email:        req.Email,
		IdProvinsi:   req.IdProvinsi,   
		IdKota:       req.IdKota,       
		JenisKelamin: req.JenisKelamin, 
		Tentang:      req.Tentang,      
		IsAdmin:      false,
	}

	err = uc.userRepo.Create(newUser)
	if err != nil {
		return nil, errors.New("failed to create user")
	}

	newToko := &domain.Toko{
		IdUser:   newUser.ID,
		NamaToko: fmt.Sprintf("%s Toko", newUser.Nama),
		UrlFoto:  "",
	}
	err = uc.tokoRepo.Create(newToko)
	if err != nil {
		log.Printf("Warning: failed to create toko for user %d: %v", newUser.ID, err)
	}

	newUser.KataSandi = ""
	return newUser, nil
}

func (uc *authUsecase) Login(req *domain.LoginRequest) (string, *domain.User, error) {
	user, err := uc.userRepo.FindByNoTelp(req.NoTelp)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil, errors.New("phone number or password is incorrect")
		}
		return "", nil, errors.New("failed to find user")
	}

	isValid := helper.CheckPasswordHash(req.KataSandi, user.KataSandi)
	if !isValid {
		return "", nil, errors.New("phone number or password is incorrect")
	}

	claims := &domain.JWTClaims{
		UserID:  user.ID,
		Email:   user.Email,
		IsAdmin: user.IsAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)), 
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token, err := uc.jwtAuth.GenerateToken(claims)
	if err != nil {
		return "", nil, errors.New("failed to generate token")
	}

	user.KataSandi = ""
	return token, user, nil
}