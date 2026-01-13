package common

import "time"

type SQLModel struct {
	Id        int        `json:"id"                   gorm:"column:id;"`
	CreatedAt *time.Time `json:"created_at"           gorm:"column:created_at;"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" gorm:"column:updated_at;"`
	// tag name trong struct giúp giao thức http nhận biết nên nhận vào và chả ra gì
	// dùng *time.Time để trỏ đến giá trị của ngày hiện tại
}
