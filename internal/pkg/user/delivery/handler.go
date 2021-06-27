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
	"technopark-dbms/internal/pkg/utilities"
)

type userHandler struct {
	userUsecase domain.UserUsecase
}

func NewUserHandler(r *router.Router, uc domain.UserUsecase) {
	h := userHandler{
		userUsecase: uc,
	}
	s := r.Group("/api/user")

	s.POST("/{nickname}/create", h.userCreateHandler)
	s.GET("/{nickname}/profile", h.userGetProfileHandler)
	s.POST("/{nickname}/profile", h.userUpdateProfileHandler)
}

func (handler *userHandler) userCreateHandler(ctx *fasthttp.RequestCtx) {
	parsedUser := &domain.User{}
	err := json.Unmarshal(ctx.PostBody(), parsedUser)
	if err != nil {
		log.WithError(err).Error(errors.JSONUnmarshallError)
		utilities.Resp(ctx, http.StatusBadRequest, errors.JSONDecodeErrorMessage)
		return
	}

	nickname := ctx.UserValue("nickname").(string)

	createdUser, err, alreadyCreatedUsers := handler.userUsecase.CreateUser(nickname, *parsedUser)
	if err != nil {
		log.WithError(err).Error("user creation error")
		if err == user.AlreadyExistsError {
			utilities.Resp(ctx, user.CodeFromError(err), alreadyCreatedUsers)
		} else {
			utilities.Resp(ctx, user.CodeFromError(err), errors.JSONErrorMessage(err))
		}
		return
	}

	utilities.Resp(ctx, fasthttp.StatusCreated, createdUser)
}

func (handler *userHandler) userGetProfileHandler(ctx *fasthttp.RequestCtx) {
	nickname := ctx.UserValue("nickname").(string)

	foundUser, err := handler.userUsecase.GetProfile(nickname)
	if err != nil {
		log.WithError(err).Error("user get details error")
		utilities.Resp(ctx, user.CodeFromError(err), errors.JSONErrorMessage(err))
		return
	}
	utilities.Resp(ctx, http.StatusOK, foundUser)
}

func (handler *userHandler) userUpdateProfileHandler(ctx *fasthttp.RequestCtx) {
	parsedUser := &domain.User{}
	err := json.Unmarshal(ctx.PostBody(), parsedUser)
	if err != nil {
		log.WithError(err).Error(errors.JSONUnmarshallError)
		utilities.Resp(ctx, fasthttp.StatusInternalServerError, errors.JSONDecodeErrorMessage)
		return
	}

	nickname := ctx.UserValue("nickname").(string)

	updatedUser, err := handler.userUsecase.UpdateUser(nickname, *parsedUser)
	if err != nil {
		log.WithError(err).WithFields(log.Fields{"nickname": nickname}).Error("user updating error")
		utilities.Resp(ctx, user.CodeFromError(err), errors.JSONErrorMessage(err))
		return
	}
	utilities.Resp(ctx, fasthttp.StatusOK, updatedUser)
}
