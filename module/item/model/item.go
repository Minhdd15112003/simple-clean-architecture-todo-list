package model

import (
	"errors"
	"social-todo-list/common"
	"strings"
)

var (
	ErrTitleCannotBeEmpty = errors.New("title cannot be empty")
	ErrItemIsDeleted      = errors.New("item is deleted")
)

const (
	EntityName = "Item"
)

type TodoItem struct {
	//Embedding models
	common.SQLModel
	Title       string        `json:"title"       gorm:"column:title;"`
	Description string        `json:"description" gorm:"column:description;"`
	Status      string        `json:"status"      gorm:"column:status;"`
	Image       *common.Image `json:"image"       gorm:"column:image;type:json"`
	User_id     int           `json:"-" gorm:"column:user_id;type:int"`
	// User        model.User    `json:"user,omitempty" gorm:"column:user";`
}

func (TodoItem) TableName() string {
	return "todo_items"
}

type TodoItemCreation struct {
	Id          int           `json:"id"          gorm:"column:id;"`
	Title       string        `json:"title"       gorm:"column:title;"`
	Description string        `json:"description" gorm:"column:description;"`
	Image       *common.Image `json:"image"       gorm:"column:image;type:json"`
	User_id     int           `json:"-" gorm:"column:user_id;type:int"`
}

func (i *TodoItemCreation) Validate() error {
	i.Title = strings.TrimSpace(i.Title)
	if i.Title == "" {
		return ErrTitleCannotBeEmpty
	}
	return nil
}

func (TodoItemCreation) TableName() string {
	return TodoItem{}.TableName()
}

type TodoItemUpdation struct {
	Title       *string `json:"title"       gorm:"column:title;"`
	Description *string `json:"description" gorm:"column:description;"`
	Status      string  `json:"status"      gorm:"column:status;"`
}

func (TodoItemUpdation) TableName() string {
	return TodoItem{}.TableName()
}
