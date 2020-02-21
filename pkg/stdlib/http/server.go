package http

import (
	"log"
	"net"
	"net/http"
)

type Middleware interface {
	Middleware() func(http.Handler) http.Handler
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
	h := c.Handler
	for _, mw := range c.Middleware {
		adapter := mw.Middleware()
		h = adapter(h)
	}
	c.Server = http.Server{
		Handler: h,
	}
	go func() {
		c.serving = true
		if err := c.Server.Serve(c.Listener); err != nil {
			c.serving = false
			log.Println("http server stopped:", err)
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
