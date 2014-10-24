package userrepo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test(t *testing.T) {
	/* create repo */
	repo := NewUserRepo()
	assert.Equal(t, 0, len(repo.FetchUsers()))

	user := NewUser("1234", "Marc")
	{
		/* create new user */
        repo.StoreUser(user)
		users := repo.FetchUsers()
		assert.Equal(t, 1, len(users))
		assert.Equal(t, user.Name, users[0].Name)
		assert.Equal(t, user.SessionId, users[0].SessionId)
		found := repo.FetchUser("1234")
		assert.Equal(t, user.SessionId, found.SessionId)
	}

	{
		/* create new user again */
        repo.StoreUser(user)
		users := repo.FetchUsers()
		assert.Equal(t, 1, len(users))
		assert.Equal(t, user.Name, users[0].Name)
		assert.Equal(t, user.SessionId, users[0].SessionId)
	}

	user2 := NewUser("4321", "Eva")
	{
		/* create another user again */
        repo.StoreUser(user2)
		users := repo.FetchUsers()
		assert.Equal(t, 2, len(users))
		assert.Equal(t, user.Name, users[0].Name)
		assert.Equal(t, user.SessionId, users[0].SessionId)
		assert.Equal(t, user2.Name, users[1].Name)
		assert.Equal(t, user2.SessionId, users[1].SessionId)
		found := repo.FetchUser("4321")
		assert.Equal(t, user2.SessionId, found.SessionId)
	}

    {
		user := NewUser("1234", "Marc")
		repo.RemoveUser(user)
		users := repo.FetchUsers()
		assert.Equal(t, 1, len(users))
		found := repo.FetchUser("4321")
		assert.Equal(t, user2.SessionId, found.SessionId)
		u := repo.FetchUser("1234")
		assert.Nil(t, u)
    }
}
