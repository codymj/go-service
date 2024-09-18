package router

import (
	"github.com/julienschmidt/httprouter"
	"go-service.codymj.io/cmd/app/router/user"
	"go-service.codymj.io/cmd/app/util"
)

const (
	ApiVersion = "/v1"
)

type Router struct {
	Router *httprouter.Router
}

func New() *Router {
	return &Router{}
}

func (r *Router) Setup(services util.Services) error {
	r.Router = httprouter.New()

	uh := user.New(services)
	uh.Init(r.Router, ApiVersion)

	return nil
}
