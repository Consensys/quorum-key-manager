package authenticator

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/consensys/quorum-key-manager/src/infra/log/testutils"

	mockauth "github.com/consensys/quorum-key-manager/src/auth/authenticator/mock"
	mockmanager "github.com/consensys/quorum-key-manager/src/auth/manager/mock"
	"github.com/consensys/quorum-key-manager/src/auth/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type testHandler struct {
	t *testing.T
}

func (h *testHandler) ServeHTTP(_ http.ResponseWriter, req *http.Request) {
	reqCtx := UserContextFromContext(req.Context())
	assert.NotNil(h.t, reqCtx, "UserContext should have been set on context")
}

func TestMiddleware(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := testutils.NewMockLogger(ctrl)

	auth1 := mockauth.NewMockAuthenticator(ctrl)
	policyMngr := mockmanager.NewMockManager(ctrl)

	mid := NewMiddleware(policyMngr, logger, auth1)

	t.Run("authentication rejected", func(t *testing.T) {
		h := mid.Then(&testHandler{t})
		auth1.EXPECT().Authenticate(gomock.Any()).Return(nil, fmt.Errorf("test invalid auth"))

		req, _ := http.NewRequest(http.MethodGet, "", nil)
		rec := httptest.NewRecorder()

		h.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusUnauthorized, rec.Code, "Status code should match")
		assert.Equal(t, []byte(`test invalid auth`), rec.Body.Bytes()[:(rec.Body.Len()-1)], "Body should match")
	})

	t.Run("authentication accepted", func(t *testing.T) {
		h := mid.Then(&testHandler{t})
		user := &types.UserInfo{
			Username: "test-username",
			Groups: []string{
				"group-test1",
				"group-test2",
			},
		}
		auth1.EXPECT().Authenticate(gomock.Any()).Return(user, nil)

		group1 := &types.Group{
			Policies: []string{
				"policy1.A",
				"policy1.B",
			},
		}
		policyMngr.EXPECT().Group(gomock.Any(), "group-test1").Return(group1, nil)
		policyMngr.EXPECT().Policy(gomock.Any(), "policy1.A").Return(&types.Policy{}, nil)
		policyMngr.EXPECT().Policy(gomock.Any(), "policy1.B").Return(&types.Policy{}, nil)

		group2 := &types.Group{
			Policies: []string{
				"policy2.A",
			},
		}
		policyMngr.EXPECT().Group(gomock.Any(), "group-test2").Return(group2, nil)
		policyMngr.EXPECT().Policy(gomock.Any(), "policy2.A").Return(&types.Policy{}, nil)

		policyMngr.EXPECT().Group(gomock.Any(), "system:authenticated").Return(nil, fmt.Errorf("not found"))

		req, _ := http.NewRequest(http.MethodGet, "", nil)
		rec := httptest.NewRecorder()

		h.ServeHTTP(rec, req)
	})

	t.Run("authentication ignored", func(t *testing.T) {
		h := mid.Then(&testHandler{t})
		auth1.EXPECT().Authenticate(gomock.Any()).Return(nil, nil)

		policyMngr.EXPECT().Group(gomock.Any(), "system:unauthenticated").Return(nil, fmt.Errorf("not found"))

		req, _ := http.NewRequest(http.MethodGet, "", nil)
		rec := httptest.NewRecorder()

		h.ServeHTTP(rec, req)
	})
}
