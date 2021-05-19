package aws

import (
	"context"
	"fmt"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities"
	"testing"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/log"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/infra/aws/mocks"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities/testutils"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/secrets"
	"github.com/aws/aws-sdk-go/aws/awserr"
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

	putSecretInput := &secretsmanager.PutSecretValueInput{
		SecretId:     &id,
		SecretString: &value,
	}

	s.T().Run("should update secret if already exists", func(t *testing.T) {

		s.mockVault.EXPECT().CreateSecret(gomock.Any(), createInput).Return(&secretsmanager.CreateSecretOutput{}, awserr.New(secretsmanager.ErrCodeResourceExistsException, "", nil))
		s.mockVault.EXPECT().PutSecretValue(gomock.Any(), putSecretInput).Return(&secretsmanager.PutSecretValueOutput{}, nil)
		s.mockVault.EXPECT().TagSecretResource(gomock.Any(), tagInput).Return(&secretsmanager.TagResourceOutput{}, nil)
		s.mockVault.EXPECT().DescribeSecret(gomock.Any(), descSecretInput).Return(descSecretOutput, nil)

		secret, err := s.secretStore.Set(ctx, id, value, attributes)

		assert.NoError(t, err)
		assert.Equal(t, value, secret.Value)

		assert.ObjectsAreEqual(attributes.Tags, secret.Tags)
		assert.Equal(t, version, secret.Metadata.Version)
	})

}

func (s *awsSecretStoreTestSuite) TestGet() {
	ctx := context.Background()
	id := "my-secret"
	version := "some-version"
	secretValue := "secret-value"

	expectedSecret := &entities.Secret{
		ID:    id,
		Value: secretValue,
	}

	getSecretInput := &secretsmanager.GetSecretValueInput{
		SecretId: &id,
	}

	getSecretOutput := &secretsmanager.GetSecretValueOutput{
		Name:         &id,
		SecretString: &secretValue,
		VersionId:    &version,
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
	}

	s.T().Run("should get a secret successfully", func(t *testing.T) {
		s.mockVault.EXPECT().GetSecret(gomock.Any(), getSecretInput).Return(getSecretOutput, nil)
		s.mockVault.EXPECT().DescribeSecret(gomock.Any(), descSecretInput).Return(descSecretOutput, nil)
		retValue, err := s.secretStore.Get(ctx, id, "")
		assert.NoError(t, err)
		assert.Equal(t, retValue.Value, expectedSecret.Value)
		assert.Equal(t, retValue.ID, expectedSecret.ID)
	})
}

func (s *awsSecretStoreTestSuite) TestList() {
	ctx := context.Background()
	sec3, sec4 := "my-secret3", "my-secret4"
	expected := []string{sec3, sec4}
	secretsList := []*secretsmanager.SecretListEntry{{Name: &sec3}, {Name: &sec4}}

	s.T().Run("should list all secret ids successfully", func(t *testing.T) {

		listInput := &secretsmanager.ListSecretsInput{}
		listOutput := &secretsmanager.ListSecretsOutput{
			SecretList: secretsList,
		}

		s.mockVault.EXPECT().ListSecrets(gomock.Any(), listInput).Return(listOutput, nil)
		ids, err := s.secretStore.List(ctx)

		assert.NoError(t, err)
		assert.Equal(t, expected, ids)
	})

	s.T().Run("should return empty list if result is nil", func(t *testing.T) {
		s.mockVault.EXPECT().ListSecrets(gomock.Any(), gomock.Any()).Return(&secretsmanager.ListSecretsOutput{}, nil)
		ids, err := s.secretStore.List(ctx)

		assert.NoError(t, err)
		assert.Empty(t, ids)
	})

	s.T().Run("should fail if list fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("error")

		s.mockVault.EXPECT().ListSecrets(gomock.Any(), gomock.Any()).Return(&secretsmanager.ListSecretsOutput{}, expectedErr)
		ids, err := s.secretStore.List(ctx)

		assert.Nil(t, ids)
		assert.Equal(t, expectedErr, err)
	})
}
