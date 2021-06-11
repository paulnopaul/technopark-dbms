package utilities

import (
	"DBMSForum/internal/pkg/errors"
	"github.com/valyala/fasthttp"
)

type SortType string

var SortTypes = map[SortType]bool{
	SortType("flat"):        true,
	SortType("tree"):        true,
	SortType("parent_tree"): true,
}

type ArrayOutParams struct {
	Limit    int32
	HasLimit bool
	Since    string
	Desc     bool
	Sort     SortType
	HasSort  bool
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
		res.Since = string(queryArgs.Peek("limit"))
	}

	res.Desc = queryArgs.GetBool("desc")
	if queryArgs.Has("sort") {
		parsedString := SortType(queryArgs.Peek("sort"))
		_, sortExists := SortTypes[parsedString]
		if !sortExists {
			return nil, errors.WrongSortType
		}
	}

	return res, nil
}
