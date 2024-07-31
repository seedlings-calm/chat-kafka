package model

// 群
type Groups struct {
	Id         int    `json:"id" gorm:"primaryKey;autoIncrement;column:id"`
	Title      string `json:"title" gorm:"column:title"`             // 群名
	Uuid       string `json:"uuid" gorm:"column:uuid"`               // 群号
	CreateUser int    `json:"create_user" gorm:"column:create_user"` //创建者
	Leader     int    `json:"leader" gorm:"column:leader"`           //群主
	ModelTime
}

func (ug Groups) TableName() string {
	return TableNamePre + "_group"
}

// 群成员
type GroupMembers struct {
	Id      int `json:"id" gorm:"primaryKey;autoIncrement;column:id"`
	GroupId int `json:"group_id"`
	UserId  int `json:"user_id"`
	ModelTime
}

func (ug GroupMembers) TableName() string {
	return TableNamePre + "_group_members"
}
