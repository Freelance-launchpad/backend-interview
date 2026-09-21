package main

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/Freelance-launchpad/backend-interview/common/japplication"
	"github.com/Freelance-launchpad/backend-interview/common/jdb"
	"github.com/Freelance-launchpad/backend-interview/common/jgin"
	"github.com/Freelance-launchpad/backend-interview/common/jhttp/httpauth"
	"github.com/Freelance-launchpad/backend-interview/common/jhttp/v2"
	"github.com/Freelance-launchpad/backend-interview/services/activity-reports/internal/domain"
	adminhandlers "github.com/Freelance-launchpad/backend-interview/services/activity-reports/internal/handlers/admin"
	userhandlers "github.com/Freelance-launchpad/backend-interview/services/activity-reports/internal/handlers/user"
	"github.com/Freelance-launchpad/backend-interview/services/activity-reports/internal/infrastructure/postgres"
	"github.com/Freelance-launchpad/backend-interview/services/activity-reports/internal/usecases"
	usersclient "github.com/Freelance-launchpad/backend-interview/services/users/pkg/client"
	transactor "github.com/Thiht/transactor/stdlib"
)

// These variables are overridden at compile time
var (
	appName   = "activity-reports"
	buildDate = "unknown"
)

type dependencies struct {
	authenticator              *jgin.Auth
	activityReportService      domain.ActivityReportService
	activityReportStore        domain.ActivityReportStore
	usersClient                usersclient.Client
	userActivityReportHandler  *userhandlers.ActivityReportHandler
	adminActivityReportHandler *adminhandlers.ActivityReportHandler
}

func initDependencies(app japplication.Application[configuration]) (dependencies, error) {
	legacyAuthenticator, err := jgin.NewLegacyAuth(&app.Configuration.Authorization)
	if err != nil {
		return dependencies{}, fmt.Errorf("failed to initialize authenticator: %w", err)
	}
	authenticator := jgin.NewAuth(app.Configuration.JWTPublicKey).
		WithBannedJTIs(app.Configuration.BannedJTIs).
		WithFallback(legacyAuthenticator)

	db, err := jdb.OpenDatabase(app.Configuration.Database)
	if err != nil {
		return dependencies{}, err
	}
	transactor, dbGetter := transactor.NewTransactor(db, transactor.NestedTransactionsSavepoints)

	activityReportStore := postgres.NewActivityReportStore(dbGetter)

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

	activityReportService := usecases.NewActivityReportService(activityReportStore, transactor, nil, usersClient)

	userActivityReportHandler := userhandlers.NewActivityReportHandler(activityReportStore, activityReportService)
	adminActivityReportHandler := adminhandlers.NewActivityReportHandler(activityReportStore)

	return dependencies{
		authenticator:              authenticator,
		activityReportService:      activityReportService,
		activityReportStore:        activityReportStore,
		usersClient:                usersClient,
		userActivityReportHandler:  userActivityReportHandler,
		adminActivityReportHandler: adminActivityReportHandler,
	}, nil
}

func initApp(app japplication.Application[configuration], deps dependencies) *jgin.Engine {
	router := jgin.New(app.Address, app.Name, app.Version, buildDate)

	{
		v3 := router.Group("/activity-reports/v3", deps.authenticator.Authenticate())
		{
			activityReports := v3.Group("/activity-reports", jgin.IsLifeOffer())
			activityReports.POST("", deps.userActivityReportHandler.CreateActivityReports)
			activityReports.GET("/missing", deps.userActivityReportHandler.GetMissingActivityReports)
		}

		{
			admin := v3.Group("/admin", jgin.RequirePermissions("activity-report.all.read"))
			admin.GET("/offers/:offer_id/activity-reports", deps.adminActivityReportHandler.GetActivityReports)
		}
	}

	return router
}

// @title activity-reports
// @version 1.0.0
// @description This api exposes activity-reports endpoints.
// @BasePath /activity-reports
func main() {
	app := japplication.New[configuration](appName)

	deps, err := initDependencies(app)
	if err != nil {
		slog.Error("failed to init dependencies", slog.Any("error", err))
		os.Exit(1)
	}

	router := initApp(app, deps)
	app.AddServices(router)
	app.Run()
}
