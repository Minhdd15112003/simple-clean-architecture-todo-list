package storage

import (
	"context"
	"errors"
	"social-todo-list/common"
	"social-todo-list/module/item/model"

	"gorm.io/gorm"
)

func (s *sqlStore) GetItems(
	ctx context.Context,
	fitter *model.Fitter,
	paging *common.Paging,
	moreKeys ...string,
) ([]model.TodoItem, error) {
	var itemData []model.TodoItem

	db := s.db.Table(model.TodoItem{}.TableName()).Where("status <> ?", "Deleted")

	if f := fitter; f != nil {
		if v := f.Status; v != "" {
			db = db.Where(&model.TodoItem{
				Status: v,
			})
		}
	}

	if err := db.Table(model.TodoItem{}.TableName()).Select("id").Count(&paging.Total).Error; err != nil {
		return nil, err
	}

	if err := db.
		Select("*").
		Offset((paging.Page - 1) * paging.Limit).
		Limit(paging.Limit).
		Order("id DESC").
		Find(&itemData).Error; err != nil {
		return nil, common.ErrDB(err)
	}

	return itemData, nil
}

func (s *sqlStore) GetItem(
	ctx context.Context,
	cond map[string]interface{},
) (*model.TodoItem, error) {
	var itemData model.TodoItem

	if err := s.db.Where(cond).Where("status <> ?", "Deleted").First(&itemData).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.RecordNotFound
		}
		return nil, common.ErrDB(err)
	}
	return &itemData, nil
}

func (s *sqlStore) CreateItem(ctx context.Context, data *model.TodoItemCreation) error {
	if err := s.db.Create(&data).Error; err != nil {
		return common.ErrDB(err)
	}
	return nil
}

func (s *sqlStore) UpdateItem(
	ctx context.Context,
	cond map[string]interface{},
	itemData *model.TodoItemUpdation,
) error {
	if err := s.db.Where(cond).Updates(&itemData).Error; err != nil {
		return common.ErrDB(err)
	}
	return nil
}

func (s *sqlStore) DeleteItem(ctx context.Context, cond map[string]interface{}) error {
	const (
		Deleted = "Deleted"
	)

	if err := s.db.Table(model.TodoItem{}.TableName()).Where(cond).Updates(map[string]interface{}{
		"status": Deleted,
	}).Error; err != nil {
		return common.ErrDB(err)
	}
	return nil
}
