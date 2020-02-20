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
	for _, m := range c.Middleware {
		n.Use(m)
	}
	c.Server = http.Server{
		Handler: n,
	}
	go func() {
		c.serving = true
		if err := c.Server.Serve(c.Listener); err != nil {
			c.serving = false
			log.Println("http server stopped")
		}
	}()
}

func (c *Server) ComponentDisable() {
	if c.serving {
		err := c.Server.Close()
		if err != nil {
			log.Fatal("SERVER DISABLE", err)
		}
	}
}
