package userrepo

import (
)

type User struct {
	SessionId string
	Name      string
}

func NewUser(sessionId string, name string) *User {
	user := new(User)
	user.SessionId = sessionId
	user.Name = name
	return user
}

type UserRepoI interface {
	FetchUser(sessionId string) *User
	FetchUsers() []User
	StoreUser(user *User)
	RemoveUser(user *User)
}

type UserRepo struct {
	users []User
}

func NewUserRepo() *UserRepo {
	userrepo := new(UserRepo)
	userrepo.users = []User{}

	return userrepo
}

func (repo *UserRepo) FetchUser(sessionId string) *User {

	for _, user := range repo.users {
		if user.SessionId == sessionId {
			return &user
		}
	}

	return nil
}

func (repo *UserRepo) FetchUsers() []User {
	return repo.users
}

func (repo *UserRepo) StoreUser(user *User) {
	found := repo.FetchUser(user.SessionId)
	if found == nil {
		repo.users = append(repo.users, *user)
	}
}

func (repo *UserRepo) RemoveUser(toBeRemoved *User) {
	newUsers := []User{}
	for _, user := range repo.users {
		if user.SessionId != toBeRemoved.SessionId {
			newUsers = append(newUsers, user)
		}
	}
	repo.users = newUsers
}
