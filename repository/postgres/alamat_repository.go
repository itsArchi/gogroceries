package postgres

import (
	"gogroceries/domain"

	"gorm.io/gorm"
)

type postgresAlamatRepository struct {
	db *gorm.DB
}

func NewPostgresAlamatRepository(db *gorm.DB) domain.AlamatRepository {
	return &postgresAlamatRepository{db}
}

func (r *postgresAlamatRepository) Create(alamat *domain.Alamat) error {
	return r.db.Create(alamat).Error
}

func (r *postgresAlamatRepository) Update(alamat *domain.Alamat) error {
	return r.db.Save(alamat).Error
}

func (r *postgresAlamatRepository) Delete(id, userID uint) error {
	return r.db.Where("id = ? AND id_user = ?", id, userID).Delete(&domain.Alamat{}).Error
}

func (r *postgresAlamatRepository) FindByID(id uint) (*domain.Alamat, error) {
	var alamat domain.Alamat
	err := r.db.First(&alamat, id).Error
	return &alamat, err
}

func (r *postgresAlamatRepository) FindAllByUserID(userID uint, filter domain.AlamatFilter, offset, limit int) ([]domain.Alamat, int64, error) {
	var alamats []domain.Alamat
	var total int64

	query := r.db.Model(&domain.Alamat{}).Where("id_user = ?", userID)

	if filter.JudulAlamat != "" {
		query = query.Where("judul_alamat ILIKE ?", "%"+filter.JudulAlamat+"%")
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&alamats).Error
	if err != nil {
		return nil, 0, err
	}

	return alamats, total, nil
}

func (r *postgresAlamatRepository) FindByIDAndUserID(id uint, userID uint) (*domain.Alamat, error) {
	var alamat domain.Alamat
	err := r.db.Where("id = ? AND id_user = ?", id, userID).First(&alamat).Error
	return &alamat, err
}