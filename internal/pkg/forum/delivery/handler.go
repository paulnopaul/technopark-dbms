package delivery

import (
	"DBMSForum/internal/pkg/domain"
	"fmt"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"net/http"
)

type forumHandler struct {
	forumUsecase domain.ForumManager
}

func NewForumHandler(r *router.Router, fu domain.ForumManager) {
	h := forumHandler{
		forumUsecase: fu,
	}
	s := r.Group("/forum")

	s.POST("/create", h.forumCreateHandler)

	s.GET("/{slug}/details", h.forumDetailsHandler)
	s.POST("/{slug}/create", h.forumCreateThreadHanlder)
	s.GET("/{slug}/users", h.forumGetUsersHandler)
	s.GET("/{slug}/threads", h.forumGetThreadsHandler)
}

// Create
/*
curl --header "Content-Type: application/json" \
--request POST \
--data '{"user":"newUser","title":"newForum","slug":"new-forum"}' \
http://localhost:5000/forum/create
*/
func (handler *forumHandler) forumCreateHandler(ctx *fasthttp.RequestCtx) {
	body := ctx.PostBody()
	fmt.Println(string(body))
	ctx.SetStatusCode(http.StatusOK)
}

func (handler *forumHandler) forumDetailsHandler(ctx *fasthttp.RequestCtx) {
	slugValue := ctx.UserValue("slug")
	fmt.Println(slugValue)
	ctx.SetStatusCode(http.StatusOK)
}

func (handler *forumHandler) forumCreateThreadHanlder(ctx *fasthttp.RequestCtx) {
	slugValue := ctx.UserValue("slug")
	fmt.Println(slugValue)
	ctx.SetStatusCode(http.StatusOK)
}

func (handler *forumHandler) forumGetUsersHandler(ctx *fasthttp.RequestCtx) {
	slugValue := ctx.UserValue("slug")
	fmt.Println(slugValue)
	ctx.SetStatusCode(http.StatusOK)
}

func (handler *forumHandler) forumGetThreadsHandler(ctx *fasthttp.RequestCtx) {
	slugValue := ctx.UserValue("slug")
	fmt.Println(slugValue)
	ctx.SetStatusCode(http.StatusOK)
}
