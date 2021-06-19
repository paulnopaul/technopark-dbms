package delivery

import (
	"technopark-dbms/internal/pkg/domain"
	"technopark-dbms/internal/pkg/errors"
	"encoding/json"
	"github.com/fasthttp/router"
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"net/http"
)

type serviceHandler struct {
	serviceUsecase domain.ServiceUsecase
}

func NewServiceHandler(r *router.Router, su domain.ServiceUsecase) {
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
		// TODO message + status
		return
	}

	ctx.SetStatusCode(http.StatusOK)
}

func (handler *serviceHandler) serviceStatusHandler(ctx *fasthttp.RequestCtx) {
	status, err := handler.serviceUsecase.Status()
	if err != nil {
		log.WithError(err).Error("service get status error")
		// TODO message + status
		return
	}

	if err = json.NewEncoder(ctx).Encode(status); err != nil {
		log.WithError(err).Error(errors.JSONEncodeError)
		ctx.Error(errors.JSONEncodeErrorMessage, fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(http.StatusOK)
}
