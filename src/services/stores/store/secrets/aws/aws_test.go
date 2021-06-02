package aws

import (
	"context"
	"fmt"
	"testing"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/log"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/services/stores/infra/aws/mocks"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/services/stores/store/entities"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/services/stores/store/entities/testutils"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/services/stores/store/secrets"
	"github.com/aws/aws-sdk-go/aws/awserr"
	_ "github.com/ConsenSysQuorum/quorum-key-manager/src/infra/aws"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/infra/aws/mocks"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities"
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
	version := "50.0.1"
	value := "my-value1"
	attributes := testutils.FakeAttributes()

	createOutput := &secretsmanager.CreateSecretOutput{
		Name:      &id,
		VersionId: &version,
	}

	metadata := &entities.Metadata{
		Version: version,
	}

	s.Run("should set a new secret successfully", func() {
		s.mockVault.EXPECT().CreateSecret(gomock.Any(), id, value).Return(createOutput, nil)
		s.mockVault.EXPECT().TagSecretResource(gomock.Any(), id, attributes.Tags).Return(&secretsmanager.TagResourceOutput{}, nil)
		s.mockVault.EXPECT().DescribeSecret(gomock.Any(), id).Return(attributes.Tags, metadata, nil)

		secret, err := s.secretStore.Set(ctx, id, value, attributes)

		assert.NoError(s.T(), err)
		assert.Equal(s.T(), value, secret.Value)

		assert.ObjectsAreEqual(attributes.Tags, secret.Tags)
		assert.Equal(s.T(), version, secret.Metadata.Version)
		assert.False(s.T(), secret.Metadata.Disabled)
		assert.True(s.T(), secret.Metadata.ExpireAt.IsZero())
		assert.True(s.T(), secret.Metadata.DeletedAt.IsZero())
	})

	s.Run("should fail when too many tags", func() {
		tooManyTags := map[string]string{}

		for i := 0; i <= maxTagsAllowed; i++ {
			tooManyTags[fmt.Sprintf("tag%d", i)] = fmt.Sprintf("value%d", i)
		}
		attributes.Tags = tooManyTags
		s.mockVault.EXPECT().CreateSecret(gomock.Any(), id, value).Return(createOutput, nil)
		secret, err := s.secretStore.Set(ctx, id, value, attributes)

		// tags back to normal
		attributes.Tags = testutils.FakeTags()

		assert.NotNil(s.T(), err)
		assert.True(s.T(), errors.IsInvalidParameterError(err))
		assert.Nil(s.T(), secret)
	})

	s.Run("should fail with describe error", func() {
		expectedErr := fmt.Errorf("any error")
		s.mockVault.EXPECT().CreateSecret(gomock.Any(), id, value).Return(createOutput, nil)
		s.mockVault.EXPECT().TagSecretResource(gomock.Any(), id, attributes.Tags).Return(&secretsmanager.TagResourceOutput{}, nil)
		s.mockVault.EXPECT().DescribeSecret(gomock.Any(), id).Return(testutils.FakeTags(), testutils.FakeMetadata(), expectedErr)

		secret, err := s.secretStore.Set(ctx, id, value, attributes)

		assert.Equal(s.T(), err, expectedErr)
		assert.Nil(s.T(), secret)

	})

	s.Run("should fail with tag error", func() {
		expectedErr := fmt.Errorf("any error")
		s.mockVault.EXPECT().CreateSecret(gomock.Any(), id, value).Return(createOutput, nil)
		s.mockVault.EXPECT().TagSecretResource(gomock.Any(), id, attributes.Tags).Return(nil, expectedErr)

		secret, err := s.secretStore.Set(ctx, id, value, attributes)

		assert.Equal(s.T(), err, expectedErr)
		assert.Nil(s.T(), secret)

	})

	s.Run("should fail with same error if write fails", func() {
		expectedErr := fmt.Errorf("error")
		s.mockVault.EXPECT().CreateSecret(gomock.Any(), id, value).Return(&secretsmanager.CreateSecretOutput{}, expectedErr)
		s.mockVault.EXPECT().TagSecretResource(gomock.Any(), id, attributes.Tags).Return(&secretsmanager.TagResourceOutput{}, nil)
		s.mockVault.EXPECT().DescribeSecret(gomock.Any(), id).Return(testutils.FakeTags(), testutils.FakeMetadata(), nil)

		secret, err := s.secretStore.Set(ctx, id, value, attributes)

		assert.Nil(s.T(), secret)
		assert.Equal(s.T(), expectedErr, err)
	})

	s.Run("should update secret if already exists", func() {

		s.mockVault.EXPECT().CreateSecret(gomock.Any(), id, value).Return(&secretsmanager.CreateSecretOutput{}, errors.AlreadyExistsError("already exists"))
		s.mockVault.EXPECT().PutSecretValue(gomock.Any(), id, value).Return(&secretsmanager.PutSecretValueOutput{}, nil)
		s.mockVault.EXPECT().TagSecretResource(gomock.Any(), id, attributes.Tags).Return(&secretsmanager.TagResourceOutput{}, nil)
		s.mockVault.EXPECT().DescribeSecret(gomock.Any(), id).Return(testutils.FakeTags(), testutils.FakeMetadata(), nil)

		secret, err := s.secretStore.Set(ctx, id, value, attributes)

		assert.NoError(s.T(), err)
		assert.Equal(s.T(), value, secret.Value)

		assert.ObjectsAreEqual(attributes.Tags, secret.Tags)
		assert.Equal(s.T(), version, secret.Metadata.Version)
	})

}

func (s *awsSecretStoreTestSuite) TestGet() {
	ctx := context.Background()
	id := "my-secret-get"
	version := "some-version"
	secretValue := "secret-value"

	expectedSecret := &entities.Secret{
		ID:    id,
		Value: secretValue,
	}

	getSecretOutput := &secretsmanager.GetSecretValueOutput{
		Name:         &id,
		SecretString: &secretValue,
		VersionId:    &version,
	}

	currentMark := CurrentVersionMark

	s.T().Run("should get a secret successfully", func(t *testing.T) {
		s.mockVault.EXPECT().GetSecret(gomock.Any(), id, "").Return(getSecretOutput, nil)
		s.mockVault.EXPECT().DescribeSecret(gomock.Any(), id).Return(testutils.FakeTags(), testutils.FakeMetadata(), nil)
		retValue, err := s.secretStore.Get(ctx, id, "")
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), retValue.Value, expectedSecret.Value)
		assert.Equal(s.T(), retValue.ID, expectedSecret.ID)
	})

	s.Run("should fail with get error", func() {
		expectedErr := errors.NotFoundError("secret not found")
		s.mockVault.EXPECT().GetSecret(gomock.Any(), id, version).Return(getSecretOutput, expectedErr)

		retValue, err := s.secretStore.Get(ctx, id, version)
		assert.Nil(s.T(), retValue)
		assert.Equal(s.T(), err, expectedErr)
	})

	s.Run("should fail with describe error", func() {
		expectedErr := errors.NotFoundError("secret not found")
		s.mockVault.EXPECT().GetSecret(gomock.Any(), id, version).Return(getSecretOutput, nil)
		s.mockVault.EXPECT().DescribeSecret(gomock.Any(), id).Return(testutils.FakeTags(), testutils.FakeMetadata(), expectedErr)
		retValue, err := s.secretStore.Get(ctx, id, version)
		assert.Nil(s.T(), retValue)
		assert.Equal(s.T(), err, expectedErr)
	})
}

func (s *awsSecretStoreTestSuite) TestGetDeleted() {
	s.Run("should fail with not implemented error", func() {
		ctx := context.Background()
		id := "some-id"
		expectedError := errors.ErrNotImplemented
		_, err := s.secretStore.GetDeleted(ctx, id)

		assert.Equal(s.T(), err, expectedError)
	})
}

func (s *awsSecretStoreTestSuite) TestDeleted() {
	ctx := context.Background()
	id := "my-secret-deleted"
	destroy := false

	deleteOutput := &secretsmanager.DeleteSecretOutput{
		Name: &id,
	}

	s.Run("should delete secret successfully", func() {
		s.mockVault.EXPECT().DeleteSecret(gomock.Any(), id, destroy).Return(deleteOutput, nil)

		err := s.secretStore.Delete(ctx, id)
		assert.NoError(s.T(), err)
	})
}

func (s *awsSecretStoreTestSuite) TestDestroy() {
	ctx := context.Background()
	id := "my-secret"
	destroy := true

	deleteOutput := &secretsmanager.DeleteSecretOutput{
		Name: &id,
	}

	s.Run("should destroy secret successfully", func() {
		s.mockVault.EXPECT().DeleteSecret(gomock.Any(), id, destroy).Return(deleteOutput, nil)

		err := s.secretStore.Destroy(ctx, id)

		assert.NoError(s.T(), err)
	})

	s.Run("should fail to destroy secret with internal error", func() {

		expectedError := errors.InternalError("internal error")
		s.mockVault.EXPECT().DeleteSecret(gomock.Any(), id, destroy).Return(nil, expectedError)

		err := s.secretStore.Destroy(ctx, id)

		assert.Error(s.T(), err)
		assert.Equal(s.T(), err, expectedError)
	})

	s.Run("should fail to destroy secret with internal error", func() {
		expectedError := errors.InternalError("internal error")
		s.mockVault.EXPECT().DeleteSecret(gomock.Any(), id, destroy).Return(nil, expectedError)

		err := s.secretStore.Destroy(ctx, id)

		assert.Error(s.T(), err)
		assert.Equal(s.T(), err, expectedError)
	})

	s.Run("should fail to destroy secret with not found error", func() {

		expectedError := errors.NotFoundError("resource was not found")
		s.mockVault.EXPECT().DeleteSecret(gomock.Any(), id, destroy).Return(nil, expectedError)

		err := s.secretStore.Destroy(ctx, id)

		assert.Error(s.T(), err)
		assert.Equal(s.T(), err, expectedError)
	})

	s.Run("should fail to destroy secret with not found error", func() {

		expectedError := errors.InternalError("internal error, limit exceeded")
		s.mockVault.EXPECT().DeleteSecret(gomock.Any(), id, destroy).Return(nil, expectedError)

		err := s.secretStore.Destroy(ctx, id)

		assert.Error(s.T(), err)
		assert.Equal(s.T(), err, expectedError)
	})

	s.Run("should fail to destroy secret with invalid parameter error", func() {

		expectedError := errors.InvalidParameterError("invalid parameter")
		s.mockVault.EXPECT().DeleteSecret(gomock.Any(), id, destroy).Return(nil, expectedError)

		err := s.secretStore.Destroy(ctx, id)

		assert.Error(s.T(), err)
		assert.Equal(s.T(), err, expectedError)
	})

	s.Run("should fail to destroy secret with invalid request error", func() {

		expectedError := errors.InvalidRequestError("invalid request")
		s.mockVault.EXPECT().DeleteSecret(gomock.Any(), id, destroy).Return(nil, expectedError)

		err := s.secretStore.Destroy(ctx, id)

		assert.Error(s.T(), err)
		assert.Equal(s.T(), err, expectedError)
	})
}

func (s *awsSecretStoreTestSuite) TestList() {
	ctx := context.Background()
	sec3, sec4 := "my-secret3", "my-secret4"
	expected := []string{sec3, sec4}
	secretsList := []*secretsmanager.SecretListEntry{{Name: &sec3}, {Name: &sec4}}

	s.Run("should list all secret ids successfully", func() {

		listOutput := &secretsmanager.ListSecretsOutput{
			SecretList: secretsList,
		}

		s.mockVault.EXPECT().ListSecrets(gomock.Any(), int64(0), "").Return(listOutput, nil)
		ids, err := s.secretStore.List(ctx)

		assert.NoError(s.T(), err)
		assert.Equal(s.T(), expected, ids)
	})

	s.Run("should list all secret ids successfully with a nextToken", func() {

		nextToken := "next"
		listOutput := &secretsmanager.ListSecretsOutput{
			SecretList: secretsList,
			NextToken:  &nextToken,
		}

		s.mockVault.EXPECT().ListSecrets(gomock.Any(), int64(0), "").Return(listOutput, nil)
		listOutput.NextToken = nil
		s.mockVault.EXPECT().ListSecrets(gomock.Any(), int64(0), nextToken).Return(listOutput, nil)
		ids, err := s.secretStore.List(ctx)

		assert.NoError(s.T(), err)
		assert.Equal(s.T(), expected, ids)
	})

	s.Run("should return empty list if result is nil", func() {
		s.mockVault.EXPECT().ListSecrets(gomock.Any(), int64(0), "").Return(&secretsmanager.ListSecretsOutput{}, nil)
		ids, err := s.secretStore.List(ctx)

		assert.NoError(s.T(), err)
		assert.Empty(s.T(), ids)
	})

	s.Run("should fail if list fails", func() {
		expectedErr := fmt.Errorf("error")

		s.mockVault.EXPECT().ListSecrets(gomock.Any(), int64(0), "").Return(&secretsmanager.ListSecretsOutput{}, expectedErr)
		ids, err := s.secretStore.List(ctx)

		assert.Nil(s.T(), ids)
		assert.Equal(s.T(), expectedErr, err)
	})
}

func (s *awsSecretStoreTestSuite) TestListDeleted() {
	s.Run("should fail with not implemented error", func() {
		ctx := context.Background()
		expectedError := errors.ErrNotImplemented
		_, err := s.secretStore.ListDeleted(ctx)

		assert.Equal(s.T(), err, expectedError)
	})
}

func ToSecretmanagerTags(tags map[string]string) []*secretsmanager.Tag {
	var fakeSecretsTags []*secretsmanager.Tag

	for key, value := range tags {
		k, v := key, value
		var in = secretsmanager.Tag{
			Key:   &k,
			Value: &v,
		}
		fakeSecretsTags = append(fakeSecretsTags, &in)
	}
	return fakeSecretsTags
}
