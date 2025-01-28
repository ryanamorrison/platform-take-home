package store

import "context"

// GetItem does ...
func (s *DBStore) GetItem(ctx context.Context, id uint) (*Item, error) {
	var item Item

	err := s.DB.WithContext(ctx).First(&item, id).Error

	return &item, err
}

// GetItems does ...
func (s *DBStore) GetItems(ctx context.Context) ([]Item, error) {
	var items []Item
	err := s.DB.WithContext(ctx).Find(&items).Error

	return items, err
}

// CreateItem does ...
func (s *DBStore) CreateItem(ctx context.Context, name, description string) (uint, error) {
	item := Item{
		Name:        name,
		Description: description,
	}

	err := s.DB.WithContext(ctx).Create(&item).Error

	return item.ID, err
}
