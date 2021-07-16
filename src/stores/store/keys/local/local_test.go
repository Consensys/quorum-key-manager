package local

import (
	"testing"

	testutils2 "github.com/consensys/quorum-key-manager/src/infra/log/testutils"
	"github.com/consensys/quorum-key-manager/src/stores/store/database/mock"
	mocksecrets "github.com/consensys/quorum-key-manager/src/stores/store/secrets/mock"

	"github.com/consensys/quorum-key-manager/src/stores/store/keys"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type localKeyStoreTestSuite struct {
	suite.Suite
	keyStore        keys.Store
	mockDB          *mock.MockDatabase
	mockKeys        *mock.MockKeys
	mockSecretStore *mocksecrets.MockStore
}

func TestLocalKeyStore(t *testing.T) {
	s := new(localKeyStoreTestSuite)
	suite.Run(t, s)
}

func (s *localKeyStoreTestSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	s.mockDB = mock.NewMockDatabase(ctrl)
	s.mockKeys = mock.NewMockKeys(ctrl)
	s.mockSecretStore = mocksecrets.NewMockStore(ctrl)

	s.keyStore = New(s.mockSecretStore, s.mockDB, testutils2.NewMockLogger(ctrl))
}

func (s *localKeyStoreTestSuite) TestCreate() {

}

func (s *localKeyStoreTestSuite) TestImport() {

}

func (s *localKeyStoreTestSuite) TestGet() {

}

func (s *localKeyStoreTestSuite) TestList() {

}

func (s *localKeyStoreTestSuite) TestSign() {

}
