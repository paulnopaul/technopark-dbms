package utilities

import (
	"github.com/mailru/easyjson"
	"github.com/valyala/fasthttp"
)

func Resp(ctx *fasthttp.RequestCtx, code int, v easyjson.Marshaler) {
	_, _ = easyjson.MarshalToWriter(v, ctx)
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(code)
}
