package userrepo

;

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserFound(t *testing.T) {
	repo := NewUserRepo()

	user := NewUser("1234", "Marc");

	repo.StoreUser(user)

	userAgain := repo.FetchUser("1234");

	assert.Equal(t, user.sessionId, userAgain.sessionId)
	assert.Equal(t, user.name, userAgain.name)
}

func TestUserNotFound(t *testing.T) {
	repo := NewUserRepo()
	user := NewUser("4321", "Eva");
	repo.StoreUser(user)

	userAgain := repo.FetchUser("1234");
	assert.Nil(t, userAgain)
}
