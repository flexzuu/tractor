package http

import (
	"log"
	"net"
	"net/http"

	"github.com/urfave/negroni"
)

type Server struct {
	Listener net.Listener
	Handler  http.Handler

	s *http.Server
}

func (c *Server) ComponentEnable() {
	if c.Listener == nil || c.Handler == nil {
		return
	}
	log.Println("starting http server")
	n := negroni.New()
	n.UseHandler(c.Handler)
	c.s = &http.Server{
		Handler: n,
	}
	go func() {
		if err := c.s.Serve(c.Listener); err != nil {
			log.Fatal(err)
		}
	}()
}

func (c *Server) ComponentDisable() {
	if c.s != nil {
		c.s.Close()
	}
}
