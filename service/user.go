package service

import (
	"fmt"
)

type User struct {
	Id string
}

type UserService struct {
	DbHits int
	users  []*User
}

func (s *UserService) GetOne(id string) (*User, error) {
	s.DbHits++

	for _, u := range s.users {
		if u.Id == id {
			return u, nil
		}
	}

	return nil, fmt.Errorf("user %s not found", id)
}

func NewUserService(usersCount int) *UserService {
	userService := &UserService{
		users: []*User{},
	}

	for i := 0; i < usersCount; i++ {
		userService.users = append(userService.users, &User{Id: fmt.Sprintf("user_%d", i+1)})
	}
	return userService
}
