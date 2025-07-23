package role

import "gorm.io/gorm"

type Role struct {
	Name string `gorm:"unique"`
	gorm.Model
}
