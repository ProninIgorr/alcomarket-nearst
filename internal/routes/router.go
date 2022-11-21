package routes

import (
	"github.com/ProninIgorr/alcomarket-nearst/internal/handlers"
	"github.com/ProninIgorr/alcomarket-nearst/internal/store"
	"github.com/ProninIgorr/alcomarket-nearst/pkg/helpers/probs"
	"github.com/ProninIgorr/alcomarket-nearst/pkg/middleware"
	"github.com/d-kolpakov/logger"
	"github.com/go-chi/chi"
	"github.com/heptiolabs/healthcheck"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"net/http"
)

func Start(serviceName string, appVersion string, logger *logger.Logger, db *pgxpool.Pool, s *store.Store) {
	serviceUtlPrefix := "/" + serviceName

	r := chi.NewRouter()
	m := middleware.New(logger)
	h := &handlers.Handler{Db: db, L: logger, ServiceName: serviceName, AppVersion: appVersion, Store: s}

	// Health check
	hc := healthcheck.NewHandler()
	hc.AddReadinessCheck("readiness-check", func() error {
		return probs.GetReadinessErr()
	})
	hc.AddLivenessCheck("liveness-check", func() error {
		return probs.GetLivenessErr()
	})

	r.Use(m.RecoverMiddleware)
	r.Route(serviceUtlPrefix, func(r chi.Router) {
		//System endpoints
		r.Group(func(r chi.Router) {
			r.Get("/", h.HomeRouteHandler)
			r.Get(serviceUtlPrefix+"/ready/", hc.ReadyEndpoint)
			r.Get(serviceUtlPrefix+"/live/", hc.LiveEndpoint)
		})

		// Application endpoints
		r.Group(func(r chi.Router) {
			r.Use(m.ContextRequestMiddleware, m.LogRequests)
			r.Get("/endpoint/", h.InternalEndpoint)
			//Public endpoints
			r.Route("/public", func(r chi.Router) {
				r.Get("/stores/", h.Stores)
				r.Get("/slow/", h.Slow)
			})
		})
	})

	log.Fatal(http.ListenAndServe(":8080", r))
}
