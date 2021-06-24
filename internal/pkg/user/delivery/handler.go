package delivery

import (
	"encoding/json"
	"github.com/fasthttp/router"
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"net/http"
	"technopark-dbms/internal/pkg/domain"
	"technopark-dbms/internal/pkg/errors"
	"technopark-dbms/internal/pkg/user"
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
	parsedUser := &domain.User{}
	err := json.Unmarshal(ctx.PostBody(), parsedUser)
	if err != nil {
		log.WithError(err).Error(errors.JSONUnmarshallError)
		ctx.Error(errors.JSONDecodeErrorMessage, fasthttp.StatusInternalServerError)
		return
	}

	nickname := ctx.UserValue("nickname").(string)

	createdUser, err := handler.userUsecase.CreateUser(nickname, *parsedUser)
	if err != nil {
		log.WithError(err).Error("user creation error")
		ctx.Error(errors.JSONErrorMessage(err), user.CodeFromError(err))
		return
	}

	if err = json.NewEncoder(ctx).Encode(createdUser); err != nil {
		log.WithError(err).Error(errors.JSONEncodeError)
		ctx.Error(errors.JSONEncodeErrorMessage, fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(http.StatusCreated)
}

func (handler *userHandler) userGetProfileHandler(ctx *fasthttp.RequestCtx) {
	nickname := ctx.UserValue("nickname").(string)

	createdUser, err := handler.userUsecase.GetProfile(nickname)
	if err != nil {
		log.WithError(err).Error("user get details error")
		ctx.Error(errors.JSONErrorMessage(err), user.CodeFromError(err))
		return
	}

	if err = json.NewEncoder(ctx).Encode(createdUser); err != nil {
		log.WithError(err).Error(errors.JSONEncodeError)
		ctx.Error(errors.JSONEncodeErrorMessage, fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(http.StatusCreated)
}

func (handler *userHandler) userUpdateProfileHandler(ctx *fasthttp.RequestCtx) {
	parsedUser := &domain.User{}
	err := json.Unmarshal(ctx.PostBody(), parsedUser)
	if err != nil {
		log.WithError(err).Error(errors.JSONUnmarshallError)
		ctx.Error(errors.JSONDecodeErrorMessage, fasthttp.StatusInternalServerError)
		return
	}

	nickname := ctx.UserValue("nickname").(string)

	createdUser, err := handler.userUsecase.UpdateProfile(nickname, *parsedUser)
	if err != nil {
		log.WithError(err).Error("user updating error")
		ctx.Error(errors.JSONErrorMessage(err), user.CodeFromError(err))
		return
	}

	if err = json.NewEncoder(ctx).Encode(createdUser); err != nil {
		log.WithError(err).Error(errors.JSONEncodeError)
		ctx.Error(errors.JSONEncodeErrorMessage, fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(http.StatusOK)
}
