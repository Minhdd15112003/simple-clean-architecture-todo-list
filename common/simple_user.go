package common

type SimpleUser struct {
	SQLModel
	LastName  string `json:"last_name"  gorm:"column:last_name;"`
	FirstName string `json:"first_name" gorm:"column:first_name;"`
}

func (SimpleUser) TableName() string {
	return "users"
}

func (su *SimpleUser) Mask() {
	su.SQLModel.Mask(DbTypeUser)
}
