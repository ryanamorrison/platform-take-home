package store

import "gorm.io/gorm"

// Item is ...
type Item struct {
	gorm.Model

	Name        string
	Description string
}
