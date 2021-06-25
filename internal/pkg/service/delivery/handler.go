package delivery

import (
	"encoding/json"
	"github.com/fasthttp/router"
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"net/http"
	"technopark-dbms/internal/pkg/domain"
	"technopark-dbms/internal/pkg/errors"
	"technopark-dbms/internal/pkg/utilities"
)

type serviceHandler struct {
	serviceUsecase domain.ServiceUsecase
}

func NewServiceHandler(r *router.Group, su domain.ServiceUsecase) {
	h := serviceHandler{
		serviceUsecase: su,
	}
	s := r.Group("/service")

	s.POST("/clear", h.serviceClearHandler)
	s.GET("/status", h.serviceStatusHandler)
}

func (handler *serviceHandler) serviceClearHandler(ctx *fasthttp.RequestCtx) {
	err := handler.serviceUsecase.Clear()
	if err != nil {
		log.WithError(err).Error("service clear error")
		return
	}

	ctx.SetStatusCode(http.StatusOK)
}

func (handler *serviceHandler) serviceStatusHandler(ctx *fasthttp.RequestCtx) {
	status, err := handler.serviceUsecase.Status()
	if err != nil {
		log.WithError(err).Error("service get status error")
		utilities.Resp(ctx, fasthttp.StatusInternalServerError, errors.JSONErrorMessage(err))
		return
	}

	body, err := json.Marshal(status)
	if err != nil {
		log.WithError(err).Error(errors.JSONEncodeError)
		utilities.Resp(ctx, fasthttp.StatusInternalServerError, errors.JSONEncodeErrorMessage)
		return
	}
	utilities.Resp(ctx, fasthttp.StatusOK, body)
}
