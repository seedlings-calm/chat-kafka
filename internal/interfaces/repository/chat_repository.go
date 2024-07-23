package repository

import (
	"strconv"

	"github.com/seedlings-calm/chat-kafka/internal/domain/model"
)

type ChatRepository interface {
	GetGroupUsers(groupID string) []string
	SaveMessage(message model.Message) error
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

func (r *MySQLChatRepository) SaveMessage(message model.Message) error {
	return nil
}
