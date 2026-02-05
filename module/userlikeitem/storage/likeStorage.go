package storage

import (
	"context"
	"social-todo-list/common"
	"social-todo-list/module/userlikeitem/model"
	"time"

	"github.com/btcsuite/btcutil/base58"
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

	// Sửa time layout - dùng format chuẩn
	timeLayout := "2006-01-02T15:04:05.999999"
	if v := paging.FakeCursor; v != "" {
		timeCreated, err := time.Parse(timeLayout, string(base58.Decode(v)))
		if err != nil {
			return nil, common.ErrDB(err)
		}

		db = db.Where("created_at < ?", timeCreated.Format("2006-01-02 16:04:05.999999"))
	} else {
		db = db.Offset((paging.Page - 1) * paging.Limit)
	}

	if err := db.
		Select("*").
		Limit(paging.Limit).
		Order("created_at DESC").
		Preload("User").
		Find(&data).Error; err != nil {
		return nil, common.ErrDB(err)
	}

	users := make([]common.SimpleUser, len(data))
	for i := range users {
		users[i] = *data[i].User
		users[i].UpdatedAt = nil
		users[i].CreatedAt = data[i].CreatedAt
	}
	if len(users) > 0 {
		users[len(users)-1].Mask()
		paging.NextCursor = base58.Encode([]byte(users[len(data)-1].CreatedAt.Format(timeLayout)))
	}
	return users, nil
}

func (s *sqlStore) GetLikeItem(ctx context.Context, ids []int) (map[int]int, error) {
	result := make(map[int]int)
	type sqlData struct {
		ItemId int `gorm:"column:item_id;type:int"`
		Count  int `gorm:"column:count;type:int"`
	}
	var listLike []sqlData
	if err := s.db.Table(model.Like{}.TableName()).Select("item_id, COUNT(item_id) as count").
		Where("item_id in (?)", ids).
		Group("item_id").Find(&listLike).Error; err != nil {
		return nil, common.ErrDB(err)
	}
	for _, item := range listLike {
		result[item.ItemId] = item.Count
	}

	return result, nil
}
