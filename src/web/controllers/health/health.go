package health

import (
	"net/http"

	"logger/web/controllers"
	"logger/web/server"

	"github.com/joaopandolfi/blackwhale/handlers"
	"github.com/joaopandolfi/blackwhale/remotes/jaeger"
	"github.com/opentracing/opentracing-go"
)

// --- Health ---

type controller struct {
	s *server.Server
}

// New Health controller
func New() controllers.Controller {
	return &controller{
		s: nil,
	}
}

// Health route
func (c *controller) health(w http.ResponseWriter, r *http.Request) {
	_, span := jaeger.StartSpanFromRequest(opentracing.GlobalTracer(), r, "health")
	defer span.Finish()

	w.Header().Set("Access-Control-Allow-Origin", "*")
	handlers.Response(w, true, http.StatusOK)
}
