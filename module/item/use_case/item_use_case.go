package usecase

import (
	"context"
	"social-todo-list/common"
	"social-todo-list/module/item/model"
)

// handler => use case -> repo > storage

type ItemStorage interface {
	GetItem(ctx context.Context, cond map[string]interface{}) (*model.TodoItem, error)
	GetItems(
		ctx context.Context,
		fitter *model.Fitter,
		paging *common.Paging,
		moreKeys ...string,
	) ([]model.TodoItem, error)
	CreateItem(ctx context.Context, data *model.TodoItemCreation) error
	UpdateItem(
		ctx context.Context,
		cond map[string]interface{},
		itemData *model.TodoItemUpdation,
	) error
	DeleteItem(ctx context.Context, cond map[string]interface{}) error
}

type itemUseCase struct {
	store ItemStorage
}

// constructor
func NewItemUseCase(store ItemStorage) *itemUseCase {
	return &itemUseCase{
		store: store, // Context giữa UseCase và Storage
	}
}

func (useCase *itemUseCase) GetItems(
	ctx context.Context,
	fitter *model.Fitter,
	paging *common.Paging,
	moreKeys ...string,
) ([]model.TodoItem, error) {
	itemData, err := useCase.store.GetItems(ctx, fitter, paging)
	if err != nil {
		return nil, common.ErrCannotListEntity(model.EntityName, err)
	}
	return itemData, nil
}

func (useCase *itemUseCase) GetItemByID(ctx context.Context, id int) (*model.TodoItem, error) {

	itemData, err := useCase.store.GetItem(ctx, map[string]interface{}{"id": id})
	if err != nil {
		return nil, common.ErrCannotGetEntity(model.EntityName, err)
	}
	return itemData, nil
}

func (useCase *itemUseCase) CreateNewItem(ctx context.Context, data *model.TodoItemCreation) error {
	if err := data.Validate(); err != nil {
		return err
	}
	if err := useCase.store.CreateItem(ctx, data); err != nil {
		return common.ErrCannotCreateEntity(model.EntityName, err)
	}
	return nil
}

func (useCase *itemUseCase) UpdateItem(
	ctx context.Context,
	id int,
	data *model.TodoItemUpdation,
) error {

	itemData, err := useCase.store.GetItem(ctx, map[string]interface{}{"id": id})

	if err != nil {
		return common.ErrCannotGetEntity(model.EntityName, err)
	}

	if itemData.Status == "Deleted" {
		return model.ErrItemIsDeleted
	}

	if err := useCase.store.UpdateItem(ctx, map[string]interface{}{"id": id}, data); err != nil {
		return common.ErrCannotUpdateEntity(model.EntityName, err)
	}

	return nil
}

func (useCase *itemUseCase) DeleteItem(ctx context.Context, id int) error {

	itemData, err := useCase.store.GetItem(ctx, map[string]interface{}{"id": id})

	if err != nil {
		return common.ErrCannotGetEntity(model.EntityName, err)
	}

	if itemData.Status == "Deleted" {
		return model.ErrItemIsDeleted
	}

	if err := useCase.store.DeleteItem(ctx, map[string]interface{}{"id": id}); err != nil {
		return common.ErrCannotDeleteEntity(model.EntityName, err)
	}

	return nil
}
