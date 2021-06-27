package middlewares

import (
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"time"
)

func Logging(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		start := time.Now()
		next(ctx)
		dur := time.Since(start)
		if dur > 90*time.Millisecond {
			log.Warn(string(ctx.RequestURI()) + " " + dur.String())
		} else {
			log.Info(string(ctx.RequestURI()))
		}
	}
}
