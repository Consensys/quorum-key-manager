package handlers

import (
	"testing"

	"github.com/ConsenSysQuorum/quorum-key-manager/src/core/mocks"
	mockstoremanager "github.com/ConsenSysQuorum/quorum-key-manager/src/core/store-manager/mock"
	mocketh1 "github.com/ConsenSysQuorum/quorum-key-manager/src/store/eth1/mock"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"
)

const (
	eth1StoreName = "Eth1Store"
)

type eth1HandlerTestSuite struct {
	suite.Suite
	keyStore *mocketh1.MockStore
	router   *mux.Router
}

func TestEth1Handler(t *testing.T) {
	s := new(eth1HandlerTestSuite)
	suite.Run(t, s)
}

func (s *eth1HandlerTestSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	backend := mocks.NewMockBackend(ctrl)
	storeManager := mockstoremanager.NewMockStoreManager(ctrl)
	s.keyStore = mocketh1.NewMockStore(ctrl)

	backend.EXPECT().StoreManager().Return(storeManager).AnyTimes()
	storeManager.EXPECT().GetEth1StoreByAddr(gomock.Any(), eth1StoreName).Return(s.keyStore, nil).AnyTimes()

	s.router = NewKeysHandler(backend)
}

// TODO: Unit tests
