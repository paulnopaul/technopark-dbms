package usecase

import (
	"DBMSForum/internal/pkg/domain"
	"DBMSForum/internal/pkg/utilities"
	"database/sql"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

type forumUsecase struct {
	DB *sql.DB
}

func NewForumUsecase(db *sql.DB) domain.ForumUsecase {
	return &forumUsecase{
		DB: db,
	}
}

func (u *forumUsecase) Create(f domain.Forum) (*domain.Forum, error) {
	query := "insert into forums(title, username, slug) values ($1, $2, $3) returning title, username, slug, posts, threads;"
	newF := &domain.Forum{}
	err := u.DB.QueryRow(query, f.Title, f.User, f.Slug).Scan(&newF.Title, &newF.User, &newF.Slug, &newF.Posts, &newF.Threads)
	if err != nil {
		return nil, err
	}
	return newF, err
}

func (u *forumUsecase) Details(slug string) (*domain.Forum, error) {
	query := "select title, username, slug, posts, threads from forums where slug = $1"
	f := &domain.Forum{}
	err := u.DB.QueryRow(query, slug).Scan(&f.Title, &f.User, &f.Slug, &f.Posts, &f.Threads)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (u *forumUsecase) CreateThread(slug string, t domain.Thread) (*domain.Thread, error) {
	createThreadQuery := "insert into threads(title, author, message) values ($1, $2, $3) returning id, title, author,  message, votes, created;"

	newT := &domain.Thread{}
	err := u.DB.QueryRow(createThreadQuery, t.Title, t.Author, t.Message, slug).Scan(&newT.ID, &newT.Title, &newT.Author, &newT.Message, &newT.Votes, &newT.Created)
	if err != nil {
		return nil, err
	}

	createFTQuery := "insert into f_t(f_slug, t_id) values ($1, $2)"
	_, err = u.DB.Exec(createFTQuery, slug, newT.ID)
	if err != nil {
		return nil, err
	}

	newT.Forum = slug
	return newT, nil
}

func (u *forumUsecase) Users(slug string, params utilities.ArrayOutParams) ([]domain.User, error) {
	getForumUsersQuery := "select nickname, fullname, about, email from f_u join users u on f_u.u_nick = u.nickname where f_u.f_slug = $1 and u.nickname > $2 order by u_nick %s limit $3;"

	if params.Desc {
		getForumUsersQuery = fmt.Sprintf(getForumUsersQuery, "desc")
	} else {
	}
	u.DB.Query(getForumUsersQuery, slug, params.Since, params.Limit)
	return nil, nil
}

func (u *forumUsecase) Threads(slug string, params utilities.ArrayOutParams) ([]domain.Thread, error) {
	return nil, errors.New("unimplemented")
}
