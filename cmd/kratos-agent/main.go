package main

import (
	"context"
	"os"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/omalloc/contrib/kratos/health"
	trace "github.com/omalloc/contrib/kratos/tracing"
	"github.com/omalloc/contrib/kratos/zap"
	"go.uber.org/automaxprocs/maxprocs"

	"github.com/omalloc/kratos-agent/internal/conf"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// To render a whole-file example, we need a package-level declaration.
	_ = ""
	// Name is the name of the compiled software.
	Name string = "kratos-agent"
	// Version is the version of the compiled software.
	Version string = "v0.0.0"
	// GitHash is the git-hash of the compiled software.
	GitHash string = "no set hash"
	// Built is build-time the compiled software.
	Built string = "0"
	// flagconf is the config flag.
	flagconf string
	// flagverbose is the verbose flag.
	flagverbose bool
	// service unique id
	id string
)

func init() {
	_, _ = maxprocs.Set(maxprocs.Logger(nil))

	hostname, _ := os.Hostname()

	rootCmd.PersistentFlags().StringVar(&flagconf, "conf", "../../configs", "config path")
	rootCmd.PersistentFlags().BoolVarP(&flagverbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringVar(&id, "id", hostname, "service unique-id; default: hostname")

	rootCmd.AddCommand(versionCmd)
}

func newApp(logger log.Logger, registrar registry.Registrar, gs *grpc.Server, hs *http.Server, hh *health.Server) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Registrar(registrar),
		kratos.Server(
			gs,
			hs,
			hh,
		),
		kratos.AfterStart(func(ctx context.Context) error {
			log.Infof("started with pid %d", os.Getpid())
			return nil
		}),
	)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
	logger := log.With(zap.NewLogger(zap.Verbose(flagverbose)),
		"ts", log.DefaultTimestamp,
		// "caller", log.DefaultCaller,
		"service.id", id,
		"service.name", Name,
		"service.version", Version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)
	log.SetLogger(logger)

	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	trace.InitTracer(
		trace.WithEndpoint(bc.Tracing.Endpoint),
		trace.WithServiceName(Name),
	)

	app, cleanup, err := wireApp(&bc, bc.Server, bc.Data, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
