package usecase

import (
	"errors"
	"gogroceries/domain" 

	"gorm.io/gorm"
)

type categoryUsecase struct {
	categoryRepo domain.CategoryRepository
}

func NewCategoryUsecase(cr domain.CategoryRepository) domain.CategoryUsecase {
	return &categoryUsecase{
		categoryRepo: cr,
	}
}

func (uc *categoryUsecase) CreateCategory(req *domain.CreateCategoryRequest) (*domain.Category, error) {
	newCategory := &domain.Category{
		NamaCategory: req.NamaCategory,
	}
	err := uc.categoryRepo.Create(newCategory)
	if err != nil {
		return nil, err
	}
	return newCategory, nil
}

func (uc *categoryUsecase) GetAllCategories() ([]domain.Category, error) {
	return uc.categoryRepo.FindAll()
}

func (uc *categoryUsecase) GetCategoryByID(id uint) (*domain.Category, error) {
	category, err := uc.categoryRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category tidak ditemukan")
		}
		return nil, err
	}
	return category, nil
}

func (uc *categoryUsecase) UpdateCategory(id uint, req *domain.UpdateCategoryRequest) (*domain.Category, error) {
	category, err := uc.categoryRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category tidak ditemukan")
		}
		return nil, err
	}

	category.NamaCategory = req.NamaCategory
	err = uc.categoryRepo.Update(category)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (uc *categoryUsecase) DeleteCategory(id uint) error {
	_, err := uc.categoryRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("category tidak ditemukan")
		}
		return err
	}
	return uc.categoryRepo.Delete(id)
}