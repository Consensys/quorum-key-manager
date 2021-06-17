package app

import (
	"context"
	"testing"

	"github.com/consensysquorum/quorum-key-manager/pkg/log/mock"
	"github.com/golang/mock/gomock"

	"github.com/consensysquorum/quorum-key-manager/pkg/common"
	"github.com/consensysquorum/quorum-key-manager/pkg/http/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type ServiceInterface1 interface {
	common.Runnable
	Method1()
}

type ServiceStruct1 struct {
	name string
}

func (srv *ServiceStruct1) Start(context.Context) error { return nil }
func (srv *ServiceStruct1) Stop(context.Context) error  { return nil }
func (srv *ServiceStruct1) Error() error                { return nil }
func (srv *ServiceStruct1) Close() error                { return nil }
func (srv *ServiceStruct1) Method1()                    {}

type ServiceInterface2 interface {
	common.Runnable
	Method2()
}

type ServiceStruct2 struct {
	name string
}

func (srv ServiceStruct2) Start(context.Context) error { return nil }
func (srv ServiceStruct2) Stop(context.Context) error  { return nil }
func (srv ServiceStruct2) Error() error                { return nil }
func (srv ServiceStruct2) Close() error                { return nil }
func (srv ServiceStruct2) Method2()                    {}

type TestConfig struct {
	name string
}

func TestRegisterServiceConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := New(&Config{HTTP: &server.Config{}}, mock.NewMockLogger(ctrl))

	cfg := &TestConfig{name: "test"}
	err := app.RegisterServiceConfig(cfg)
	require.NoError(t, err, "RegisterServiceConfig Ptr config must not error")

	cfg2 := new(TestConfig)
	err = app.ServiceConfig(cfg2)
	require.NoError(t, err, "ServiceConfig must not error")
	assert.Equal(t, cfg, cfg2, "ServiceConfig should match")

	err = app.RegisterServiceConfig(TestConfig{name: "test"})
	assert.Error(t, err, "RegisterServiceConfig not pointer must error")

	err = app.RegisterServiceConfig(common.ToPtr("test").(*string))
	assert.Error(t, err, "RegisterServiceConfig pointer to not struct must error")

	err = app.RegisterServiceConfig(nil)
	assert.Error(t, err, "RegisterServiceConfig nil must error")

	err = app.RegisterServiceConfig(new(ServiceInterface1))
	assert.Error(t, err, "RegisterServiceConfig pointer to interface must error")
}

func TestRegisterService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := New(&Config{HTTP: &server.Config{}}, mock.NewMockLogger(ctrl))

	srv1 := &ServiceStruct1{name: "test-srv1"}
	err := app.RegisterService(srv1)
	require.NoError(t, err, "RegisterService srv1 must not error")

	srv2 := ServiceStruct2{name: "test-srv2"}
	err = app.RegisterService(srv2)
	require.NoError(t, err, "RegisterService srv2 must not error")

	srvPtr1 := new(ServiceStruct1)
	err = app.Service(srvPtr1)
	require.NoError(t, err, "Service srvPtr1 must not error")
	assert.Equal(t, srv1, srvPtr1, "Service srvPtr1 should match")

	srvInterface1 := new(ServiceInterface1)
	err = app.Service(srvInterface1)
	require.NoError(t, err, "Service srvInterface1 must not error")
	assert.Equal(t, srv1, *srvInterface1, "Service srvInterface1 should match")

	srvPtr2 := new(ServiceInterface2)
	err = app.Service(srvPtr2)
	require.NoError(t, err, "Service srvPtr2 must not error")
	assert.Equal(t, srv2, *srvPtr2, "Service srvPtr2 should match")
}
