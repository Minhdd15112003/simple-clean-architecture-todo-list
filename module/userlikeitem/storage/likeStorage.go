package storage

import (
	"context"
	"social-todo-list/common"
	"social-todo-list/module/userlikeitem/model"

	"gorm.io/gorm"
)

func (s *sqlStore) Create(ctx context.Context, data *model.Like) error {
	if err := s.db.Create(data).Error; err != nil {
		return common.ErrDB(err)
	}
	return nil
}

func (s *sqlStore) Find(ctx context.Context, userID, itemID int) (*model.Like, error) {
	var data model.Like
	if err := s.db.Where("user_id = ? and item_id = ?", userID, itemID).First(&data).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.RecordNotFound
		}
		return nil, common.ErrDB(err)
	}
	return &data, nil
}

func (s *sqlStore) Delete(ctx context.Context, userID, itemID int) error {
	var data model.Like
	if err := s.db.Table(data.TableName()).
		Where("user_id = ? and item_id = ?", userID, itemID).
		Delete(nil).Error; err != nil {
		return common.ErrDB(err)
	}
	return nil
}

func (s *sqlStore) ListUsers(
	ctx context.Context,
	itemID int,
	paging *common.Paging,
) ([]common.SimpleUser, error) {
	var data []model.Like
	db := s.db.Where("item_id = ?", itemID)
	if err := db.Table(model.Like{}.TableName()).Select("user_id").Count(&paging.Total).Error; err != nil {
		return nil, common.ErrDB(err)
	}

	if err := db.
		Select("*").
		Offset((paging.Page - 1) * paging.Limit).
		Limit(paging.Limit).
		Order("created_at DESC").
		Preload("User").
		Find(&data).Error; err != nil {
		return nil, common.ErrDB(err)
	}

	users := make([]common.SimpleUser, len(data))
	for i := range users {
		users[i] = *data[i].User
	}
	return users, nil
}
