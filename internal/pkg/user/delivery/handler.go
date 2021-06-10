package delivery

import (
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"net/http"
)

type userHandler struct {
}

func NewUserHandler(r *router.Router) {
	h := userHandler{}
	s := r.Group("/user")

	s.POST("/{nickname}/create", h.userCreateHandler)
	s.GET("/{nickname}/profile", h.userGetProfileHandler)
	s.POST("/{nickname}/profile", h.userUpdateProfileHandler)
}

func (handler *userHandler) userCreateHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(http.StatusOK)
}

func (handler *userHandler) userGetProfileHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(http.StatusOK)
}

func (handler *userHandler) userUpdateProfileHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(http.StatusOK)
}
