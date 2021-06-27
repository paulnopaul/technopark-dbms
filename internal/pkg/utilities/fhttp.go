package utilities

import (
	"github.com/mailru/easyjson"
	"github.com/valyala/fasthttp"
)

func Resp(ctx *fasthttp.RequestCtx, code int, v interface{}) {
	_, _ = easyjson.MarshalToWriter(v.(easyjson.Marshaler), ctx)
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(code)
}
