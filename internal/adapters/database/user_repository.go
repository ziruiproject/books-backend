package database

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"starter/internal/core/pagination"
	"starter/internal/core/role"
	"starter/internal/core/user"
)

type UserRepository struct {
}

func NewUserRepository() user.Repository {
	return &UserRepository{}
}

func (repository *UserRepository) Save(db *gorm.DB, user *user.User) error {
	var role role.Role
	result := db.FirstOrCreate(&role, user.Role)
	if result.Error != nil {
		log.Error().
			Err(result.Error).
			Msgf("Failed to save role")
		return result.Error
	}
	user.Role = role

	result = db.Create(user)
	if result.Error != nil {
		log.Error().
			Err(result.Error).
			Msgf("Failed to save user")
		return result.Error
	}

	return nil
}

func (repository *UserRepository) Update(db *gorm.DB, user *user.User) error {
	result := db.Save(user)
	if result.Error != nil {
		log.Error().
			Err(result.Error).
			Msgf("Failed to save user")
		return result.Error
	}

	return nil
}

func (repository *UserRepository) Delete(db *gorm.DB, id int) error {
	result := db.Delete(&user.User{}, id)
	if result.Error != nil {
		log.Error().
			Err(result.Error).
			Msgf("Failed to delete user")
		return result.Error
	}

	return nil
}

func (repository *UserRepository) FindAll(db *gorm.DB, params pagination.Request) ([]user.User, int64, error) {
	var users []user.User
	var count int64

	db.Model(&users).Count(&count)

	result := db.
		Debug().
		Order(fmt.Sprintf("%s %s", params.OrderBy, params.SortBy)).
		Limit(params.Limit).
		Offset((params.Page - 1) * params.Limit).
		Find(&users)

	if result.Error != nil {
		log.Error().
			Err(result.Error).
			Msgf("Failed to find all users")
	}

	return users, count, result.Error
}

func (repository *UserRepository) FindByID(db *gorm.DB, id int) (user.User, error) {
	var user user.User
	result := db.Preload("Role").First(&user, id)
	if result.Error != nil {
		log.Error().
			Err(result.Error).
			Msgf("Failed to find user")
		return user, result.Error
	}

	return user, nil
}

func (repository *UserRepository) FindByEmail(db *gorm.DB, email string) (user.User, error) {
	var user user.User
	result := db.Where("email = ?", email).Preload("Role").First(&user)
	if result.Error != nil {
		log.Error().
			Err(result.Error).
			Msgf("Failed to find user")
		return user, result.Error
	}
	return user, nil
}
