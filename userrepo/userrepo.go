package userrepo

import (
	"fmt"
	"sync"
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
	user.SentMessages = make([]Message, 10)
	user.ReceivedMessages = make([]Message, 10)
	return user
}

func (user User) String() string {
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

func (msg Message) String() string {
	return fmt.Sprintf("{ Message: OriginatorSessionId: %s, MessageText: %s, Timestamp: %s }",
		msg.OriginatorSessionId, msg.MessageText, msg.Timestamp.String())
}

type UserRepoI interface {
	FetchUser(sessionId string) (User, bool)
	FetchUsers() []User
	StoreUser(user *User)
	RemoveUser(user *User)
}

type UserRepo struct {
	m     sync.Mutex
	users map[string]User
}

func NewUserRepo() *UserRepo {
	userrepo := new(UserRepo)
	userrepo.users = make(map[string]User)

	return userrepo
}

func (this UserRepo) FetchUser(sessionId string) (User, bool) {
	this.m.Lock()
	defer this.m.Unlock()

	user, ok := this.users[sessionId]
	return user, ok
}

func (this UserRepo) FetchUsers() []User {
	this.m.Lock()
	defer this.m.Unlock()

	list := make([]User, 0, len(this.users))
	for _, user := range this.users {
		list = append(list, user)
	}
	return list
}

func (this *UserRepo) StoreUser(user *User) {
	this.m.Lock()
	defer this.m.Unlock()

	this.users[user.SessionId] = *user
}

func (this *UserRepo) RemoveUser(toBeRemoved *User) {
	this.m.Lock()
	defer this.m.Unlock()

	delete(this.users, toBeRemoved.SessionId)
}
