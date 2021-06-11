package delivery

import (
	"DBMSForum/internal/pkg/domain"
	"DBMSForum/internal/pkg/errors"
	"DBMSForum/internal/pkg/utilities"
	"encoding/json"
	"github.com/fasthttp/router"
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"net/http"
)

type forumHandler struct {
	forumUsecase domain.ForumUsecase
}

func NewForumHandler(r *router.Router, fu domain.ForumUsecase) {
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
		ctx.Error(errors.JSONDecodeErrorMessage, fasthttp.StatusInternalServerError)
		return
	}

	createdForum, err := handler.forumUsecase.Create(*parsedForum)
	if err != nil {
		log.WithError(err).Error("forum creation error")
		// TODO error + message
		return
	}

	if err = json.NewEncoder(ctx).Encode(createdForum); err != nil {
		log.WithError(err).Error(errors.JSONEncodeError)
		ctx.Error(errors.JSONEncodeErrorMessage, fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(http.StatusCreated)
}

func (handler *forumHandler) forumDetailsHandler(ctx *fasthttp.RequestCtx) {
	slugValue := ctx.UserValue("slug").(string)
	// TODO check slug function (maybe in middlewares)

	forumDetails, err := handler.forumUsecase.Details(slugValue)
	if err != nil {
		log.WithError(err).Error("forum get details error")
		// todo error + message
		return
	}

	if err = json.NewEncoder(ctx).Encode(forumDetails); err != nil {
		log.WithError(err).Error(errors.JSONEncodeError)
		ctx.Error(errors.JSONEncodeErrorMessage, fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(http.StatusOK)
}

func (handler *forumHandler) forumCreateThreadHandler(ctx *fasthttp.RequestCtx) {
	slugValue := ctx.UserValue("slug").(string)
	parsedThread := &domain.Thread{}
	err := json.Unmarshal(ctx.PostBody(), parsedThread)
	if err != nil {
		log.WithError(err).Error(errors.JSONUnmarshallError)
		ctx.Error(errors.JSONDecodeErrorMessage, fasthttp.StatusInternalServerError)
		return
	}

	createdThread, err := handler.forumUsecase.CreateThread(slugValue, *parsedThread)
	if err != nil {
		log.WithError(err).Error("forum create thread error")
		// todo error + message
		return
	}

	if err = json.NewEncoder(ctx).Encode(createdThread); err != nil {
		log.WithError(err).Error(errors.JSONEncodeError)
		ctx.Error(errors.JSONEncodeErrorMessage, fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(http.StatusOK)
}

func (handler *forumHandler) forumGetUsersHandler(ctx *fasthttp.RequestCtx) {
	slugValue := ctx.UserValue("slug").(string)
	params, err := utilities.NewArrayOutParams(ctx.QueryArgs())
	if err != nil {
		log.WithError(err).Error(errors.QuerystringParseError)
		ctx.Error(errors.JSONQuerystringErrorMessage, errors.CodeFromJSONMessage(errors.JSONQuerystringErrorMessage))
		return
	}

	foundUsers, err := handler.forumUsecase.Users(slugValue, *params)
	if err != nil {
		log.WithError(err).Error("forum get users error")
		// todo error + message
		return
	}

	if err = json.NewEncoder(ctx).Encode(foundUsers); err != nil {
		log.WithError(err).Error(errors.JSONEncodeError)
		ctx.Error(errors.JSONEncodeErrorMessage, fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(http.StatusOK)
}

func (handler *forumHandler) forumGetThreadsHandler(ctx *fasthttp.RequestCtx) {
	slugValue := ctx.UserValue("slug").(string)
	params, err := utilities.NewArrayOutParams(ctx.URI().QueryArgs())
	if err != nil {
		log.WithError(err).Error(errors.QuerystringParseError)
		ctx.Error(errors.JSONQuerystringErrorMessage, errors.CodeFromJSONMessage(errors.JSONQuerystringErrorMessage))
		return
	}

	foundUsers, err := handler.forumUsecase.Threads(slugValue, *params)
	if err != nil {
		log.WithError(err).Error("forum get users error")
		// todo error + message
		return
	}

	if err = json.NewEncoder(ctx).Encode(foundUsers); err != nil {
		log.WithError(err).Error(errors.JSONEncodeError)
		ctx.Error(errors.JSONEncodeErrorMessage, fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(http.StatusOK)
}
