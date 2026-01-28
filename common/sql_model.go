package common

import "time"

type SQLModel struct {
	Id        int        `json:"-"                   gorm:"column:id;"`
	FakeID    *UID       `json:"id" gorm:"-"`
	CreatedAt *time.Time `json:"created_at"           gorm:"column:created_at;"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" gorm:"column:updated_at;"`
	// tag name trong struct giúp giao thức http nhận biết nên nhận vào và chả ra gì
	// dùng *time.Time để trỏ đến giá trị của ngày hiện tại
}

func (sqlModel *SQLModel) Mask(dbType DbType) {
	uid := NewUID(uint32(sqlModel.Id), int(dbType), 1)
	sqlModel.FakeID = &uid
}
