package utilities

import (
	"github.com/valyala/fasthttp"
)

type ArrayOutParams struct {
	Limit int32
	Since string
	Desc  bool
	Sort  string
}

func NewArrayOutParams(queryArgs *fasthttp.Args) (*ArrayOutParams, error) {
	res := &ArrayOutParams{}

	if queryArgs.Has("limit") {
		parsedLimit, err := queryArgs.GetUint("limit")
		if err != nil {
			return nil, err
		}
		res.Limit = int32(parsedLimit)
	} else {
		res.Limit = 100
	}

	if queryArgs.Has("since") {
		res.Since = string(queryArgs.Peek("since"))
	}

	res.Desc = queryArgs.GetBool("desc")

	if queryArgs.Has("sort") {
		res.Sort = string(queryArgs.Peek("sort"))
	}
	return res, nil
}
