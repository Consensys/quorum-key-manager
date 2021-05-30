package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	mockstoremanager "github.com/ConsenSysQuorum/quorum-key-manager/src/services/stores/manager/mock"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type storeHandlerTestSuite struct {
	suite.Suite

	storeManager *mockstoremanager.MockManager
	ctrl         *gomock.Controller
	router       *mux.Router
}

func TestStoresHandler(t *testing.T) {
	s := new(storeHandlerTestSuite)
	suite.Run(t, s)
}

func (s *storeHandlerTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())

	s.storeManager = mockstoremanager.NewMockManager(s.ctrl)

	s.router = mux.NewRouter()
	NewStoresHandler(s.storeManager).Register(s.router)
}

func (s *storeHandlerTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *storeHandlerTestSuite) TestTest() {
	s.Run("should execute request successfully with version", func() {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodGet, "/stores/test", nil)

		s.router.ServeHTTP(rw, httpRequest)

		expectedBody := []byte("OK")
		assert.Equal(s.T(), expectedBody, rw.Body.Bytes())
		assert.Equal(s.T(), http.StatusOK, rw.Code)
	})
}
