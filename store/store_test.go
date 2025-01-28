package store

import (
	"context"
	"testing"
)

func TestSQLLiteCreateItem(t *testing.T) {
	ctx := context.Background()
	s, err := NewSQLiteBackedStore()
	if err != nil {
		t.Error(err)
	}
	_, err = s.CreateItem(ctx, "some_name", "some_description")
	if err != nil {
		t.Error(err)
	}
}
