package http

import (
	"log"
	"net"
	"net/http"

	"github.com/urfave/negroni"
)

type Middleware interface {
	ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}

type Server struct {
	http.Server

	Listener   net.Listener
	Handler    http.Handler
	Middleware []Middleware

	serving bool
}

func (c *Server) ComponentEnable() {
	if c.Listener == nil || c.Handler == nil {
		return
	}
	log.Println("starting http server")
	n := negroni.New()
	n.UseHandler(c.Handler)
	c.Server = http.Server{
		Handler: n,
	}
	go func() {
		c.serving = true
		if err := c.Server.Serve(c.Listener); err != nil {
			c.serving = false
			log.Fatal(err)
		}
	}()
}

func (c *Server) ComponentDisable() {
	if c.serving {
		c.Server.Close()
	}
}
