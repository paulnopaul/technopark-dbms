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

func NewUserHandler(r *router.Group, uc domain.UserUsecase) {
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
		utilities.Resp(ctx, http.StatusBadRequest, errors.JSONDecodeErrorMessage)
		return
	}

	nickname := ctx.UserValue("nickname").(string)

	createdUser, err, alreadyCreatedUsers := handler.userUsecase.CreateUser(nickname, *parsedUser)
	if err != nil {
		log.WithError(err).Error("user creation error")
		var errorMsg []byte
		if err == user.AlreadyExistsError {
			errorMsg, _ = json.Marshal(alreadyCreatedUsers)
		} else {
			errorMsg = errors.JSONErrorMessage(err)
		}
		utilities.Resp(ctx, user.CodeFromError(err), errorMsg)
		return
	}
	body, err := json.Marshal(createdUser)
	if err != nil {
		log.WithError(err).Error(errors.JSONEncodeError)
		utilities.Resp(ctx, http.StatusInternalServerError, errors.JSONEncodeErrorMessage)
		return
	}

	utilities.Resp(ctx, fasthttp.StatusCreated, body)
}

func (handler *userHandler) userGetProfileHandler(ctx *fasthttp.RequestCtx) {
	nickname := ctx.UserValue("nickname").(string)

	foundUser, err := handler.userUsecase.GetProfile(nickname)
	if err != nil {
		log.WithError(err).Error("user get details error")
		utilities.Resp(ctx, user.CodeFromError(err), errors.JSONErrorMessage(err))
		return
	}

	body, err := json.Marshal(foundUser)
	if err != nil {
		log.WithError(err).Error(errors.JSONEncodeError)
		utilities.Resp(ctx, fasthttp.StatusInternalServerError, errors.JSONEncodeErrorMessage)
		return
	}

	ctx.Success("application/json", body)
	ctx.SetStatusCode(http.StatusOK)
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

	updatedUser, err := handler.userUsecase.UpdateProfile(nickname, *parsedUser)
	if err != nil {
		log.WithError(err).WithFields(log.Fields{"nickname": nickname}).Error("user updating error")
		utilities.Resp(ctx, user.CodeFromError(err), errors.JSONErrorMessage(err))
		return
	}

	body, err := json.Marshal(updatedUser)
	if err != nil {
		log.WithError(err).Error(errors.JSONEncodeError)
		utilities.Resp(ctx, fasthttp.StatusInternalServerError, errors.JSONEncodeErrorMessage)
		return
	}

	utilities.Resp(ctx, fasthttp.StatusOK, body)
}
