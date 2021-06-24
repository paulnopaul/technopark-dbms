package utilities

import (
	"github.com/valyala/fasthttp"
)

func Resp(ctx *fasthttp.RequestCtx, code int, data []byte) {
	ctx.SetContentType("application/json")
	ctx.SetBody(data)
	ctx.SetStatusCode(code)
}
