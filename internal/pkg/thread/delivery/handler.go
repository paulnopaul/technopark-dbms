package delivery

import (
	"encoding/json"
	"fmt"
	"github.com/fasthttp/router"
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"net/http"
	"technopark-dbms/internal/pkg/domain"
	"technopark-dbms/internal/pkg/errors"
	"technopark-dbms/internal/pkg/forum"
	"technopark-dbms/internal/pkg/thread"
	"technopark-dbms/internal/pkg/utilities"
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
	slugOrId := utilities.NewSlugOrId(ctx.UserValue("slug_or_id").(string))
	var parsedPosts []domain.Post
	err := json.Unmarshal(ctx.PostBody(), &parsedPosts)
	if err != nil {
		log.WithError(err).Error(errors.JSONUnmarshallError)
		utilities.Resp(ctx, http.StatusBadRequest, errors.JSONDecodeErrorMessage)
		return
	}

	createdPosts, err := handler.threadUsecase.CreatePosts(slugOrId, parsedPosts)
	responseStatus := fasthttp.StatusCreated
	if err != nil {
		log.WithError(err).Error("post creation error")
		if err == thread.AlreadyExists {
			responseStatus = fasthttp.StatusConflict
		} else {
			utilities.Resp(ctx, fasthttp.StatusInternalServerError, errors.JSONErrorMessage(err))
			return
		}
	}

	body, err := json.Marshal(createdPosts)
	if err != nil {
		log.WithError(err).Error(errors.JSONEncodeError)
		utilities.Resp(ctx, fasthttp.StatusInternalServerError, errors.JSONEncodeErrorMessage)
		return
	}

	utilities.Resp(ctx, responseStatus, body)
}

func (handler *threadHandler) threadGetDetailsHandler(ctx *fasthttp.RequestCtx) {
	slugOrId := utilities.NewSlugOrId(ctx.UserValue("slug_or_id").(string))
	forumDetails, err := handler.threadUsecase.GetThreadDetails(slugOrId)
	if err != nil {
		log.WithError(err).Error("thread get details error")
		if err == forum.NotFound {
			utilities.Resp(ctx,
				fasthttp.StatusNotFound,
				errors.JSONMessage(fmt.Sprintf("Can't find thread with slug: %s", slugOrId.Slug)))
			return
		} else {
			utilities.Resp(ctx,
				fasthttp.StatusInternalServerError,
				errors.JSONErrorMessage(err))
			return
		}
	}

	body, err := json.Marshal(forumDetails)
	if err != nil {
		log.WithError(err).Error(errors.JSONEncodeError)
		utilities.Resp(ctx, fasthttp.StatusInternalServerError, errors.JSONEncodeErrorMessage)
		return
	}

	utilities.Resp(ctx, fasthttp.StatusOK, body)
}

func (handler *threadHandler) threadUpdateDetailsHandler(ctx *fasthttp.RequestCtx) {
	slugOrId := utilities.NewSlugOrId(ctx.UserValue("slug_or_id").(string))
	parsedThread := &domain.Thread{}
	err := json.Unmarshal(ctx.PostBody(), parsedThread)
	if err != nil {
		log.WithError(err).Error(errors.JSONUnmarshallError)
		utilities.Resp(ctx, http.StatusBadRequest, errors.JSONDecodeErrorMessage)
		return
	}

	updatedThread, err := handler.threadUsecase.UpdateThreadDetails(slugOrId, *parsedThread)
	if err != nil {
		log.WithError(err).Error("thread update error")
		if err == thread.NotFound {
			utilities.Resp(ctx, fasthttp.StatusNotFound, errors.JSONErrorMessage(err))
			return
		} else {
			utilities.Resp(ctx, fasthttp.StatusInternalServerError, errors.JSONErrorMessage(err))
			return
		}
	}

	body, err := json.Marshal(updatedThread)
	if err != nil {
		log.WithError(err).Error(errors.JSONEncodeError)
		utilities.Resp(ctx, fasthttp.StatusInternalServerError, errors.JSONEncodeErrorMessage)
		return
	}

	utilities.Resp(ctx, http.StatusOK, body)
}

func (handler *threadHandler) threadGetPostsHandler(ctx *fasthttp.RequestCtx) {
	slugOrId := utilities.NewSlugOrId(ctx.UserValue("slug_or_id").(string))
	params, err := utilities.NewArrayOutParams(ctx.QueryArgs())
	if err != nil {
		log.WithError(err).Error(errors.QuerystringParseError)
		utilities.Resp(ctx, fasthttp.StatusInternalServerError, errors.JSONErrorMessage(err))
		return
	}

	foundPosts, err := handler.threadUsecase.GetThreadPosts(slugOrId, *params)
	if err != nil {
		log.WithError(err).Error("post find error")
		if err == thread.NotFound {
			utilities.Resp(ctx, fasthttp.StatusNotFound, errors.JSONErrorMessage(err))
			return
		} else {
			utilities.Resp(ctx, fasthttp.StatusInternalServerError, errors.JSONErrorMessage(err))
			return
		}
	}

	body, err := json.Marshal(foundPosts)
	if err != nil {
		log.WithError(err).Error(errors.JSONEncodeError)
		utilities.Resp(ctx, http.StatusInternalServerError, errors.JSONEncodeErrorMessage)
		return
	}

	utilities.Resp(ctx, fasthttp.StatusOK, body)
}

func (handler *threadHandler) threadVoteHandler(ctx *fasthttp.RequestCtx) {
	slugOrId := utilities.NewSlugOrId(ctx.UserValue("slug_or_id").(string))
	parsedVote := &domain.Vote{}
	err := json.Unmarshal(ctx.PostBody(), parsedVote)
	if err != nil {
		log.WithError(err).Error(errors.JSONUnmarshallError)
		utilities.Resp(ctx, http.StatusBadRequest, errors.JSONDecodeErrorMessage)
		return
	}

	votedThread, err := handler.threadUsecase.VoteThread(slugOrId, *parsedVote)
	if err != nil {
		log.WithError(err).Error("vote creation error")
		if err == thread.NotFound || err == thread.AuthorNotExists {
			utilities.Resp(ctx, fasthttp.StatusNotFound, errors.JSONErrorMessage(err))
			return
		} else {
			utilities.Resp(ctx, fasthttp.StatusInternalServerError, errors.JSONErrorMessage(err))
			return
		}
	}

	body, err := json.Marshal(votedThread)
	if err != nil {
		log.WithError(err).Error(errors.JSONEncodeError)
		utilities.Resp(ctx, http.StatusInternalServerError, errors.JSONEncodeErrorMessage)
		return
	}

	utilities.Resp(ctx, fasthttp.StatusOK, body)
}
