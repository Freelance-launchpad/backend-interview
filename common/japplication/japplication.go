package japplication

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
	_ "time/tzdata"

	"github.com/DataDog/dd-trace-go/v2/profiler"
	"github.com/Freelance-launchpad/backend-interview/common/confiture"
	"github.com/Freelance-launchpad/backend-interview/common/jenv"
	jlogv2 "github.com/Freelance-launchpad/backend-interview/common/jlog/v2"
	"github.com/Freelance-launchpad/backend-interview/common/jmonitoring"
)

type Service interface {
	Start() error
	Stop(ctx context.Context) error
}

type Application[T any] struct {
	Name          string
	Version       string
	Address       string
	Configuration T

	services []Service
}

const defaultConfigurationPath = "/etc/configuration.yaml"
const usage = "path of the service configuration file"

var configFilePath = flag.String("config", defaultConfigurationPath, usage)

func New[T any](appName string) Application[T] {

	version := os.Getenv("APP_VERSION")
	address := os.Getenv("APP_ADDR")
	monitoring, _ := strconv.ParseBool(os.Getenv("MONITORING_ENABLED"))
	profiling, _ := strconv.ParseBool(os.Getenv("PROFILING_ENABLED"))

	time.Local = time.UTC

	if jenv.InProd() && profiling {
		err := profiler.Start(
			profiler.WithLogStartup(false),
			profiler.WithEnv(jenv.Name()),
			profiler.WithService(appName),
			profiler.WithVersion(version),
			profiler.WithProfileTypes(
				profiler.CPUProfile,
				profiler.HeapProfile,
				profiler.GoroutineProfile,
				profiler.BlockProfile,
				profiler.MutexProfile,
			),
		)
		if err != nil {
			slog.Error("failed to init profiler", slog.Any("error", err))
			os.Exit(1)
		}
	}

	flag.Parse()
	jlogv2.Init(appName, version)

	var configuration T
	if err := confiture.Load(&configuration, *configFilePath); err != nil {
		slog.Error("failed to load configuration", slog.Any("error", err))
		os.Exit(1)
	}

	if err := jmonitoring.Init(appName, version, monitoring); err != nil {
		slog.Error("failed to init monitoring", slog.Any("error", err))
		os.Exit(1)
	}

	return Application[T]{
		Name:          appName,
		Version:       version,
		Address:       address,
		Configuration: configuration,
	}
}

func (app *Application[T]) AddServices(s ...Service) {
	app.services = append(app.services, s...)
}

func (app *Application[T]) start() {
	for _, service := range app.services {
		go func(s Service) {
			if err := s.Start(); err != nil {
				slog.Error("service.Run error", slog.Any("error", err))
				os.Exit(1)
			}
		}(service)
	}
}

func (app *Application[T]) stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, service := range app.services {
		if err := service.Stop(ctx); err != nil {
			slog.Error("service.Shutdown error", slog.Any("error", err))
			os.Exit(1)
		}
	}

	if jenv.InProd() {
		profiler.Stop()
	}

	jmonitoring.Stop()
}

func (app *Application[T]) Run() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	app.start()
	slog.Info(app.Name + " is running")
	<-ctx.Done()

	stop()
	slog.Info("shutting down gracefully, press Ctrl+C again to force")

	app.stop()
	slog.Info(app.Name + " exited")
}
