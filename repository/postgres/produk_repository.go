package postgres

import (
	"errors"
	"gogroceries/domain"

	"gorm.io/gorm"
)

type postgresProdukRepository struct {
	db *gorm.DB
}

func NewPostgresProdukRepository(db *gorm.DB) domain.ProdukRepository {
	return &postgresProdukRepository{db}
}

func (r *postgresProdukRepository) Create(produk *domain.Produk, fotoUrls []string) (*domain.Produk, error) {
	tx := r.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	if err := tx.Create(produk).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	for _, url := range fotoUrls {
		foto := domain.FotoProduk{
			IdProduk: produk.ID,
			Url:      url,
		}
		if err := tx.Create(&foto).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		produk.FotoProduk = append(produk.FotoProduk, foto)
	}

	if err := tx.Preload("Category").Preload("Toko").First(produk, produk.ID).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return produk, nil
}

func (r *postgresProdukRepository) Update(produk *domain.Produk) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.Save(produk).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (r *postgresProdukRepository) Delete(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id_produk = ?", id).Delete(&domain.FotoProduk{}).Error; err != nil {
			return err
		}

		return tx.Delete(&domain.Produk{}, id).Error
	})
}
func (r *postgresProdukRepository) FindByID(id uint) (*domain.Produk, error) {
	var produk domain.Produk
	err := r.db.Preload("Toko").
		Preload("Category").
		Preload("FotoProduk").
		First(&produk, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil 
	}

	return &produk, err
}

func (r *postgresProdukRepository) FindBySlug(slug string) (*domain.Produk, error) {
	var produk domain.Produk
	err := r.db.Where("slug = ?", slug).
		Preload("Toko").
		Preload("Category").
		Preload("FotoProduk").
		First(&produk).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &produk, err
}

func (r *postgresProdukRepository) FindAll(filter domain.ProdukFilter, offset, limit int) ([]domain.Produk, int64, error) {
	var produk []domain.Produk
	var total int64

	query := r.db.Model(&domain.Produk{}).
		Preload("Toko").
		Preload("Category").
		Preload("FotoProduk")

	if filter.NamaProduk != "" {
		query = query.Where("LOWER(nama_produk) LIKE LOWER(?)", "%"+filter.NamaProduk+"%")
	}

	if filter.CategoryID > 0 {
		query = query.Where("id_category = ?", filter.CategoryID)
	}

	if filter.TokoID > 0 {
		query = query.Where("id_toko = ?", filter.TokoID)
	}

	if filter.MinHarga > 0 {
		query = query.Where("harga_konsumen >= ?", filter.MinHarga)
	}

	if filter.MaxHarga > 0 {
		query = query.Where("harga_konsumen <= ?", filter.MaxHarga)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&produk).Error

	return produk, total, err
}

func (r *postgresProdukRepository) UpdateStok(tx *gorm.DB, produkID uint, kuantitas int) error {
	return tx.Model(&domain.Produk{}).
		Where("id = ?", produkID).
		Update("stok", gorm.Expr("stok + ?", kuantitas)).
		Error
}

func (r *postgresProdukRepository) FindFotoByProdukID(produkID uint) ([]domain.FotoProduk, error) {
	var fotos []domain.FotoProduk
	err := r.db.Where("id_produk = ?", produkID).
		Find(&fotos).Error

	if err != nil {
		return nil, err
	}

	return fotos, nil
}

func (r *postgresProdukRepository) FindByTokoID(tokoID uint) ([]domain.Produk, error) {
	var produk []domain.Produk
	err := r.db.Where("id_toko = ?", tokoID).
		Preload("Toko").
		Preload("Category").
		Preload("FotoProduk").
		Order("created_at DESC").
		Find(&produk).Error

	return produk, err
}

func (r *postgresProdukRepository) FindByIDs(ids []uint) ([]domain.Produk, error) {
	var produks []domain.Produk
	
	err := r.db.Preload("Toko").
		Preload("Category").
		Where("id IN (?)", ids).
		Find(&produks).Error

	return produks, err
}