package delivery

import (
	"DBMSForum/internal/pkg/domain"
	"fmt"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"net/http"
)

type Handler struct {
	forumUsecase domain.ForumUsecase
}

func NewForumHandler(r *router.Router, fu domain.ForumUsecase) {
	h := Handler{
		forumUsecase: fu,
	}
	s := r.Group("/forum")

	s.POST("/create", h.Create)

	s.GET("/{slug}/details", h.Details)
	s.POST("/{slug}/create", h.CreateThread)
	s.GET("/{slug}/users", h.Users)
	s.GET("/{slug}/threads", h.Threads)
}

func (handler *Handler) Create(ctx *fasthttp.RequestCtx) {
	body := ctx.PostBody()
	fmt.Println(string(body))
	ctx.SetStatusCode(http.StatusOK)
}

func (handler *Handler) Details(ctx *fasthttp.RequestCtx) {
	slugValue := ctx.UserValue("slug")
	fmt.Println(slugValue)
	ctx.SetStatusCode(http.StatusOK)
}

func (handler *Handler) CreateThread(ctx *fasthttp.RequestCtx) {
	slugValue := ctx.UserValue("slug")
	fmt.Println(slugValue)
	ctx.SetStatusCode(http.StatusOK)
}

func (handler *Handler) Users(ctx *fasthttp.RequestCtx) {
	slugValue := ctx.UserValue("slug")
	fmt.Println(slugValue)
	ctx.SetStatusCode(http.StatusOK)
}

func (handler *Handler) Threads(ctx *fasthttp.RequestCtx) {
	slugValue := ctx.UserValue("slug")
	fmt.Println(slugValue)
	ctx.SetStatusCode(http.StatusOK)
}
