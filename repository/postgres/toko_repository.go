package postgres

import (
	"gogroceries/domain"

	"gorm.io/gorm"
)

type postgresTokoRepository struct {
	db *gorm.DB
}

func NewPostgresTokoRepository(db *gorm.DB) domain.TokoRepository {
	return &postgresTokoRepository{db}
}

func (r *postgresTokoRepository) Create(toko *domain.Toko) error {
	return r.db.Create(toko).Error
}

func (r *postgresTokoRepository) FindByUserID(userID uint) (*domain.Toko, error) {
	var toko domain.Toko
	err := r.db.Where("id_user = ?", userID).First(&toko).Error
	if err != nil {
		return nil, err
	}
	return &toko, nil
}

func (r *postgresTokoRepository) FindByID(id uint) (*domain.Toko, error) {
    var toko domain.Toko
    err := r.db.First(&toko, id).Error
    if err != nil {
        return nil, err
    }
    return &toko, nil
}

func (r *postgresTokoRepository) Update(toko *domain.Toko) error {
	return r.db.Save(toko).Error
}

func (r *postgresTokoRepository) Delete(toko *domain.Toko) error {
	return r.db.Delete(toko).Error
}

func (r *postgresTokoRepository) FindAll(filter domain.TokoFilter, offset, limit int) ([]domain.Toko, int64, error) {
	var tokos []domain.Toko
	var total int64

	query := r.db.Model(&domain.Toko{})

	if filter.NamaToko != "" {
		query = query.Where("nama_toko ILIKE ?", "%"+filter.NamaToko+"%")
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Offset(offset).Limit(limit).Find(&tokos).Error
	if err != nil {
		return nil, 0, err
	}

	return tokos, total, nil
}

func (r *postgresTokoRepository) FindByIDToko(id uint) (*domain.Toko, error) {
    return r.FindByID(id)
}