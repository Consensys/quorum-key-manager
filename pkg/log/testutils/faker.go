package testutils

import (
	"github.com/consensys/quorum-key-manager/pkg/log/mock"
	"github.com/golang/mock/gomock"
)

func NewMockLogger(ctrl *gomock.Controller) *mock.MockLogger {
	mockLogger := mock.NewMockLogger(ctrl)

	mockLogger.EXPECT().Debug(gomock.Any(), gomock.Any()).Return(mockLogger).AnyTimes()
	mockLogger.EXPECT().Info(gomock.Any(), gomock.Any()).Return(mockLogger).AnyTimes()
	mockLogger.EXPECT().Warn(gomock.Any(), gomock.Any()).Return(mockLogger).AnyTimes()
	mockLogger.EXPECT().Error(gomock.Any(), gomock.Any()).Return(mockLogger).AnyTimes()
	mockLogger.EXPECT().Panic(gomock.Any(), gomock.Any()).Return(mockLogger).AnyTimes()
	mockLogger.EXPECT().Fatal(gomock.Any(), gomock.Any()).Return(mockLogger).AnyTimes()
	mockLogger.EXPECT().With(gomock.Any()).Return(mockLogger).AnyTimes()
	mockLogger.EXPECT().WithError(gomock.Any()).Return(mockLogger).AnyTimes()
	mockLogger.EXPECT().WithComponent(gomock.Any()).Return(mockLogger).AnyTimes()
	mockLogger.EXPECT().Write(gomock.Any()).Return(0, nil).AnyTimes()

	return mockLogger
}
