package server

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/usk81/aveo"

	"github.com/usk81/ashihara/service/holidays/interface/transport/di"
	httprouter "github.com/usk81/ashihara/service/holidays/interface/transport/http/router"
	"github.com/usk81/ashihara/shared/infrastructure/mysql"
	smw "github.com/usk81/ashihara/shared/interface/transport/http/middleware"
	"github.com/usk81/ashihara/shared/interface/transport/http/router"
	"github.com/usk81/ashihara/shared/interface/transport/http/server"
	"github.com/usk81/ashihara/shared/utils/logger"
)

type (
	Config struct {
		Address string

		Logger *slog.Logger

		Middlewares chi.Middlewares

		Context context.Context

		HTTPClient *http.Client
	}
)

var (
	ignoreUAs = []string{
		// kubernetes
		"kube-probe",
		// AWS
		"ELB-HealthChecker",
		// GCP
		"GoogleHC",
		// Azure
		"Edge Health Probe",
		"Load Balancer Agent",
	}
)

func Run(conf Config) {
	cf := conf
	setInitialConfig(&cf)

	r, err := newRouter(&cf)
	if err != nil {
		panic(err)
	}

	srv, err := server.New(cf.Address, cf.Logger, r)
	if err != nil {
		panic(err)
	}
	srv.Start()
}

func setInitialConfig(conf *Config) {
	if conf.Address == "" {
		conf.Address = ":80"
	}
	if conf.Logger == nil {
		conf.Logger = logger.New()
	}
	if len(conf.Middlewares) == 0 {
		conf.Middlewares = []func(http.Handler) http.Handler{
			middleware.RequestID,
			middleware.RealIP,
			middleware.Compress(5, "application/json"),
			smw.Logger(conf.Logger, ignoreUAs),
			middleware.Recoverer,
		}
	}
	if conf.Context == nil {
		conf.Context = context.Background()
	}
	if conf.HTTPClient == nil {
		conf.HTTPClient = &http.Client{
			Transport: &http.Transport{
				DialContext: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
				}).DialContext,
				TLSHandshakeTimeout:   10 * time.Second,
				ResponseHeaderTimeout: 10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
			},
			Timeout: 60 * time.Second,
		}
	}
}

func newRouter(conf *Config) (mux *chi.Mux, err error) {
	mux = router.Setup(conf.Middlewares...)

	// create database connections
	dr, err := mysql.ConnectReader(conf.Context, aveo.NewOs())
	if err != nil {
		return nil, err
	}
	dw, err := mysql.ConnectWriter(conf.Context, aveo.NewOs())
	if err != nil {
		return nil, err
	}

	// register more routes over here...
	httprouter.NewHolidays(di.Holidays(conf.HTTPClient, dr, dw, conf.Logger)).Route(mux) // nolint: errcheck

	router.LogRoutes(mux, conf.Logger)
	return mux, nil
}
