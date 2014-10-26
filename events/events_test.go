package events

import (
	"github.com/raboof/microchat/userrepo"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserCreatedEvent(t *testing.T) {
	repo := userrepo.NewUserRepo()
	assert.Equal(t, 0, len(repo.FetchUsers()))

	eventListener := NewDomainEventListener(repo)

	eventListener.HandleEvent("key", "UserLoggedIn,Marc,12345", "user", 1, 1)
	assert.Equal(t, 1, len(repo.FetchUsers()))
	found, exists := repo.FetchUser("12345")
	assert.True(t, exists)
	assert.Equal(t, "12345", found.SessionId)
	assert.Equal(t, "Marc", found.Name)
}

func TestUserRemovedEvent(t *testing.T) {
	repo := userrepo.NewUserRepo()
	assert.Equal(t, 0, len(repo.FetchUsers()))
	repo.StoreUser(userrepo.NewUser("12345", "Marc"))
	assert.Equal(t, 1, len(repo.FetchUsers()))

	eventListener := NewDomainEventListener(repo)

	eventListener.HandleEvent("key", "UserLoggedOut,Marc,12345", "user", 1, 1)
	assert.Equal(t, 0, len(repo.FetchUsers()))
	_, exists := repo.FetchUser("12345")
	assert.False(t, exists)
}

func TestUnsupportedEvent(t *testing.T) {
	repo := userrepo.NewUserRepo()
	assert.Equal(t, 0, len(repo.FetchUsers()))
	repo.StoreUser(userrepo.NewUser("12345", "Marc"))
	assert.Equal(t, 1, len(repo.FetchUsers()))

	eventListener := NewDomainEventListener(repo)

	eventListener.HandleEvent("key", "UserCreated,Marc,12345", "user", 1, 1)
	assert.Equal(t, 1, len(repo.FetchUsers()))
}
