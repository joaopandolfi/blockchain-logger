package health

import "logger/web/server"

// SetupRouter -
func (c *controller) SetupRouter(s *server.Server) {
	c.s = s
	c.s.R.HandleFunc("/", c.health).Methods("POST", "GET", "HEAD")
	c.s.R.HandleFunc("/health", c.health).Methods("POST", "GET", "HEAD")
}
