package authenticator

import (
	"fmt"
	"testing"

	"github.com/consensys/quorum-key-manager/src/auth/authenticator/mock"
	"github.com/consensys/quorum-key-manager/src/auth/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestFirst(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	auth1 := mock.NewMockAuthenticator(ctrl)
	auth2 := mock.NewMockAuthenticator(ctrl)

	auth := First(auth1, auth2)

	t.Run("1st authenticator accepts request", func(t *testing.T) {
		authUser := &types.UserInfo{}
		auth1.EXPECT().Authenticate(gomock.Any()).Return(&types.UserInfo{}, nil)
		u, err := auth.Authenticate(nil)
		assert.NoError(t, err, "Authenticate should not error")
		assert.Equal(t, authUser, u, "UserInfo should match")
	})

	t.Run("1st authenticator rejects request", func(t *testing.T) {
		authErr := fmt.Errorf("invalid auth")
		auth1.EXPECT().Authenticate(gomock.Any()).Return(nil, authErr)
		u, err := auth.Authenticate(nil)
		assert.Error(t, err, "Authenticate should not error")
		assert.Equal(t, authErr, err, "Error should match")
		assert.Nil(t, u)
	})

	t.Run("1st authenticator ignores request", func(t *testing.T) {
		authUser := &types.UserInfo{}
		auth1.EXPECT().Authenticate(gomock.Any()).Return(nil, nil)
		auth2.EXPECT().Authenticate(gomock.Any()).Return(authUser, nil)
		u, err := auth.Authenticate(nil)
		assert.NoError(t, err, "Authenticate should not error")
		assert.Equal(t, authUser, u, "UserInfo should match")
	})

	t.Run("Both authenticator ignores request", func(t *testing.T) {
		auth1.EXPECT().Authenticate(gomock.Any()).Return(nil, nil)
		auth2.EXPECT().Authenticate(gomock.Any()).Return(nil, nil)
		u, err := auth.Authenticate(nil)
		assert.NoError(t, err, "Authenticate should not error")
		assert.Nil(t, u, "UserInfo should match")
	})
}
