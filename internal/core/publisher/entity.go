package publisher

import "gorm.io/gorm"

type Publisher struct {
	Name string
	gorm.Model
}
