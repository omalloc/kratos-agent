package main

import (
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
	"github.com/omalloc/contrib/kratos/zap"
	"go.uber.org/automaxprocs/maxprocs"

	"github.com/omalloc/kratos-agent/internal/conf"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// To render a whole-file example, we need a package-level declaration.
	_ = ""
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// GitHash is the git-hash of the compiled software.
	GitHash string
	// Built is build-time the compiled software.
	Built string
	// flagconf is the config flag.
	flagconf    string
	flagverbose bool

	id, _ = os.Hostname()
)

func init() {
	maxprocs.Set(maxprocs.Logger(nil))

	rootCmd.PersistentFlags().StringVar(&flagconf, "conf", "../../configs", "config path")
	rootCmd.PersistentFlags().BoolVarP(&flagverbose, "verbose", "v", false, "verbose output")

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
	)
}

func main() {
	rootCmd.Execute()
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
