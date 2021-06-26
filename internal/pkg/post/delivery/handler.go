package delivery

import (
	"encoding/json"
	"github.com/fasthttp/router"
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"strconv"
	"strings"
	"technopark-dbms/internal/pkg/domain"
	"technopark-dbms/internal/pkg/errors"
	"technopark-dbms/internal/pkg/post"
	"technopark-dbms/internal/pkg/utilities"
)

type postHandler struct {
	postUsecase domain.PostUsecase
}

func NewPostHandler(r *router.Group, pu domain.PostUsecase) {
	h := postHandler{
		postUsecase: pu,
	}
	s := r.Group("/post")

	s.GET("/{id:[0-9]+}/details", h.postGetDetailsHandler)
	s.POST("/{id:[0-9]+}/details", h.postUpdateDetailsHandler)
}

func parseRelated(queryArgs *fasthttp.Args) (userRelated, forumRelated, threadRelated bool) {
	values := strings.Split(string(queryArgs.Peek("related")), ",")
	if values != nil {
		for _, value := range values {
			if string(value) == "user" {
				userRelated = true
			} else if string(value) == "forum" {
				forumRelated = true
			} else if string(value) == "thread" {
				threadRelated = true
			}
		}
	}
	return
}

type postFull struct {
	Post   *domain.Post   `json:"post"`
	Forum  *domain.Forum  `json:"forum,omitempty"`
	Thread *domain.Thread `json:"thread,omitempty"`
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

	foundPost, foundForum, foundThread, foundUser, err := handler.postUsecase.GetPostDetails(postId, userRelated, forumRelated, threadRelated)
	if err != nil {
		log.WithError(err).Error("post get details error")
		if err == post.NotFoundError {
			utilities.Resp(ctx, fasthttp.StatusNotFound, errors.JSONErrorMessage(err))
			return
		}
		utilities.Resp(ctx, fasthttp.StatusInternalServerError, errors.JSONErrorMessage(err))
		return
	}

	postFull := postFull{foundPost, foundForum, foundThread, foundUser}
	utilities.Resp(ctx, fasthttp.StatusOK, postFull)
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

	foundPost, err := handler.postUsecase.UpdatePostDetails(postId, *parsedPost)
	if err != nil {
		log.WithError(err).Error("forum update details error")
		if err == post.NotFoundError {
			utilities.Resp(ctx, fasthttp.StatusNotFound, errors.JSONErrorMessage(err))
			return
		}
		utilities.Resp(ctx, fasthttp.StatusInternalServerError, errors.JSONEncodeErrorMessage)
		return
	}
	utilities.Resp(ctx, fasthttp.StatusOK, foundPost)
}
