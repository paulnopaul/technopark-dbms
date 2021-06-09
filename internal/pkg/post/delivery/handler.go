package delivery

import (
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"net/http"
)

type postHandler struct {
}

func NewServiceHandler(r *router.Router) {
	h := postHandler{}
	s := r.Group("/post")

	s.GET("/{id}/details", h.postGetDetailsHandler)
	s.POST("/{id}/details", h.postUpdateDetailsHandler)
}

func (handler *postHandler) postGetDetailsHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(http.StatusOK)
}

func (handler *postHandler) postUpdateDetailsHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(http.StatusOK)
}
