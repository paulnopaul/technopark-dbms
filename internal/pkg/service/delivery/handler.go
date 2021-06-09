package delivery

import (
	"DBMSForum/internal/pkg/domain"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"net/http"
)

type serviceHandler struct {
}

func NewServiceHandler(r *router.Router, fu domain.ForumManager) {
	h := serviceHandler{}
	s := r.Group("/service")

	s.POST("/clear", h.serviceClearHandler)
	s.GET("/status", h.serviceStatusHandler)
}

func (handler *serviceHandler) serviceClearHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(http.StatusOK)
}

func (handler *serviceHandler) serviceStatusHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(http.StatusOK)
}
