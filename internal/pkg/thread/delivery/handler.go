package delivery

import (
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"net/http"
	"technopark-dbms/internal/pkg/domain"
)

type threadHandler struct {
	threadUsecase domain.ThreadUsecase
}

func NewThreadHandler(r *router.Router, tu domain.ThreadUsecase) {
	h := threadHandler{
		threadUsecase: tu,
	}
	s := r.Group("/thread")

	s.POST("/{slug_or_id}/create", h.threadCreateHandler)
	s.GET("/{slug_or_id}/details", h.threadGetDetailsHandler)
	s.POST("/{slug_or_id}/details", h.threadUpdateDetailsHandler)
	s.GET("/{slug_or_id}/posts", h.threadGetPostsHandler)
	s.POST("/{slug_or_id}/vote", h.threadVoteHandler)
}

func (handler *threadHandler) threadCreateHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(http.StatusOK)
}

func (handler *threadHandler) threadGetDetailsHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(http.StatusOK)
}

func (handler *threadHandler) threadUpdateDetailsHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(http.StatusOK)
}

func (handler *threadHandler) threadGetPostsHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(http.StatusOK)
}

func (handler *threadHandler) threadVoteHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(http.StatusOK)
}
