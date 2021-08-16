package keys

import (
	"context"
	"testing"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/infra/log/testutils"
	mock2 "github.com/consensys/quorum-key-manager/src/stores/database/mock"
	testutils2 "github.com/consensys/quorum-key-manager/src/stores/entities/testutils"
	"github.com/consensys/quorum-key-manager/src/stores/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestEncryptKey(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mock.NewMockKeyStore(ctrl)
	db := mock2.NewMockKeys(ctrl)
	logger := testutils.NewMockLogger(ctrl)
	data := []byte("0x123")
	result := []byte("0x456")

	connector := NewConnector(store, db, logger)

	t.Run("should encrypt data successfully", func(t *testing.T) {
		key := testutils2.FakeKey()

		store.EXPECT().Encrypt(gomock.Any(), key.ID, data).Return(result, nil)

		rResult, err := connector.Encrypt(ctx, key.ID, data)

		assert.NoError(t, err)
		assert.Equal(t, rResult, result)
	})

	t.Run("should fail to encrypt data if encrypt fails", func(t *testing.T) {
		key := testutils2.FakeKey()
		expectedErr := errors.UnauthorizedError("not authorized")

		store.EXPECT().Encrypt(gomock.Any(), key.ID, data).Return(nil, expectedErr)

		_, err := connector.Encrypt(ctx, key.ID, data)

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})
}
