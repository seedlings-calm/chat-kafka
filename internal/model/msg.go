package model

// 消息
type Messages struct {
	Id      int    `json:"id" gorm:"primaryKey;autoIncrement;column:id"`
	Uuid    string `json:"uuid" gorm:"column:uuid"`        // 消息ID
	Content string `json:"content" gorm:"type:text"`       //消息内容
	FromId  int    `json:"fromId" gorm:"column:from_id"`   //发送人
	ToId    int    `json:"toId" gorm:"column:to_id"`       //接收人
	GroupId int    `json:"groupId" gorm:"column:group_id"` //群组ID
	ModelTime
}

func (m Messages) TableName() string {
	return TableNamePre + "_messages"
}

type MessagesLog struct{}

func (ml MessagesLog) TableName() string {
	return TableNamePre + "_messages_log"
}
