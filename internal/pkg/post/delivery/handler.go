package delivery

import (
	"technopark-dbms/internal/pkg/domain"
	"technopark-dbms/internal/pkg/errors"
	"encoding/json"
	"github.com/fasthttp/router"
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"net/http"
	"strconv"
)

type postHandler struct {
	postUsecase domain.PostUsecase
}

func NewPostHandler(r *router.Router, pu domain.PostUsecase) {
	h := postHandler{}
	s := r.Group("/post")

	s.GET("/{id:[0-9]+}/details", h.postGetDetailsHandler)
	s.POST("/{id:[0-9]+}/details", h.postUpdateDetailsHandler)
}

func parseRelated(queryArgs *fasthttp.Args) (userRelated, forumRelated, threadRelated bool) {
	userRelated = queryArgs.Has("user")
	forumRelated = queryArgs.Has("forum")
	threadRelated = queryArgs.Has("thread")
	return
}

func (handler *postHandler) postGetDetailsHandler(ctx *fasthttp.RequestCtx) {
	postId, err := strconv.ParseInt(ctx.UserValue("id").(string), 10, 64)
	if err != nil {
		log.WithError(err).Error(errors.URLParamsError)
		ctx.Error(errors.JSONURLParamsErrorMessage, errors.CodeFromJSONMessage(errors.JSONURLParamsErrorMessage))
		return
	}
	userRelated, forumRelated, threadRelated := parseRelated(ctx.QueryArgs())

	foundPost, err := handler.postUsecase.GetDetails(postId, userRelated, forumRelated, threadRelated)
	if err != nil {
		log.WithError(err).Error("post get details error")
		// todo error + message
		return
	}

	if err = json.NewEncoder(ctx).Encode(foundPost); err != nil {
		log.WithError(err).Error(errors.JSONEncodeError)
		ctx.Error(errors.JSONEncodeErrorMessage, fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(http.StatusOK)

}

func (handler *postHandler) postUpdateDetailsHandler(ctx *fasthttp.RequestCtx) {
	postId, err := strconv.ParseInt(ctx.UserValue("id").(string), 10, 64)
	if err != nil {
		log.WithError(err).Error(errors.URLParamsError)
		ctx.Error(errors.JSONURLParamsErrorMessage, errors.CodeFromJSONMessage(errors.JSONURLParamsErrorMessage))
		return
	}
	parsedPost := &domain.Post{}
	err = json.Unmarshal(ctx.PostBody(), parsedPost)
	if err != nil {
		log.WithError(err).Error(errors.JSONUnmarshallError)
		ctx.Error(errors.JSONDecodeErrorMessage, fasthttp.StatusInternalServerError)
		return
	}

	foundPost, err := handler.postUsecase.UpdateDetails(postId, *parsedPost)
	if err != nil {
		log.WithError(err).Error("forum get details error")
		// todo error + message
		return
	}

	if err = json.NewEncoder(ctx).Encode(foundPost); err != nil {
		log.WithError(err).Error(errors.JSONEncodeError)
		ctx.Error(errors.JSONEncodeErrorMessage, fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(http.StatusCreated)
}
