package userrepo

import (
	"fmt"
	"time"
)

type User struct {
	SessionId        string
	Name             string
	SentMessages     []Message
	ReceivedMessages []Message
}

func NewUser(sessionId string, name string) *User {
	user := new(User)
	user.SessionId = sessionId
	user.Name = name
	user.SentMessages = []Message{}
	user.ReceivedMessages = []Message{}
	return user
}

func (user *User) String() string {
	return fmt.Sprintf("{ User: SessionId: %s, Name: %s }",
		user.SessionId, user.Name)
}

func (user *User) AddMsgReceived(msg *Message) {
	user.ReceivedMessages = append(user.ReceivedMessages, *msg)
}

func (user *User) AddMsgSent(msg *Message) {
	user.SentMessages = append(user.SentMessages, *msg)
}

type Message struct {
	OriginatorSessionId string
	MessageText         string
	Timestamp           time.Time
}

func NewMessage(originatorSessionId string, messageText string) *Message {
	msg := new(Message)
	msg.OriginatorSessionId = originatorSessionId
	msg.MessageText = messageText
	msg.Timestamp = time.Now()

	return msg
}

func (msg *Message) String() string {
	return fmt.Sprintf("{ Message: OriginatorSessionId: %s, MessageText: %s, Timestamp: %s }",
		msg.OriginatorSessionId, msg.MessageText, msg.Timestamp.String())
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
