package handler

import (
	"github.kostyadodich/demo/pkg/repository"
)

type Item struct {
	itemRepo *repository.Item
}

func NewItem(itemRepo *repository.Item) *Item {
	return &Item{
		itemRepo: itemRepo,
	}
}
