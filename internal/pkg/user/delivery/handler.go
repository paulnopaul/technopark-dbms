package delivery

import (
	"DBMSForum/internal/pkg/domain"
	"github.com/fasthttp/router"
)

type Handler struct {
}

func NewUserHandler(r *router.Router, fu domain.ForumManager) {
	h := Handler{}
	s := r.Group("/user")


}
