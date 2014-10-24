package userrepo

import (
	"log"
)

type User struct {
	SessionId string
	Name      string
}

func NewUser(sessionId string, name string) *User {
	user := new(User)
	user.SessionId = sessionId
	user.Name = name
	log.Printf("User with session %s and name %s created", sessionId, name)
	return user
}

type UserRepoI interface {
	FetchUser(sessionId string) *User
	FetchUsers() []User
	StoreUser(user *User)
}

type UserRepo struct {
	users []User
}

func NewUserRepo() *UserRepo {
	userrepo := new(UserRepo)
	userrepo.users = []User{}

	log.Printf("Repo created")

	return userrepo
}

func (repo *UserRepo) FetchUser(sessionId string) *User {

	log.Printf("Repo sise: %d\n", len(repo.users))
	for _, user := range repo.users {
		log.Printf("*** Session %s found: %s", user.SessionId, user.Name)
		if user.SessionId == sessionId {
			log.Printf("User with session %s found: %s", sessionId, user.Name)
			return &user
		}
	}
	log.Printf("User with session %s NOT found", sessionId)

	log.Printf("Repo sise: %d\n", len(repo.users))
	for _, user := range repo.users {
		log.Printf("*** Session %s found: %s", user.SessionId, user.Name)
		if user.SessionId == sessionId {
			log.Printf("User with session %s found: %s", sessionId, user.Name)
			return &user
		}
	}
	log.Printf("User with session %s NOT found", sessionId)

	return nil
}

func (repo *UserRepo) FetchUsers() []User {
	log.Printf("%s Users found", len(repo.users))
	return repo.users
}

func (repo *UserRepo) StoreUser(user *User) {
	repo.users = append(repo.users, *user)
	log.Printf("User %s added", user.Name)
}
