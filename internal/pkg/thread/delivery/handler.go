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
	"technopark-dbms/internal/pkg/post"
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
	s := r.Group("/api/thread")

	s.POST("/{slug_or_id}/create", h.threadCreatePostsHandler)
	s.GET("/{slug_or_id}/details", h.threadGetDetailsHandler)
	s.POST("/{slug_or_id}/details", h.threadUpdateDetailsHandler)
	s.GET("/{slug_or_id}/posts", h.threadGetPostsHandler)
	s.POST("/{slug_or_id}/vote", h.threadVoteHandler)
}

func (handler *threadHandler) threadCreatePostsHandler(ctx *fasthttp.RequestCtx) {
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
		if err == thread.NotFound {
			utilities.Resp(ctx, fasthttp.StatusNotFound, errors.JSONErrorMessage(err))
			return
		} else if err == post.InvalidParentError {
			utilities.Resp(ctx, fasthttp.StatusConflict, errors.JSONErrorMessage(err))
			return
		} else if err == thread.AuthorNotExists {
			utilities.Resp(ctx, fasthttp.StatusNotFound, errors.JSONErrorMessage(err))
			return
		} else {
			utilities.Resp(ctx, fasthttp.StatusInternalServerError, errors.JSONErrorMessage(err))
			return
		}
	}
	_ = json.NewEncoder(ctx).Encode(createdPosts)
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(responseStatus)
}

func (handler *threadHandler) threadGetDetailsHandler(ctx *fasthttp.RequestCtx) {
	slugOrId := utilities.NewSlugOrId(ctx.UserValue("slug_or_id").(string))
	threadDetails, err := handler.threadUsecase.GetThreadDetails(slugOrId)
	if err != nil {
		log.WithError(err).Error("thread get details error")
		if err == thread.NotFound {
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
	utilities.Resp(ctx, fasthttp.StatusOK, threadDetails)
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
	utilities.Resp(ctx, http.StatusOK, updatedThread)
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
	utilities.Resp(ctx, fasthttp.StatusOK, foundPosts)
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

	votedThread, err := handler.threadUsecase.CreateThreadVote(slugOrId, *parsedVote)
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
	utilities.Resp(ctx, fasthttp.StatusOK, votedThread)
}
