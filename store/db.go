package store

import (
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DBStore is ...
type DBStore struct {
	*gorm.DB
}

// NewSQLiteBackedStore does ...
func NewSQLiteBackedStore() (*DBStore, error) {
	db, err := gorm.Open(sqlite.Open("tables.db"))

	if err != nil {
		return nil, err
	}
	return &DBStore{db}, nil
}

// NewPostgresBackedStore does ...
func NewPostgresBackedStore(dsn string) (*DBStore, error) {
	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		return nil, err
	}
	return &DBStore{db}, nil
}

// Migrate does ...
func (s *DBStore) Migrate() error {
	return s.DB.AutoMigrate(&Item{})
}
