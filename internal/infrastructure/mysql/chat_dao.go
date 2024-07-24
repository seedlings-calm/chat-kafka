package mysql

import (
	"strconv"
)

type ChatRepository interface {
	GetGroupUsers(groupID string) []string
	SaveMessage(message string) error
}
type MySQLChatRepository struct {
}

func NewMySQLChatRepository() ChatRepository {
	return &MySQLChatRepository{}
}

func (r *MySQLChatRepository) GetGroupUsers(groupID string) []string {
	var users []string
	for i := 1; i <= 10; i++ {
		users = append(users, strconv.Itoa(i))
	}
	users = append(users, "35")
	return users
}

func (r *MySQLChatRepository) SaveMessage(message string) error {
	return nil
}
