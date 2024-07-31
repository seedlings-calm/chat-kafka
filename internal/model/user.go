package model

type User struct {
	Id       int    `json:"id" gorm:"primaryKey;autoIncrement;column:id"`
	UserNo   string `json:"user_no" gorm:"column:user_no"`     // 展示的用户id
	Nickname string `json:"nickname" gorm:"column:nickname"`   // 昵称
	Mobile   string `json:"mobile" gorm:"column:mobile"`       // 手机号
	Password string `json:"password" gorm:"column:password"`   // 密码
	UserSalt string `json:"user_salt" gorm:"column:user_salt"` // 加密盐
	Status   int    `json:"status" gorm:"column:status"`       // 用户状态 1正常 2冻结
	ModelTime
}

func (u User) TableName() string {
	return TableNamePre + "_user"
}
