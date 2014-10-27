package main

import (
	"github.com/raboof/microchat/userrepo"
	"github.com/stretchr/testify/assert"
	"testing"
)

type FakeForwarder struct {
	ForwardUserLoggedInCalled  bool
	ForwardUserLoggedOutCalled bool
	ForwardMsgSentCalled       bool
}

func (this *FakeForwarder) ForwardUserLoggedIn(user userrepo.User) {
	this.ForwardUserLoggedInCalled = true
}

func (this *FakeForwarder) ForwardUserLoggedOut(user userrepo.User) {
	this.ForwardUserLoggedOutCalled = true
}

func (this *FakeForwarder) ForwardMsgSent(msg userrepo.Message) {
	this.ForwardMsgSentCalled = true
}

func TestUserCreatedEvent(t *testing.T) {
	repo := userrepo.NewUserRepo()
	assert.Equal(t, 0, len(repo.FetchUsers()))
	forwarder := new(FakeForwarder)

	handle := handleEvent(repo, forwarder)
	handle("key", "UserLoggedIn,Marc,12345", "user", 1, 1)
	assert.Equal(t, 1, len(repo.FetchUsers()))
	found, exists := repo.FetchUser("12345")
	assert.True(t, exists)
	assert.Equal(t, "12345", found.SessionId)
	assert.Equal(t, "Marc", found.Name)

	assert.True(t, forwarder.ForwardUserLoggedInCalled)
	assert.False(t, forwarder.ForwardUserLoggedOutCalled)
}

func TestUserRemovedEvent(t *testing.T) {
	repo := userrepo.NewUserRepo()
	assert.Equal(t, 0, len(repo.FetchUsers()))
	repo.StoreUser(userrepo.NewUser("12345", "Marc"))
	assert.Equal(t, 1, len(repo.FetchUsers()))
	forwarder := new(FakeForwarder)

	handle := handleEvent(repo, forwarder)
	handle("key", "UserLoggedOut,Marc,12345", "user", 1, 1)
	assert.Equal(t, 0, len(repo.FetchUsers()))
	_, exists := repo.FetchUser("12345")
	assert.False(t, exists)
	assert.True(t, forwarder.ForwardUserLoggedOutCalled)
	assert.False(t, forwarder.ForwardUserLoggedInCalled)
}

func TestUnsupportedEvent(t *testing.T) {
	repo := userrepo.NewUserRepo()
	assert.Equal(t, 0, len(repo.FetchUsers()))
	repo.StoreUser(userrepo.NewUser("12345", "Marc"))
	assert.Equal(t, 1, len(repo.FetchUsers()))
	forwarder := new(FakeForwarder)

	handle := handleEvent(repo, forwarder)
	handle("key", "UserCreated,Marc,12345", "user", 1, 1)
	assert.Equal(t, 1, len(repo.FetchUsers()))

	assert.False(t, forwarder.ForwardUserLoggedInCalled)
	assert.False(t, forwarder.ForwardUserLoggedOutCalled)
}
