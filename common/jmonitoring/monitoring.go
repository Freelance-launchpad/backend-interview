package jmonitoring

import (
	"fmt"
	"log/slog"
	"regexp"
	"strings"
	"time"

	"github.com/DataDog/datadog-go/v5/statsd"
	"github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"

	"github.com/Freelance-launchpad/backend-interview/common/jenv"
)

// Configuration describes the configuration needed for the monitoring
type Configuration struct {
	Enabled bool `yaml:"enabled"`
}

// Note: LogLevel values are not exposed at the moment, see: https://github.com/DataDog/dd-trace-go/discussions/4822
const (
	tracerLevelDebug tracer.LogLevel = iota
	tracerLevelInfo
	tracerLevelWarn
	tracerLevelError
)

var defaultClient statsd.ClientInterface = new(statsd.NoOpClient)

// Init starts the statsd client if enabled and in the correct environment
func Init(appName, version string, enabled bool) error {
	if !enabled || jenv.InLocal() {
		return nil
	}

	var err error
	defaultClient, err = statsd.New("")
	if err != nil {
		return fmt.Errorf("statsd.New() error: %w", err)
	}

	environment := jenv.Name()
	if err := tracer.Start(
		tracer.WithLogStartup(false),
		tracer.WithEnv(environment),
		tracer.WithService(appName),
		tracer.WithServiceVersion(version),
		tracer.WithRuntimeMetrics(),
		tracer.WithSendRetries(2),
		tracer.WithLogger(tracer.AdaptLogger(logTracer)),
	); err != nil {
		return fmt.Errorf("tracer.Start() error: %w", err)
	}

	return nil
}

func Stop() {
	tracer.Stop()

	if err := defaultClient.Close(); err != nil {
		slog.Error("statsd.Close error", slog.Any("error", err))
	}
}

func logTracer(lvl tracer.LogLevel, msg string, _ ...any) {
	if lvl == tracerLevelError && isErrorDemotedToWarning(msg) {
		lvl = tracerLevelWarn
	}

	switch lvl {
	case tracerLevelDebug:
		slog.Debug(msg)
	case tracerLevelInfo:
		slog.Info(msg)
	case tracerLevelWarn:
		slog.Warn(msg)
	case tracerLevelError:
		slog.Error(msg)
	default:
		slog.Warn("tracer: unknown log level", slog.Any("level", lvl), slog.String("msg", msg))
	}
}

var reLostTraces = regexp.MustCompile(`lost \d+ traces`)

func isErrorDemotedToWarning(msg string) bool {
	return reLostTraces.MatchString(msg) ||
		strings.Contains(msg, "Error sending stats payload:")
}

// Count tracks how many times something happened per second.
func Count(name string, value int64, tags []string, rate float64) {
	err := defaultClient.Count(name, value, tags, rate)
	if err != nil {
		slog.Error("statsd.Count error", slog.Any("error", err))
	}
}

// Timing sends timing information, it is an alias for TimeInMilliseconds
func Timing(name string, value time.Duration, tags []string, rate float64) {
	err := defaultClient.Timing(name, value, tags, rate)
	if err != nil {
		slog.Error("statsd.Timing error", slog.Any("error", err))
	}
}

// Histogram tracks the statistical distribution of a set of values on each host.
func Histogram(name string, value float64, tags []string, rate float64) {
	err := defaultClient.Histogram(name, value, tags, rate)
	if err != nil {
		slog.Error("statsd.Histogram error", slog.Any("error", err))
	}
}

// Gauge measures the value of a metric at a particular time.
func Gauge(name string, value float64, tags []string, rate float64) {
	err := defaultClient.Gauge(name, value, tags, rate)
	if err != nil {
		slog.Error("statsd.Gauge error", slog.Any("error", err))
	}
}
