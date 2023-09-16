package biz

import (
	"github.com/google/wire"
	"go.opentelemetry.io/otel"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(
	NewCRIUsecase,
	NewClusterUsecase,
	NewCrontabUsecase,
)

var tracer = otel.Tracer("kratos-agent")
