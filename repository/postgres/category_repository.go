package postgres

import (
	"gogroceries/domain"
	"gorm.io/gorm"
)

type postgresCategoryRepository struct {
	db *gorm.DB
}

func NewPostgresCategoryRepository(db *gorm.DB) domain.CategoryRepository {
	return &postgresCategoryRepository{db}
}

func (r *postgresCategoryRepository) Create(category *domain.Category) error {
	return r.db.Create(category).Error
}

func (r *postgresCategoryRepository) Update(category *domain.Category) error {
	return r.db.Save(category).Error
}

func (r *postgresCategoryRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Category{}, id).Error
}

func (r *postgresCategoryRepository) FindByID(id uint) (*domain.Category, error) {
	var category domain.Category
	err := r.db.First(&category, id).Error
	return &category, err
}

func (r *postgresCategoryRepository) FindAll() ([]domain.Category, error) {
	var categories []domain.Category
	err := r.db.Find(&categories).Error
	return categories, err
}