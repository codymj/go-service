package user

import (
	"github.com/julienschmidt/httprouter"
	"go-service.codymj.io/cmd/app/util"
)

type handler struct {
	services util.Services
}

type Handler interface {
	Init(r *httprouter.Router, apiVersion string)
}

func New(services util.Services) Handler {
	return &handler{
		services: services,
	}
}

func (h *handler) Init(r *httprouter.Router, apiVersion string) {
	usersPath := apiVersion + "/users"
	usersIdPath := usersPath + "/:id"

	r.GET(usersIdPath, h.getById)
	r.GET(usersPath, h.getByParams)
}
