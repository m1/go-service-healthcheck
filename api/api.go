package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"

	"github.com/m1/go-service-healthcheck/config"
	"github.com/m1/go-service-healthcheck/log"
	"github.com/m1/go-service-healthcheck/repositories"
	"github.com/m1/go-service-healthcheck/storage"
)

type API struct {
	Config       *config.Config
	DB           *storage.SQLiteDB
	Repositories *Repositories
	Logger       *log.Logger

	router *chi.Mux
}

type Repositories struct {
	Service repositories.ServiceRepository
}

// New ...
func New() *API {
	r := chi.NewRouter()
	return &API{router: r}
}

// Run inits the config, loggers, dbs and all the repos and
// handlers for the app.
func (a *API) Run() error {
	var err error
	a.Config, err = config.LoadConfig()
	if err != nil {
		return err
	}

	a.Logger, err = log.NewLogger(a.Config)
	if err != nil {
		return err
	}

	a.DB, err = storage.NewSQLiteDB(a.Config.DBConfig)
	if err != nil {
		return err
	}

	// register repos
	repos := &Repositories{
		Service: repositories.NewServiceRepository(a.DB),
	}
	a.Repositories = repos

	// register handlers
	internalHandler := NewInternalHandler(a)
	servicesHandler := NewServicesHandler(a)

	a.router.Use(middleware.StripSlashes)

	a.router.Route("/v1", func(api chi.Router) {
		api.Mount("/_internal", internalHandler.GetRoutes())
		api.Mount("/services", servicesHandler.GetRoutes())
	})

	var routes []string
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		route = strings.Replace(route, "/*/", "/", -1)
		routes = append(routes, fmt.Sprintf("%s %s", method, route))
		return nil
	}

	if err := chi.Walk(a.router, walkFunc); err != nil {
		fmt.Printf("Logging err: %s\n", err.Error())
	}

	a.Logger.With(zap.Any("routes", routes)).Info("routes init")

	ctx := context.Background()
	srv := http.Server{
		Addr:    fmt.Sprintf(":%v", a.Config.APIConfig.Port),
		Handler: chi.ServerBaseContext(ctx, a.router),
	}
	a.Logger.Info(fmt.Sprintf("listening on :%v", a.Config.APIConfig.Port))

	// start graceful shutdown
	idleConnections := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		if err := srv.Shutdown(context.Background()); err != nil {
			a.Logger.Error(fmt.Sprintf("HTTP server Shutdown: %v", err))
		}
		close(idleConnections)
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		a.Logger.Error(fmt.Sprintf("HTTP server ListenAndServe: %v", err))
	}

	<-idleConnections
	return nil
}
