package request

import (
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCombinePreparer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	preparer1 := NewMockPreparer(ctrl)
	preparer2 := NewMockPreparer(ctrl)

	preparer := CombinePreparer(preparer1, preparer2)

	req := &http.Request{}
	call1 := preparer1.EXPECT().Prepare(req).Return(req, nil)
	preparer2.EXPECT().Prepare(req).Return(req, nil).After(call1)

	outReq, err := preparer.Prepare(req)
	assert.NoError(t, err, "Prepare should not error")
	assert.Equal(t, req, outReq, "Return req should be correct")
}
