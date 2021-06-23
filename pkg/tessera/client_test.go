package tessera

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/consensysquorum/quorum-key-manager/pkg/http/testutils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreRaw(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	transport := testutils.NewMockRoundTripper(ctrl)
	client := NewHTTPClient(&http.Client{Transport: transport})

	tests := []struct {
		desc string

		// JSON body of the request
		payload     []byte
		privateFrom string

		prepare func()

		expectedBody []byte

		respBody    []byte
		expectedKey []byte
	}{
		{
			desc:         "Base case",
			payload:      []byte{0xab, 0xcd},
			privateFrom:  "KkOjNLmCI6r+mICrC6l+XuEDjFEzQllaMQMpWLl4y1s=",
			expectedBody: []byte(`{"payload":"q80=","privateFrom":"KkOjNLmCI6r+mICrC6l+XuEDjFEzQllaMQMpWLl4y1s="}`),
			respBody:     []byte(`{"key":"SGVsbG8sIOS4lueVjA=="}`),
			expectedKey:  []byte{0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x2c, 0x20, 0xe4, 0xb8, 0x96, 0xe7, 0x95, 0x8c},
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			m := testutils.RequestMatcher(
				t,
				"/storeraw",
				tt.expectedBody,
			)

			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader(tt.respBody)),
				Header:     make(http.Header),
			}
			resp.Header.Set("Content-Type", "application/json")
			transport.EXPECT().RoundTrip(m).Return(resp, nil)

			enclaveKey, err := client.StoreRaw(context.Background(), tt.payload, tt.privateFrom)
			require.NoError(t, err, "StoreRaw must not error")
			assert.Equal(t, tt.expectedKey, enclaveKey, "Key should be valid")
		})
	}
}
