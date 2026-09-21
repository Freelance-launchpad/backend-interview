package main

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/DataDog/dd-trace-go/v2/ddtrace/ext"
	"github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"
	"github.com/Freelance-launchpad/backend-interview/common/clients/gotenberg"
	"github.com/Freelance-launchpad/backend-interview/common/japplication"
	"github.com/Freelance-launchpad/backend-interview/common/jdb"
	"github.com/Freelance-launchpad/backend-interview/common/jhttp/httpauth"
	"github.com/Freelance-launchpad/backend-interview/common/jhttp/v2"
	"github.com/Freelance-launchpad/backend-interview/services/activity-reports/internal/domain"
	"github.com/Freelance-launchpad/backend-interview/services/activity-reports/internal/infrastructure/generator"
	"github.com/Freelance-launchpad/backend-interview/services/activity-reports/internal/infrastructure/postgres"
	"github.com/Freelance-launchpad/backend-interview/services/activity-reports/internal/usecases"
	usersclient "github.com/Freelance-launchpad/backend-interview/services/users/pkg/client"

	transactor "github.com/Thiht/transactor/stdlib"
)

//go:embed resources/template_activity_report.html
var activityReportTemplate []byte

var appName = "activity-reports-cronjobs"

type dependencies struct {
	db                     *sql.DB
	activityReportsService domain.ActivityReportService
}

func initDependencies(app japplication.Application[configuration]) (dependencies, error) {
	db, err := jdb.OpenDatabase(app.Configuration.Database)
	if err != nil {
		return dependencies{}, fmt.Errorf("failed to open database: %w", err)
	}
	transactor, dbGetter := transactor.NewTransactor(db, transactor.NestedTransactionsSavepoints)

	activityReportsStore := postgres.NewActivityReportStore(dbGetter)

	m2mAuth := httpauth.NewAuth(
		app.Name, app.Version,
		app.Configuration.AuthClientSecret,
		*app.Configuration.AuthURL.URL,
	)

	usersClient := usersclient.New(
		jhttp.NewClient(10*time.Second).
			WithBearerAuthentication(m2mAuth).
			WithUserAgent(app.Name, app.Version).
			WithTracing("users"),
		*app.Configuration.UsersURL.URL,
	)
	gotenberg := gotenberg.New(
		jhttp.NewClient(60*time.Second).
			WithUserAgent(app.Name, app.Version).
			WithTracing("gotenberg"),
		*app.Configuration.GotenbergURL.URL,
	)

	generator, err := generator.NewPDFGenerator(activityReportTemplate, usersClient, gotenberg)
	if err != nil {
		return dependencies{}, fmt.Errorf("failed to initialize generator: %w", err)
	}
	activityReportsService := usecases.NewActivityReportService(
		activityReportsStore,
		transactor,
		generator,
		usersClient,
	)

	return dependencies{
		db:                     db,
		activityReportsService: activityReportsService,
	}, nil
}

func main() {
	app := japplication.New[configuration](appName)
	defer tracer.Stop()

	deps, err := initDependencies(app)
	if err != nil {
		slog.Error("Failed to initialize dependencies", slog.Any("error", err))
		os.Exit(1)
	}
	defer deps.db.Close()

	slog.Info(appName+" is running", slog.Any("args", os.Args))

	if len(os.Args) < 3 {
		slog.Error("command arg required")
		os.Exit(1)
	}

	command := os.Args[2]

	span, ctx := tracer.StartSpanFromContext(context.Background(), "cronjobs",
		tracer.Tag(ext.ServiceName, appName),
		tracer.Tag(ext.ResourceName, command),
		tracer.Tag(ext.SpanType, "cron"),
	)

	switch command {
	case "generate-reports-pdfs":
		// REDACTED
	default:
		err = errors.New("command not recognized")
	}

	if err != nil {
		slog.ErrorContext(ctx, appName+" finished with error",
			slog.Any("error", err),
			slog.Any("command", command),
		)
	}
	span.Finish(tracer.WithError(err))
}
