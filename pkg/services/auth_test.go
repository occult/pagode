package services

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestAuthClient_Auth(t *testing.T) {
	assertNoAuth := func() {
		_, err := c.Auth.GetAuthenticatedUserID(ctx)
		assert.True(t, errors.Is(err, NotAuthenticatedError{}))
		_, err = c.Auth.GetAuthenticatedUser(ctx)
		assert.True(t, errors.Is(err, NotAuthenticatedError{}))
	}

	assertNoAuth()

	err := c.Auth.Login(ctx, usr.ID)
	require.NoError(t, err)

	uid, err := c.Auth.GetAuthenticatedUserID(ctx)
	require.NoError(t, err)
	assert.Equal(t, usr.ID, uid)

	u, err := c.Auth.GetAuthenticatedUser(ctx)
	require.NoError(t, err)
	assert.Equal(t, u.ID, usr.ID)

	err = c.Auth.Logout(ctx)
	require.NoError(t, err)

	assertNoAuth()
}
