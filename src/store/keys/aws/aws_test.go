package aws

import (
	"context"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/log"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/infra/aws/mocks"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities/testutils"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/keys"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

const (
	id = "my-key"
)

type awsKeyStoreTestSuite struct {
	suite.Suite
	mockVault *mocks.MockKmsClient
	keyStore  keys.Store
}

func TestAWSKeyStore(t *testing.T) {
	s := new(awsKeyStoreTestSuite)
	suite.Run(t, s)
}

func (s *awsKeyStoreTestSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	s.mockVault = mocks.NewMockKmsClient(ctrl)
	s.keyStore = New(s.mockVault, log.DefaultLogger())
}

// TestCreate Key creation test cases
func (s *awsKeyStoreTestSuite) TestCreate() {
	ctx := context.Background()
	attributes := testutils.FakeAttributes()
	algorithm := testutils.FakeAlgorithm()
	creationDate := time.Now()
	deletionDate := creationDate.AddDate(1, 0, 0)

	retCreateKey := kms.CreateKeyOutput{
		KeyMetadata: &kms.KeyMetadata{
			CreationDate: &creationDate,
			DeletionDate: &deletionDate,
			Enabled:      aws.Bool(true),
			KeyId:        aws.String("someId"),
		},
	}

	retGetPub := kms.GetPublicKeyOutput{}

	s.T().Run("should create a new key successfully", func(t *testing.T) {
		s.mockVault.EXPECT().CreateKey(gomock.Any(), id, gomock.Any(), gomock.Any()).
			Return(&retCreateKey, nil)
		s.mockVault.EXPECT().GetPublicKey(gomock.Any(), id).
			Return(&retGetPub, nil)

		key, err := s.keyStore.Create(ctx, id, algorithm, attributes)

		assert.NoError(t, err)
		assert.NotEmpty(t, key.Metadata.CreatedAt)
		assert.NotEmpty(t, key.Metadata.DeletedAt)
		assert.False(t, key.Metadata.Disabled)

	})
}

// TestSign Signature test cases
func (s *awsKeyStoreTestSuite) TestSign() {
	ctx := context.Background()
	msg := []byte("some sample message")
	myKeyId := "the_id"

	retSign := kms.SignOutput{
		KeyId:     &myKeyId,
		Signature: []byte("signature"),
	}

	s.T().Run("should sign a sample message", func(t *testing.T) {
		s.mockVault.EXPECT().Sign(gomock.Any(), id, msg).
			Return(&retSign, nil)

		signature, err := s.keyStore.Sign(ctx, id, msg)

		assert.NoError(t, err)
		assert.NotEmpty(t, signature)

	})
}
