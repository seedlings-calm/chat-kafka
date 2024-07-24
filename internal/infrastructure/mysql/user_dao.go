package mysql

import (
	"errors"

	domain "github.com/seedlings-calm/chat-kafka/internal/domain/user"
)

type UserRepositoryImpl struct {
	// 假设我们用一个简单的 map 来模拟数据库
	users map[string]*domain.User
}

func NewUserRepositoryImpl() *UserRepositoryImpl {
	return &UserRepositoryImpl{
		users: make(map[string]*domain.User),
	}
}

func (repo *UserRepositoryImpl) GetUser(id string) (*domain.User, error) {
	user, exists := repo.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (repo *UserRepositoryImpl) SaveUser(user *domain.User) error {
	repo.users[user.ID] = user
	return nil
}
