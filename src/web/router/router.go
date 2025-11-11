package router

import (
	"logger/config"
	"logger/web/controllers/health"
	"logger/web/controllers/log"
	"logger/web/middleware"
	"logger/web/server"

	"github.com/unrolled/secure"
)

// Router public struct
type Router struct {
	s *server.Server
}

// New Router
func New(s *server.Server) Router {
	return Router{s: s}
}

// Setup router
func (r *Router) Setup() {
	r.secure()

	r.s.R.Methods("OPTIONS").HandlerFunc(middleware.Options)

	health.New().SetupRouter(r.s)
	log.New().SetupRouter(r.s)
}

// CreateSubRouter with path
func (r *Router) createSubRouter(path string) *server.Server {
	return &server.Server{
		R:      r.s.R.PathPrefix(path).Subrouter(),
		Config: r.s.Config,
	}
}

func (r *Router) secure() {
	secureMiddleware := secure.New(config.Get().Propertyes.Security.Options)
	r.s.R.Use(secureMiddleware.Handler)
}
