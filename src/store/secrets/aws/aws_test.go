package aws

import (
	"context"
	"fmt"
	"testing"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/log"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/infra/aws/mocks"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities/testutils"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/secrets"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type awsSecretStoreTestSuite struct {
	suite.Suite
	mockVault   *mocks.MockSecretsManagerClient
	secretStore secrets.Store
}

func TestAwsSecretStore(t *testing.T) {
	s := new(awsSecretStoreTestSuite)
	suite.Run(t, s)
}

func (s *awsSecretStoreTestSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	s.mockVault = mocks.NewMockSecretsManagerClient(ctrl)

	s.secretStore = New(s.mockVault, log.DefaultLogger())
}

func (s *awsSecretStoreTestSuite) TestSet() {
	ctx := context.Background()
	id := "my-secret1"
	version := "2"
	value := "my-value1"
	attributes := testutils.FakeAttributes()

	createInput := &secretsmanager.CreateSecretInput{
		Name:         &id,
		SecretString: &value,
	}

	createOutput := &secretsmanager.CreateSecretOutput{
		Name:      &id,
		VersionId: &version,
	}

	fakeSecretsTags := []*secretsmanager.Tag{}

	for key, value := range attributes.Tags {
		k, v := key, value
		var in secretsmanager.Tag = secretsmanager.Tag{
			Key:   &k,
			Value: &v,
		}
		fakeSecretsTags = append(fakeSecretsTags, &in)
	}

	tagInput := &secretsmanager.TagResourceInput{
		SecretId: &id,
		Tags:     fakeSecretsTags,
	}

	descSecretInput := &secretsmanager.DescribeSecretInput{
		SecretId: &id,
	}

	currentMark := "AWSCURRENT"
	versionID2stages := map[string][]*string{
		version: {&currentMark},
	}

	descSecretOutput := &secretsmanager.DescribeSecretOutput{
		Name:               &id,
		VersionIdsToStages: versionID2stages,
		Tags:               fakeSecretsTags,
	}

	s.T().Run("should set a new secret successfully", func(t *testing.T) {
		s.mockVault.EXPECT().CreateSecret(gomock.Any(), createInput).Return(createOutput, nil)
		s.mockVault.EXPECT().TagSecretResource(gomock.Any(), tagInput).Return(&secretsmanager.TagResourceOutput{}, nil)
		s.mockVault.EXPECT().DescribeSecret(gomock.Any(), descSecretInput).Return(descSecretOutput, nil)

		secret, err := s.secretStore.Set(ctx, id, value, attributes)

		assert.NoError(t, err)
		assert.Equal(t, value, secret.Value)

		assert.ObjectsAreEqual(attributes.Tags, secret.Tags)
		assert.Equal(t, version, secret.Metadata.Version)
		assert.False(t, secret.Metadata.Disabled)
		assert.True(t, secret.Metadata.ExpireAt.IsZero())
		assert.True(t, secret.Metadata.DeletedAt.IsZero())
	})

	s.T().Run("should fail with same error if write fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("error")
		s.mockVault.EXPECT().CreateSecret(gomock.Any(), createInput).Return(&secretsmanager.CreateSecretOutput{}, expectedErr)
		s.mockVault.EXPECT().TagSecretResource(gomock.Any(), tagInput).Return(&secretsmanager.TagResourceOutput{}, nil)
		s.mockVault.EXPECT().DescribeSecret(gomock.Any(), descSecretInput).Return(descSecretOutput, nil)

		secret, err := s.secretStore.Set(ctx, id, value, attributes)

		assert.Nil(t, secret)
		assert.Equal(t, expectedErr, err)
	})
}
