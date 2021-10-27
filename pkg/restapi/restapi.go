package restapi

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/seggga/prometheus/pkg/http"
	"github.com/seggga/prometheus/pkg/storage"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

// RESTAPI represents a REST API business logic server.
type RESTAPI struct {
	server http.Server
	errors chan error
	logger *zap.SugaredLogger
}

// ServerConfig lets get some server parameters with toml-file.
type ServerConfig struct {
	Addr string `toml:"API_ADDR"`
}

// New returns a new instance of the REST API server.
func New(logger *zap.SugaredLogger, stor storage.CropURLStorage, srvAddr string) (*RESTAPI, error) {
	// define routes
	router := chi.NewRouter()
	router.Get("/", handler.Home(logger))
	router.Get("/{shortID}", handler.Redirect(stor, logger))
	router.Get("/links/{shortID}", handler.ViewStatistics(stor, logger))
	router.Post("/new-link", handler.NewLink(stor, logger))
	router.Post("/delete", handler.Delete(stor, logger))

	/*
		todo:
			/user/login
			/user/logout
	*/

	return &RESTAPI{
		server: http.Server{
			Addr:    srvAddr,
			Handler: router,
		},
		errors: make(chan error, 1),
		logger: logger,
	}, nil
}

// Start method starts the API server.
func (rapi *RESTAPI) Start() {

	// start monitoring service
	go http.ListenAndServe(":9102", promhttp.Handler())

	// start main service
	go func() {
		rapi.errors <- rapi.server.ListenAndServe()
		close(rapi.errors)
	}()
}

// Stop method stops API server.
func (rapi *RESTAPI) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return rapi.server.Shutdown(ctx)
}

// Notify returns a channel to notify the caller about errors.
// If you receive an error from the channel you should stop the application.
func (rapi *RESTAPI) Notify() <-chan error {
	return rapi.errors
}
