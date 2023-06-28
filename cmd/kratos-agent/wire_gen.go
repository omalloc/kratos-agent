// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/omalloc/contrib/kratos/health"
	"github.com/omalloc/contrib/kratos/registry"
	"github.com/omalloc/kratos-agent/internal/biz"
	"github.com/omalloc/kratos-agent/internal/conf"
	"github.com/omalloc/kratos-agent/internal/data"
	"github.com/omalloc/kratos-agent/internal/server"
	"github.com/omalloc/kratos-agent/internal/service"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp(bootstrap *conf.Bootstrap, confServer *conf.Server, confData *conf.Data, logger log.Logger) (*kratos.App, func(), error) {
	protobufRegistry := server.NewRegistryConfig(bootstrap)
	client, cleanup, err := registry.NewEtcd(protobufRegistry)
	if err != nil {
		return nil, nil, err
	}
	registrar := registry.NewRegistrar(client)
	criUsecase := biz.NewCRIUsecase(logger)
	agentService := service.NewAgentService(logger, criUsecase)
	grpcServer := server.NewGRPCServer(confServer, agentService, logger)
	httpServer := server.NewHTTPServer(confServer, agentService, logger)
	dataData, cleanup2, err := data.NewData(confData, logger)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	v := server.NewChecker(dataData, agentService)
	healthServer := health.NewServer(v, logger, httpServer)
	app := newApp(logger, registrar, grpcServer, httpServer, healthServer)
	return app, func() {
		cleanup2()
		cleanup()
	}, nil
}
