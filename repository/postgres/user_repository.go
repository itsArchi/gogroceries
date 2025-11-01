package postgres

import (
	"gogroceries/domain"

	"gorm.io/gorm"
)

type postgresUserRepository struct {
	db *gorm.DB
}

func NewPostgresUserRepository(db *gorm.DB) domain.UserRepository {
	return &postgresUserRepository{db}
}

func (r *postgresUserRepository) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *postgresUserRepository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("email = ?", email).First(&user).Error
	
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *postgresUserRepository) FindByNoTelp(noTelp string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("no_telp = ?", noTelp).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *postgresUserRepository) FindById(id uint) (*domain.User, error) {
	var user domain.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *postgresUserRepository) Update(user *domain.User) error {
	return r.db.Save(user).Error
}

func (r *postgresUserRepository) Delete(user *domain.User) error {
	return r.db.Delete(user).Error
}

func (r *postgresUserRepository) FindAll() ([]domain.User, error) {
	var users []domain.User
	err := r.db.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}
