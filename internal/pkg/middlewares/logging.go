package middlewares

import (
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

func Logging(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		log.Info(string(ctx.RequestURI()))
		next(ctx)
	}
}
