package postgres

import (
	"errors" 
	"fmt"    
	"gogroceries/domain"

	"gorm.io/gorm"
)

type postgresTrxRepository struct {
	db *gorm.DB
}

func NewPostgresTrxRepository(db *gorm.DB) domain.TrxRepository {
	return &postgresTrxRepository{db}
}

func (r *postgresTrxRepository) Create(trx *domain.Trx, details []domain.DetailTrx, logs []domain.LogProduk) (*domain.Trx, error) {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(trx).Error; err != nil {
			return fmt.Errorf("gagal simpan trx: %w", err)
		}

		if len(details) != len(logs) {
			return errors.New("internal: jumlah detail dan log tidak cocok")
		}

		for i := range details {
			details[i].IdTrx = trx.ID

			if err := tx.Create(&details[i]).Error; err != nil {
				return fmt.Errorf("gagal simpan detail trx produk ID %d: %w", details[i].IdProduk, err)
			}

			logs[i].IdDetailTrx = details[i].ID

			if err := tx.Create(&logs[i]).Error; err != nil {
				return fmt.Errorf("gagal simpan log produk ID %d: %w", details[i].IdProduk, err)
			}

			result := tx.Model(&domain.Produk{}).Where("id = ? AND stok >= ?", details[i].IdProduk, details[i].Kuantitas).
				UpdateColumn("stok", gorm.Expr("stok - ?", details[i].Kuantitas))

			if result.Error != nil {
				return fmt.Errorf("gagal update stok produk ID %d: %w", details[i].IdProduk, result.Error)
			}
			if result.RowsAffected == 0 {
				return fmt.Errorf("stok produk ID %d tidak mencukupi saat update", details[i].IdProduk)
			}
		}

		return nil 
	})

	if err != nil {
		return nil, err 
	}

	err = r.db.Preload("AlamatKirim").
		Preload("DetailTrx").
		Preload("DetailTrx.LogProduk").
		Preload("DetailTrx.Toko").
		First(&trx, trx.ID).Error

	if err != nil {
		fmt.Printf("Warning: gagal preload data trx setelah create: %v", err)
	}

	return trx, nil
}

func (r *postgresTrxRepository) FindByID(id uint, userID uint) (*domain.Trx, error) {
	var trx domain.Trx
	err := r.db.Preload("AlamatKirim").
		Preload("DetailTrx").
		Preload("DetailTrx.LogProduk").
		Preload("DetailTrx.Toko").
		Where("id = ? AND id_user = ?", id, userID). 
		First(&trx, id).Error
	return &trx, err
}

func (r *postgresTrxRepository) FindAllByUserID(userID uint, filter domain.TrxFilter, limit, offset int) ([]domain.Trx, int64, error) {
	var trxs []domain.Trx
	var total int64

	query := r.db.Model(&domain.Trx{}).Where("id_user = ?", userID)

	if filter.KodeInvoice != "" {
		query = query.Where("kode_invoice ILIKE ?", "%"+filter.KodeInvoice+"%")
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Preload("AlamatKirim").
		Preload("DetailTrx").
		Preload("DetailTrx.LogProduk").
		Preload("DetailTrx.Toko").
		Limit(limit).Offset(offset).Order("created_at DESC").Find(&trxs).Error

	return trxs, total, err
}

func (r *postgresTrxRepository) FindByIDAndUserID(id uint, userID uint) (*domain.Trx, error) {
	return r.FindByID(id, userID) 
}

func (r *postgresTrxRepository) UpdateStatus(id uint, status string) error {
	return r.db.Model(&domain.Trx{}).
		Where("id = ?", id).
		Update("status", status).
		Error
}

func (r *postgresTrxRepository) FindByTokoID(tokoID uint) ([]domain.Trx, error) {
	var trx []domain.Trx
	err := r.db.Joins("JOIN detail_trx ON trx.id = detail_trx.id_trx").
		Joins("JOIN produk ON detail_trx.id_produk = produk.id").
		Preload("User").
		Preload("AlamatKirim"). 
		Preload("DetailTrx").
		Preload("DetailTrx.LogProduk"). 
		Where("detail_trx.id_toko = ?", tokoID). 
		Distinct("trx.*").
		Find(&trx).Error
	return trx, err
}