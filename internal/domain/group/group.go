package domain

import "errors"

// 数据库表
type Group struct {
	ID    string
	Name  string
	Users []string
}

func (g *Group) AddUser(userID string) error {
	for _, id := range g.Users {
		if id == userID {
			return errors.New("user already in group")
		}
	}
	g.Users = append(g.Users, userID)
	return nil
}

func NewGroup(id, name string) *Group {
	return &Group{
		ID:   id,
		Name: name,
	}
}
