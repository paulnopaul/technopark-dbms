package usecase

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx"
	"technopark-dbms/internal/pkg/constants"
	"technopark-dbms/internal/pkg/domain"
	"technopark-dbms/internal/pkg/post"
	"technopark-dbms/internal/pkg/utilities"
	"time"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

type postUsecase struct {
	DB     *pgx.ConnPool
	UUCase domain.UserUsecase
	FUCase domain.ForumUsecase
	TUCase domain.ThreadUsecase
}

func (p *postUsecase) GetById(id int64) (*domain.Post, error) {
	query := "select id, parent, author, message, is_edited, forum, thread, created from posts where id = $1"
	resPost := &domain.Post{}
	var created *time.Time
	err := p.DB.QueryRow(query, id).
		Scan(&resPost.ID, &resPost.Parent, &resPost.Author, &resPost.Message, &resPost.IsEdited, &resPost.Forum, &resPost.Thread, &created)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, post.NotFoundError
		}
		return nil, err
	}
	if created != nil {
		resPost.Created = created.Format(constants.TimeLayout)
	}
	return resPost, nil
}

func NewPostUsecase(db *pgx.ConnPool) domain.PostUsecase {
	return &postUsecase{
		DB: db,
	}
}

func (p *postUsecase) GetDetails(id int64, relatedUser bool, relatedForum bool, relatedThread bool) (*domain.Post, *domain.Forum, *domain.Thread, *domain.User, error) {
	resPost, err := p.GetById(id)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	var resForum *domain.Forum
	var resThread *domain.Thread
	var resUser *domain.User
	if relatedForum {
		resForum, err = p.FUCase.Details(resPost.Forum)
		if err != nil {
			return nil, nil, nil, nil, err
		}
	}
	if relatedThread {
		resThread, err = p.TUCase.GetThreadDetails(utilities.NewSlugOrId(fmt.Sprint(resPost.Thread)))
		if err != nil {
			return nil, nil, nil, nil, err
		}
	}
	if relatedUser {
		resUser, err = p.UUCase.GetProfile(resPost.Author)
		if err != nil {
			return nil, nil, nil, nil, err
		}
	}
	return resPost, resForum, resThread, resUser, nil
}

func (p *postUsecase) UpdateDetails(id int64, postUpdate domain.Post) (*domain.Post, error) {
	_, err := p.GetById(id)
	if err != nil {
		return nil, err
	}

	query := "update posts set message = $1, is_edited = true where id = $2 returning id, parent, author, message, is_edited, forum, thread, created;"
	resPost := &domain.Post{}
	err = p.DB.QueryRow(query, postUpdate.Message, id).
		Scan(&resPost.ID, &resPost.Parent, &resPost.Author, &resPost.Message, &resPost.IsEdited, &resPost.Forum, &resPost.Thread, &resPost.Created)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, post.NotFoundError
		}
		return nil, err
	}
	return resPost, nil
}
