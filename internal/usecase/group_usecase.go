package usecase

import domain "github.com/seedlings-calm/chat-kafka/internal/domain/group"

type GroupRepository interface {
	GetGroup(id string) (*domain.Group, error)
	SaveGroup(group *domain.Group) error
}

type GroupUseCase struct {
	groupRepo GroupRepository
}

func NewGroupUseCase(repo GroupRepository) *GroupUseCase {
	return &GroupUseCase{groupRepo: repo}
}

func (uc *GroupUseCase) CreateGroup(id, name string) error {
	group := domain.NewGroup(id, name)
	return uc.groupRepo.SaveGroup(group)
}

func (uc *GroupUseCase) AddUserToGroup(userID, groupID string) error {
	group, err := uc.groupRepo.GetGroup(groupID)
	if err != nil {
		return err
	}
	if err := group.AddUser(userID); err != nil {
		return err
	}
	return uc.groupRepo.SaveGroup(group)
}
