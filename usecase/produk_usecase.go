package usecase

import (
	"errors"
	"fmt"
	"gogroceries/domain" 
	"gogroceries/internal/helper" 
	"log"
	"math"
	"time"


	"gorm.io/gorm"
)

type produkUsecase struct {
	produkRepo   domain.ProdukRepository
	tokoRepo    domain.TokoRepository
	categoryRepo domain.CategoryRepository 
}

func NewProdukUsecase(pr domain.ProdukRepository, tr domain.TokoRepository, cr domain.CategoryRepository) domain.ProdukUsecase {
	return &produkUsecase{
		produkRepo:   pr,
		tokoRepo:    tr,
		categoryRepo: cr,
	}
}

func (uc *produkUsecase) CreateProduk(req *domain.CreateProdukRequest, userID uint) (*domain.Produk, error) {
	toko, err := uc.tokoRepo.FindByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("toko tidak ditemukan untuk user ini, tidak bisa menambah produk")
		}
		return nil, fmt.Errorf("gagal mencari toko: %w", err)
	}

	_, err = uc.categoryRepo.FindByID(req.IdCategory)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("category dengan ID %d tidak ditemukan", req.IdCategory)
		}
		return nil, fmt.Errorf("gagal validasi category: %w", err)
	}

	slug := helper.Slugify(req.NamaProduk)
	_, err = uc.produkRepo.FindBySlug(slug)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("gagal cek slug: %w", err)
	}
	if err == nil { 
		slug = fmt.Sprintf("%s-%d", slug, time.Now().UnixNano())
		log.Printf("Slug duplikat, generate slug baru: %s", slug)
	}


	newProduk := &domain.Produk{
		IdToko:        toko.ID,
		IdCategory:    req.IdCategory,
		NamaProduk:    req.NamaProduk,
		Slug:          slug,
		HargaReseller: req.HargaReseller,
		HargaKonsumen: req.HargaKonsumen,
		Stok:          req.Stok,
		Deskripsi:     req.Deskripsi,
	}

	createdProduk, err := uc.produkRepo.Create(newProduk, req.Photos) 
	if err != nil {
		return nil, err
	}

	return createdProduk, nil
}

func (uc *produkUsecase) GetAllProduk(filter domain.ProdukFilter, page, limit int) ([]domain.Produk, *domain.PaginationResponse, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	offset := (page - 1) * limit

	produks, totalData, err := uc.produkRepo.FindAll(filter, limit, offset)
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
		Data:      produks,
	}

	return produks, pagination, nil
}

func (uc *produkUsecase) GetProdukByID(id uint) (*domain.Produk, error) {
	produk, err := uc.produkRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("produk tidak ditemukan")
		}
		return nil, err
	}
	return produk, nil
}

func (uc *produkUsecase) UpdateProduk(id uint, req *domain.UpdateProdukRequest, userID uint) (*domain.Produk, error) {
	produk, err := uc.produkRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("produk tidak ditemukan")
		}
		return nil, err
	}

	toko, err := uc.tokoRepo.FindByID(userID)
	if err != nil || toko.ID != produk.IdToko {
		return nil, errors.New("forbidden: anda tidak bisa mengubah produk ini")
	}

	if req.NamaProduk != "" {
		produk.NamaProduk = req.NamaProduk
		newSlug := helper.Slugify(req.NamaProduk)
        if newSlug != produk.Slug {
             _, slugErr := uc.produkRepo.FindBySlug(newSlug)
             if slugErr != nil && !errors.Is(slugErr, gorm.ErrRecordNotFound){
                return nil, fmt.Errorf("gagal cek slug baru: %w", slugErr)
             }
             if slugErr == nil { 
                newSlug = fmt.Sprintf("%s-%d", newSlug, time.Now().UnixNano())
             }
             produk.Slug = newSlug
        }
	}
	if req.IdCategory != 0 {
		_, catErr := uc.categoryRepo.FindByID(req.IdCategory)
		if catErr != nil {
			return nil, fmt.Errorf("category ID %d tidak valid", req.IdCategory)
		}
		produk.IdCategory = req.IdCategory
	}
	if req.HargaReseller > 0 { 
		produk.HargaReseller = req.HargaReseller
	}
	if req.HargaKonsumen > 0 {
		produk.HargaKonsumen = req.HargaKonsumen
	}
	if req.Stok >= 0 { 
		produk.Stok = req.Stok
	}
	if req.Deskripsi != "" {
		produk.Deskripsi = req.Deskripsi
	}


	err = uc.produkRepo.Update(produk)
	if err != nil {
		return nil, err
	}
	updatedProduk, _ := uc.produkRepo.FindByID(id) 
	return updatedProduk, nil
}

func (uc *produkUsecase) DeleteProduk(id uint, userID uint) error {
	produk, err := uc.produkRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("produk tidak ditemukan")
		}
		return err
	}

	toko, err := uc.tokoRepo.FindByID(userID)
	if err != nil || toko.ID != produk.IdToko {
		return errors.New("forbidden: anda tidak bisa menghapus produk ini")
	}

	return uc.produkRepo.Delete(id)
}