package utils

import (
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.1/keyvault"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/consensys/quorum-key-manager/pkg/common"
	"github.com/consensys/quorum-key-manager/src/stores/store/entities/testutils"
	"github.com/stretchr/testify/assert"
)

func TestParseSecretBundle(t *testing.T) {
	id := "my-secret1"
	value := "my-value1"
	version := "2"
	secretBundleID := id + "/" + version
	attributes := testutils.FakeAttributes()

	expectedCreatedAt, _ := time.Parse(time.RFC3339, "2018-03-22T02:24:06.945319214Z")
	expectedUpdatedAt, _ := time.Parse(time.RFC3339, "2018-03-22T02:24:06.945319214Z")
	
	keyBundle := &keyvault.SecretBundle{
		Value: &value,
		ID:    &secretBundleID,
		Attributes: &keyvault.SecretAttributes{
			Created: &(&struct{ x date.UnixTime }{date.NewUnixTimeFromNanoseconds(expectedCreatedAt.UnixNano())}).x,
			Updated: &(&struct{ x date.UnixTime }{date.NewUnixTimeFromNanoseconds(expectedUpdatedAt.UnixNano())}).x,
			Enabled: &(&struct{ x bool }{true}).x,
		},
		Tags: common.Tomapstrptr(attributes.Tags),
	}

	secret := ParseSecretBundle(keyBundle)
	assert.Equal(t, id, secret.ID)
	assert.Equal(t, value, secret.Value)
	assert.Equal(t, expectedCreatedAt, secret.Metadata.CreatedAt)
	assert.Equal(t, attributes.Tags, secret.Tags)
	assert.Equal(t, version, secret.Metadata.Version)
	assert.False(t, secret.Metadata.Disabled)
	assert.True(t, secret.Metadata.ExpireAt.IsZero())
	assert.True(t, secret.Metadata.DeletedAt.IsZero())
}
