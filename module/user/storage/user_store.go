package storage

import (
	"context"
	"errors"
	"social-todo-list/common"
	"social-todo-list/module/user/model"

	"gorm.io/gorm"
)

func (s *sqlStore) CreateUser(ctx context.Context, data *model.UserCreation) error {
	db := s.db.Begin()
	// data.PrepareForInsert()

	if err := db.Table(data.TableName()).Create(data).Error; err != nil {
		return common.ErrDB(err)
	}
	if err := db.Commit().Error; err != nil {
		db.Rollback()
		return common.ErrDB(err)
	}
	return nil

}

func (s *sqlStore) FindUser(ctx context.Context, conditions map[string]interface{}, moreInfor ...string) (*model.User, error) {
	db := s.db.Table(model.User{}.TableName())

	for i := range moreInfor {
		db = db.Preload(moreInfor[i])
	}

	var user model.User
	err := db.Where(conditions).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.RecordNotFound
		}
	}
	return &user, nil
}
