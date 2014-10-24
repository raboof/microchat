package userrepo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserFound(t *testing.T) {
	repo := NewUserRepo()

	user := NewUser("1234", "Marc")

	repo.StoreUser(user)

	userAgain := repo.FetchUser("1234")

	assert.Equal(t, user.SessionId, userAgain.SessionId)
	assert.Equal(t, user.Name, userAgain.Name)
}

func TestUserNotFound(t *testing.T) {
	repo := NewUserRepo()
	user := NewUser("4321", "Eva")
	repo.StoreUser(user)

	userAgain := repo.FetchUser("1234")
	assert.Nil(t, userAgain)
}
