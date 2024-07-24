package usecase

import domain "github.com/seedlings-calm/chat-kafka/internal/domain/user"

type UserRepository interface {
	GetUser(id string) (*domain.User, error)
	SaveUser(user *domain.User) error
}

type UserUseCase struct {
	userRepo UserRepository
}

func NewUserUseCase(repo UserRepository) *UserUseCase {
	return &UserUseCase{userRepo: repo}
}

func (uc *UserUseCase) AddUserToGroup(userID, groupID string) error {
	user, err := uc.userRepo.GetUser(userID)
	if err != nil {
		return err
	}
	user.AddGroup(groupID)
	return uc.userRepo.SaveUser(user)
}
