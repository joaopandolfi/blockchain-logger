package log

import "logger/web/server"

// SetupRouter -
func (c *controller) SetupRouter(s *server.Server) {
	c.s = s
	c.s.R.HandleFunc("/log", c.newLog).Methods("POST", "HEAD")
	c.s.R.HandleFunc("/validate", c.validate).Methods("GET", "HEAD")
	c.s.R.HandleFunc("/validate/{init:[0-9]+}/{end:[0-9]+}", c.validateSegment).Methods("GET", "HEAD")
}
