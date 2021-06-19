package delivery

import (
	"technopark-dbms/internal/pkg/domain"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"net/http"
)

type userHandler struct {
	userUsecase domain.UserUsecase
}

func NewUserHandler(r *router.Router, uc domain.UserUsecase) {
	h := userHandler{
		userUsecase: uc,
	}
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
