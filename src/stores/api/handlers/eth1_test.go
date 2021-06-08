package handlers

import (
	mock2 "github.com/ConsenSysQuorum/quorum-key-manager/src/stores/manager/mock"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/stores/store/eth1/mock"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"
)

type eth1HandlerTestSuite struct {
	suite.Suite

	ctrl         *gomock.Controller
	storeManager *mock2.MockManager
	keyStore     *mock.MockStore
	router       *mux.Router
}

func TestEth1Handler(t *testing.T) {
	s := new(eth1HandlerTestSuite)
	suite.Run(t, s)
}

func (s *eth1HandlerTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())

	s.storeManager = mock2.NewMockManager(s.ctrl)
	s.keyStore = mock.NewMockStore(s.ctrl)

	s.router = mux.NewRouter()
	NewStoresHandler(s.storeManager).Register(s.router)
}

func (s *eth1HandlerTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

// TODO: Unit tests
