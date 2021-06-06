package middlewares

import (
	"github.com/valyala/fasthttp"
	"log"
)

func Logging(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		log.Println(ctx.URI())
		next(ctx)
	}
}
