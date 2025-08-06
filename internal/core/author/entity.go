package author

import "gorm.io/gorm"

type Author struct {
	FirstName string
	LastName  string
	gorm.Model
}
