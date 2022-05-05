package database

import (
	// "context"
	"github.com/go-redis/redis/v8"
)

type List struct {
	db *redis.Client
	listName string
}

func (l List) AddItem() error {
	// var context = context.Background()
	return nil
}

func (l List) RemoveItem() error {
	return nil
}

func (l List) CountItems() error {
	return nil
}

func (l List) IsInList() error {
	return nil
}

func (l List) ClearList() error {
	return nil
}

func (l List) IndexOfItem() error {
	return nil
}	