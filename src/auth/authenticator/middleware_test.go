package authenticator

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/consensys/quorum-key-manager/src/infra/log/testutils"
	"github.com/stretchr/testify/require"

	mockauth "github.com/consensys/quorum-key-manager/src/auth/authenticator/mock"
	"github.com/consensys/quorum-key-manager/src/auth/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type testHandler struct {
	t        *testing.T
	userInfo *types.UserInfo
}

func (h *testHandler) ServeHTTP(_ http.ResponseWriter, req *http.Request) {
	userInfo := UserInfoContextFromContext(req.Context())
	if h.userInfo != nil {
		require.NotNil(h.t, h.userInfo)
		assert.Equal(h.t, h.userInfo.Groups, userInfo.Groups)
		assert.Equal(h.t, h.userInfo.Username, userInfo.Username)
	} else {
		require.Nil(h.t, h.userInfo)
	}
}

func TestMiddleware(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := testutils.NewMockLogger(ctrl)

	auth1 := mockauth.NewMockAuthenticator(ctrl)

	mid := NewMiddleware(logger, auth1)

	t.Run("authentication rejected", func(t *testing.T) {
		h := mid.Then(&testHandler{t, nil})
		auth1.EXPECT().Authenticate(gomock.Any()).Return(nil, fmt.Errorf("test invalid auth"))

		req, _ := http.NewRequest(http.MethodGet, "", nil)
		rec := httptest.NewRecorder()

		h.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusUnauthorized, rec.Code, "Status code should match")
		assert.Equal(t, []byte(`test invalid auth`), rec.Body.Bytes()[:(rec.Body.Len()-1)], "Body should match")
	})

	t.Run("authentication accepted", func(t *testing.T) {
		user := &types.UserInfo{
			Username: "test-username",
			Groups: []string{
				"group-test1",
				"group-test2",
			},
		}
		h := mid.Then(&testHandler{t, user})
		auth1.EXPECT().Authenticate(gomock.Any()).Return(user, nil)
		req, _ := http.NewRequest(http.MethodGet, "", nil)
		rec := httptest.NewRecorder()

		h.ServeHTTP(rec, req)
	})

	t.Run("authentication ignored", func(t *testing.T) {
		h := mid.Then(&testHandler{t, types.AnonymousUser})
		auth1.EXPECT().Authenticate(gomock.Any()).Return(nil, nil)

		req, _ := http.NewRequest(http.MethodGet, "", nil)
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
	})
}
