package router

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/usk81/toolkit/slice"
)

type (
	// HTTPRouter ...
	HTTPRouter interface {
		Route(*chi.Mux) error
	}

	// Endpoint ...
	Endpoint struct {
		Middlewares []func(http.Handler) http.Handler
		Handler     http.HandlerFunc
	}

	// EndpointPattern ...
	EndpointPattern struct {
		Pattern   string
		Endpoints map[string]Endpoint
	}

	// Route ...
	Route struct {
		Middlewares []func(http.Handler) http.Handler
		Endpoints   []EndpointPattern
		SubRoutes   []SubRoute
	}

	// SubRoute ...
	SubRoute struct {
		Pattern string
		Route   Route
	}
)

var methods = []string{
	"All",
	http.MethodGet,
	http.MethodHead,
	http.MethodPost,
	http.MethodPut,
	http.MethodPatch,
	http.MethodDelete,
	http.MethodConnect,
	http.MethodOptions,
	http.MethodTrace,
}

// Setup a new go-chi router
func Setup(middlewares ...func(http.Handler) http.Handler) *chi.Mux {
	r := chi.NewRouter()

	// A good base middleware stack
	for _, m := range middlewares {
		r.Use(m)
	}

	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	return r
}

// New sets a Route
func New(rt Route) *Route {
	return &rt
}

// Build router
func (rt *Route) Build(r *chi.Mux) (err error) {
	return build(r, rt)
}

func build(r *chi.Mux, rt *Route) (err error) {
	for _, m := range rt.Middlewares {
		r.Use(m)
	}
	for _, ep := range rt.Endpoints {
		if err = buildEndpoints(r, ep); err != nil {
			return
		}
	}
	for _, sr := range rt.SubRoutes {
		if err = buildSubroute(r, sr); err != nil {
			return
		}
	}
	return
}

func buildEndpoints(r *chi.Mux, ep EndpointPattern) (err error) {
	for method, e := range ep.Endpoints {
		if !slice.Exists(method, methods) {
			return fmt.Errorf("invalid http method : %s", method)
		}
		if method == "All" {
			r.With(e.Middlewares...).HandleFunc(ep.Pattern, e.Handler)
		} else {
			r.With(e.Middlewares...).MethodFunc(method, ep.Pattern, e.Handler)
		}
	}
	return nil
}

func buildSubroute(r *chi.Mux, sr SubRoute) (err error) {
	mux := chi.NewRouter()
	if err = build(mux, &sr.Route); err != nil {
		return
	}
	r.Mount(sr.Pattern, mux)
	return
}

// LogRoutes logs routing on Debug level when server is ONLINE
func LogRoutes(r *chi.Mux, logger *slog.Logger) {
	if err := chi.Walk(r, loggingPrintRoute(logger)); err != nil {
		logger.Error("Failed to walk routes:", slog.Any("error", err))
	}
}

func loggingPrintRoute(logger *slog.Logger) chi.WalkFunc {
	return func(method string, route string, _ http.Handler, _ ...func(http.Handler) http.Handler) error {
		route = strings.ReplaceAll(route, "/*/", "/")
		logger.Debug("Registering route", slog.String("method", method), slog.String("route", route))
		return nil
	}
}
