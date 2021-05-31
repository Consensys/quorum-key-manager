package handlers

import (
	"testing"

	mockstoremanager "github.com/ConsenSysQuorum/quorum-key-manager/src/stores/manager/mock"
	mocketh1 "github.com/ConsenSysQuorum/quorum-key-manager/src/stores/store/eth1/mock"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"
)

type eth1HandlerTestSuite struct {
	suite.Suite

	ctrl         *gomock.Controller
	storeManager *mockstoremanager.MockManager
	keyStore     *mocketh1.MockStore
	router       *mux.Router
}

func TestEth1Handler(t *testing.T) {
	s := new(eth1HandlerTestSuite)
	suite.Run(t, s)
}

func (s *eth1HandlerTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())

	s.storeManager = mockstoremanager.NewMockManager(s.ctrl)
	s.keyStore = mocketh1.NewMockStore(s.ctrl)

	s.router = mux.NewRouter()
	NewStoresHandler(s.storeManager).Register(s.router)
}

func (s *eth1HandlerTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

// TODO: Unit tests
