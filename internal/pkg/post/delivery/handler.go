package delivery

import (
	"encoding/json"
	"github.com/fasthttp/router"
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"net/http"
	"strconv"
	"technopark-dbms/internal/pkg/domain"
	"technopark-dbms/internal/pkg/errors"
	"technopark-dbms/internal/pkg/forum"
	"technopark-dbms/internal/pkg/utilities"
)

type postHandler struct {
	postUsecase domain.PostUsecase
}

func NewPostHandler(r *router.Router, pu domain.PostUsecase) {
	h := postHandler{
		postUsecase: pu,
	}
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

type postFull struct {
	Post   *domain.Post   `json:"post"`
	Forum  *domain.Forum  `json:",omitempty"`
	Thread *domain.Thread `json:",omitempty"`
	User   *domain.User   `json:"author,omitempty"`
}

func (handler *postHandler) postGetDetailsHandler(ctx *fasthttp.RequestCtx) {
	postId, err := strconv.ParseInt(ctx.UserValue("id").(string), 10, 64)
	if err != nil {
		log.WithError(err).Error(errors.URLParamsError)
		utilities.Resp(ctx, errors.CodeFromDeliveryError(errors.URLParamsError), errors.JSONURLParamsErrorMessage)
		return
	}
	userRelated, forumRelated, threadRelated := parseRelated(ctx.QueryArgs())

	foundPost, foundForum, foundThread, foundUser, err := handler.postUsecase.GetDetails(postId, userRelated, forumRelated, threadRelated)
	if err != nil {
		log.WithError(err).Error("post get details error")
		if err != forum.NotFound {

		}
		return
	}

	body, err := json.Marshal(postFull{foundPost, foundForum, foundThread, foundUser})
	if err != nil {
		log.WithError(err).Error(errors.JSONEncodeError)
		utilities.Resp(ctx, http.StatusInternalServerError, errors.JSONEncodeErrorMessage)
		return
	}

	utilities.Resp(ctx, fasthttp.StatusOK, body)
}

func (handler *postHandler) postUpdateDetailsHandler(ctx *fasthttp.RequestCtx) {
	postId, err := strconv.ParseInt(ctx.UserValue("id").(string), 10, 64)
	if err != nil {
		log.WithError(err).Error(errors.URLParamsError)
		utilities.Resp(ctx, errors.CodeFromDeliveryError(errors.URLParamsError), errors.JSONURLParamsErrorMessage)
		return
	}

	parsedPost := &domain.Post{}
	err = json.Unmarshal(ctx.PostBody(), parsedPost)
	if err != nil {
		log.WithError(err).Error(errors.JSONUnmarshallError)
		utilities.Resp(ctx, fasthttp.StatusInternalServerError, errors.JSONDecodeErrorMessage)
		return
	}

	foundPost, err := handler.postUsecase.UpdateDetails(postId, *parsedPost)
	if err != nil {
		log.WithError(err).Error("forum get details error")

		return
	}

	body, err := json.Marshal(foundPost)
	if err != nil {
		log.WithError(err).Error(errors.JSONEncodeError)
		utilities.Resp(ctx, fasthttp.StatusInternalServerError, errors.JSONEncodeErrorMessage)
		return
	}
	utilities.Resp(ctx, fasthttp.StatusCreated, body)
}
