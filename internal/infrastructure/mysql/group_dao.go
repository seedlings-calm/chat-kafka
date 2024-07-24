package mysql

import (
	"errors"

	domain "github.com/seedlings-calm/chat-kafka/internal/domain/group"
)

type GroupRepositoryImpl struct {
	// 假设我们用一个简单的 map 来模拟数据库
	groups map[string]*domain.Group
}

func NewGroupRepositoryImpl() *GroupRepositoryImpl {
	return &GroupRepositoryImpl{
		groups: make(map[string]*domain.Group),
	}
}

func (repo *GroupRepositoryImpl) GetGroup(id string) (*domain.Group, error) {
	group, exists := repo.groups[id]
	if !exists {
		return nil, errors.New("group not found")
	}
	return group, nil
}

func (repo *GroupRepositoryImpl) SaveGroup(group *domain.Group) error {
	repo.groups[group.ID] = group
	return nil
}
