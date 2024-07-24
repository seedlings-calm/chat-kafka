package domain

type User struct {
	ID       string
	Name     string
	GroupIDs []string
}

func (u *User) AddGroup(groupID string) {
	u.GroupIDs = append(u.GroupIDs, groupID)
}

func NewUser(id, name string) *User {
	return &User{
		ID:   id,
		Name: name,
	}
}
