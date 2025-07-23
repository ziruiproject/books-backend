package user

import (
	"gorm.io/gorm"
	"starter/internal/core/role"
)

type User struct {
	Name     string
	Role     role.Role
	Email    string `gorm:"unique"`
	Password string `gorm:"min:8"`
	RoleID   uint
	gorm.Model
}
