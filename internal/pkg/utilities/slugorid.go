package utilities

import "strconv"

type SlugOrId struct {
	IsSlug bool
	Slug   string
	ID     int32
}

func NewSlugOrId(idString string) SlugOrId {
	parsedID, err := strconv.ParseInt(idString, 10, 32)
	if err != nil {
		return SlugOrId{
			IsSlug: true,
			Slug:   idString,
		}
	}
	return SlugOrId{
		IsSlug: false,
		ID:     int32(parsedID),
	}
}
