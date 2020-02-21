package http

import (
	"net/http"

	"github.com/goji/httpauth"
)

type SingleUserBasicAuth struct {
	Username string
	Password string
}

func (c *SingleUserBasicAuth) Middleware() func(http.Handler) http.Handler {
	return httpauth.SimpleBasicAuth(c.Username, c.Password)
}
