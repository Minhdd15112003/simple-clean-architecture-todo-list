package model

import (
	"social-todo-list/common"
	"time"
)

const (
	EntityName = "UserLikeItem"
)

type Like struct {
	UserID    int                `json:"user_id" gorm:"column:user_id;type:int"`
	ItemID    int                `json:"item_id" gorm:"column:item_id;type:int"`
	CreatedAt *time.Time         `json:"created_at,omitempty"`
	User      *common.SimpleUser `json:"-" gorm:"foreignKey:UserID;"`
}

func (Like) TableName() string {
	return "user_like_items"
}

func ErrCannotLikeItem(err error) *common.AppError {
	return common.NewCustomError(err, "Cannot like this item", "ErrCannotLike")
}

func ErrCannotUnLikeItem(err error) *common.AppError {
	return common.NewCustomError(err, "Cannot dislike this item", "ErrCannotUnLike")
}

func ErrCannotDidNotLikeItem(err error) *common.AppError {
	return common.NewCustomError(err, "You have not liked this item", "ErrCannotDidNotLike")
}
