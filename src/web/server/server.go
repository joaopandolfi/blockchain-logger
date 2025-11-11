package server

import (
	"context"
	"net/http"

	"logger/config"

	"github.com/gorilla/mux"
	"github.com/joaopandolfi/blackwhale/utils"
)

// Server web
type Server struct {
	R      *mux.Router
	Config config.Config
	srv    *http.Server
}

// New server
func New(r *mux.Router, conf config.Config) *Server {
	// Bind to a port and pass our router in
	utils.Info("Server listenning on", conf.Propertyes.Port)
	srv := &http.Server{
		Handler:      r,
		Addr:         conf.Propertyes.Port,
		WriteTimeout: conf.Propertyes.Timeout.Write,
		ReadTimeout:  conf.Propertyes.Timeout.Read,
	}

	return &Server{
		R:      r,
		Config: conf,
		srv:    srv,
	}
}

// Start Web server
func (s *Server) Start() {

	var err error
	if config.Get().Propertyes.Security.Debug {
		err = s.srv.ListenAndServe()
	} else {
		err = s.srv.ListenAndServeTLS(config.Get().Propertyes.Security.TLSCert, config.Get().Propertyes.Security.TLSKey)
	}
	if err != nil && err != http.ErrServerClosed {
		utils.CriticalError("Fatal server error", err.Error())
	}
}

// Shutdown server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
