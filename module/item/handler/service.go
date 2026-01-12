package handler

import (
	"context"
	"social-todo-list/common"
	"social-todo-list/module/item/model"
)

type ItemUseCase interface {
	GetItemByID(ctx context.Context, id int) (*model.TodoItem, error)
	GetItems(ctx context.Context, fitter *model.Fitter, paging *common.Paging, moreKeys ...string) ([]model.TodoItem, error)
	CreateNewItem(ctx context.Context, data *model.TodoItemCreation) error
	UpdateItem(ctx context.Context, id int, itemData *model.TodoItemUpdation) error
	DeleteItem(ctx context.Context, id int) error
}

type ItemService struct {
	useCase ItemUseCase
}

func NewItemService(useCase ItemUseCase) *ItemService {
	return &ItemService{
		useCase: useCase, //Context giữa UseCase và Service
	}
}

func (s *ItemService) GetItems(ctx context.Context, fitter *model.Fitter, paging *common.Paging, moreKeys ...string) ([]model.TodoItem, error) {
	return s.useCase.GetItems(ctx, fitter, paging, moreKeys...)
}

func (s *ItemService) GetItemByID(ctx context.Context, id int) (*model.TodoItem, error) {
	return s.useCase.GetItemByID(ctx, id)
}

func (s *ItemService) CreateNewItem(ctx context.Context, data *model.TodoItemCreation) error {
	return s.useCase.CreateNewItem(ctx, data)
}

func (s *ItemService) UpdateItem(ctx context.Context, id int, itemData *model.TodoItemUpdation) error {
	return s.useCase.UpdateItem(ctx, id, itemData)
}

func (s *ItemService) DeleteItem(ctx context.Context, id int) error {
	return s.useCase.DeleteItem(ctx, id)
}
