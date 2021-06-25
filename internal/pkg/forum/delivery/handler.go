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
	"technopark-dbms/internal/pkg/utilities"
)

type forumHandler struct {
	forumUsecase domain.ForumUsecase
}

func NewForumHandler(r *router.Group, fu domain.ForumUsecase) {
	h := forumHandler{
		forumUsecase: fu,
	}
	s := r.Group("/forum")

	s.POST("/create", h.forumCreateHandler)

	s.GET("/{slug}/details", h.forumDetailsHandler)
	s.POST("/{slug}/create", h.forumCreateThreadHandler)
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
	parsedForum := &domain.Forum{}
	err := json.Unmarshal(ctx.PostBody(), parsedForum)
	if err != nil {
		log.WithError(err).Error(errors.JSONUnmarshallError)
		utilities.Resp(ctx, http.StatusBadRequest, errors.JSONDecodeErrorMessage)
		return
	}

	createdForum, err := handler.forumUsecase.CreateForum(*parsedForum)
	responseStatus := fasthttp.StatusCreated
	if err != nil {
		log.WithError(err).Error("forum creation error")
		if err == forum.AlreadyExists {
			responseStatus = fasthttp.StatusConflict
		} else if err == forum.AuthorNotExists {
			utilities.Resp(ctx, fasthttp.StatusNotFound,
				errors.JSONMessage(fmt.Sprintf("Can't find user with nickname: %s", parsedForum.User)))
			return
		} else {
			utilities.Resp(ctx, fasthttp.StatusInternalServerError, errors.JSONErrorMessage(err))
			return
		}
	}

	body, err := json.Marshal(createdForum)
	if err != nil {
		log.WithError(err).Error(errors.JSONEncodeError)
		utilities.Resp(ctx, fasthttp.StatusInternalServerError, errors.JSONEncodeErrorMessage)
		return
	}

	utilities.Resp(ctx, responseStatus, body)
}

func (handler *forumHandler) forumDetailsHandler(ctx *fasthttp.RequestCtx) {
	slugValue := ctx.UserValue("slug").(string)
	forumDetails, err := handler.forumUsecase.Details(slugValue)
	if err != nil {
		log.WithError(err).Error("forum get details error")
		if err == forum.NotFound {
			utilities.Resp(ctx,
				fasthttp.StatusNotFound,
				errors.JSONMessage(fmt.Sprintf("Can't find forum with slug: %s", slugValue)))
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

func (handler *forumHandler) forumCreateThreadHandler(ctx *fasthttp.RequestCtx) {
	slugValue := ctx.UserValue("slug").(string)
	parsedThread := &domain.Thread{}
	err := json.Unmarshal(ctx.PostBody(), parsedThread)
	if err != nil {
		log.WithError(err).Error(errors.JSONUnmarshallError)
		utilities.Resp(ctx, http.StatusBadRequest, errors.JSONDecodeErrorMessage)
		return
	}

	createdThread, err := handler.forumUsecase.CreateThread(slugValue, *parsedThread)
	respStatus := fasthttp.StatusCreated
	if err != nil {
		log.WithError(err).Error("forum create thread error")
		if err == forum.AlreadyExists {
			respStatus = fasthttp.StatusConflict
		} else if err == forum.AuthorNotExists {
			utilities.Resp(ctx, fasthttp.StatusNotFound, errors.JSONErrorMessage(err))
			return
		} else {
			utilities.Resp(ctx, fasthttp.StatusInternalServerError, errors.JSONErrorMessage(err))
			return
		}
	}

	body, err := json.Marshal(createdThread)
	if err != nil {
		log.WithError(err).Error(errors.JSONEncodeError)
		utilities.Resp(ctx, fasthttp.StatusInternalServerError, errors.JSONEncodeErrorMessage)
		return
	}

	utilities.Resp(ctx, respStatus, body)
}

func (handler *forumHandler) forumGetUsersHandler(ctx *fasthttp.RequestCtx) {
	slugValue := ctx.UserValue("slug").(string)
	params, err := utilities.NewArrayOutParams(ctx.QueryArgs())
	if err != nil {
		log.WithError(err).Error(errors.QuerystringParseError)
		utilities.Resp(ctx, fasthttp.StatusInternalServerError, errors.JSONErrorMessage(err))
		return
	}

	foundUsers, err := handler.forumUsecase.Users(slugValue, *params)
	if err != nil {
		log.WithError(err).Error("forum get users error")
		if err == forum.NotFound {
			utilities.Resp(ctx, http.StatusNotFound, errors.JSONErrorMessage(err))
			return
		}
		utilities.Resp(ctx, http.StatusInternalServerError, errors.JSONErrorMessage(err))
		return
	}

	body, err := json.Marshal(foundUsers)
	if err != nil {
		log.WithError(err).Error(errors.JSONEncodeError)
		utilities.Resp(ctx, fasthttp.StatusInternalServerError, errors.JSONEncodeErrorMessage)
		return
	}

	utilities.Resp(ctx, http.StatusOK, body)
}

func (handler *forumHandler) forumGetThreadsHandler(ctx *fasthttp.RequestCtx) {
	slugValue := ctx.UserValue("slug").(string)
	params, err := utilities.NewArrayOutParams(ctx.URI().QueryArgs())
	if err != nil {
		log.WithError(err).Error(errors.QuerystringParseError)
		utilities.Resp(ctx,
			fasthttp.StatusInternalServerError,
			errors.JSONErrorMessage(err))
		return
	}

	foundUsers, err := handler.forumUsecase.Threads(slugValue, *params)
	if err != nil {
		log.WithError(err).Error("forum get users error")
		if err == forum.NotFound {
			utilities.Resp(ctx,
				fasthttp.StatusNotFound,
				errors.JSONMessage(fmt.Sprintf("Can't find forum with slug: %s", slugValue)))
			return
		}
		utilities.Resp(ctx,
			fasthttp.StatusInternalServerError,
			errors.JSONErrorMessage(err))
	}

	body, err := json.Marshal(foundUsers)
	if err != nil {
		log.WithError(err).Error(errors.JSONEncodeError)
		utilities.Resp(ctx, fasthttp.StatusInternalServerError, errors.JSONEncodeErrorMessage)
		return
	}

	utilities.Resp(ctx, fasthttp.StatusOK, body)
}
