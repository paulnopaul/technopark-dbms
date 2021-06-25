package utilities

import (
	"encoding/json"
	"github.com/valyala/fasthttp"
)

func Resp(ctx *fasthttp.RequestCtx, code int, v interface{}) {
	_ = json.NewEncoder(ctx).Encode(v)
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(code)
}
