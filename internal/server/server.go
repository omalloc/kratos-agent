package server

import (
	"github.com/google/wire"
	"github.com/omalloc/contrib/kratos/health"
	"github.com/omalloc/contrib/kratos/registry"
	"github.com/omalloc/contrib/protobuf"

	"github.com/omalloc/kratos-agent/internal/conf"
	"github.com/omalloc/kratos-agent/internal/data"
	"github.com/omalloc/kratos-agent/internal/service"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(
	NewGRPCServer,
	NewHTTPServer,
	NewRegistryConfig,
	NewChecker,

	registry.NewEtcd,
	registry.NewRegistrar,
	registry.NewDiscovery,

	health.NewServer,
)

func NewRegistryConfig(bc *conf.Bootstrap) *protobuf.Registry {
	return bc.Registry
}

func NewTracingConfig(bc *conf.Bootstrap) *protobuf.Tracing {
	return bc.Tracing
}

func NewChecker(c1 *data.Data, agent *service.AgentService) []health.Checker {
	return []health.Checker{c1, agent}
}
